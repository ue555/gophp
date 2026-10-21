package php

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// sourceError keeps the origin of an error after execution returns to its
// caller's file. An include and a function call both restore that context.
type sourceError struct {
	filename string
	err      error
}

// A compile-time declaration failure in a required file has no include call
// frame in PHP 8.5's diagnostic trace. Runtime failures keep that frame.
func (e *evaluator) captureIncludeCompileFailure() {
	if e.failureFrames != nil {
		return
	}
	e.captureFailure()
	e.failureFrames = e.failureFrames[:len(e.failureFrames)-1]
}

func (e *sourceError) Error() string { return e.err.Error() }
func (e *sourceError) Unwrap() error { return e.err }

func withSource(err error, filename string) error {
	if err == nil {
		return nil
	}
	var existing *sourceError
	if errors.As(err, &existing) {
		return err
	}
	return &sourceError{filename: filename, err: err}
}

func (e *evaluator) evalInclude(node includeExpr) (value, error) {
	pathValue, err := e.eval(node.path)
	if err != nil {
		return value{}, err
	}
	path, err := e.toText(pathValue, node.line)
	if err != nil {
		return value{}, err
	}

	filename := e.resolveLocalInclude(path)
	if filename == "" {
		if err := e.warn(node.line, fmt.Sprintf("%s(%s): Failed to open stream: No such file or directory", node.kind, path)); err != nil {
			return value{}, err
		}
		if strings.HasPrefix(node.kind, "require") {
			return value{}, &runtimeError{node.line, "Error", fmt.Sprintf("Failed opening required '%s' (include_path='%s')", path, e.includePath)}
		}
		if err := e.warn(node.line, fmt.Sprintf("%s(): Failed opening '%s' for inclusion (include_path='%s')", node.kind, path, e.includePath)); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	// An include is a call-site in PHP traces. Keep the caller's file and
	// line while compiling and executing the child, then remove the frame on
	// every exit path (including a caught exception).
	e.frames = append(e.frames, callFrame{method: node.kind, filename: e.filename, line: node.line})
	defer func() { e.frames = e.frames[:len(e.frames)-1] }()
	if strings.HasSuffix(node.kind, "_once") && e.loadedFiles[filename] {
		return boolean(true), nil
	}
	if e.includeDepth >= 256 {
		return value{}, &runtimeError{node.line, "Error", "gophp: maximum include depth exceeded (256)"}
	}
	source, err := os.ReadFile(filename)
	if err != nil {
		return value{}, &runtimeError{node.line, "Error", err.Error()}
	}
	once := strings.HasSuffix(node.kind, "_once")
	// The *_once forms record the opened file before compilation. Ordinary
	// include/require record it only after compilation succeeds.
	if once {
		e.loadedFiles[filename] = true
	}
	tokens, err := lex(string(source))
	var statements []statement
	strictTypes := false
	if err == nil {
		p := parser{tokens: tokens, version: e.version, filename: filename}
		statements, err = p.program()
		strictTypes = p.strictTypes
	}
	if err != nil {
		var parseFailure *syntaxError
		if errors.As(err, &parseFailure) {
			// The main file is parsed before execution and reports a parse
			// diagnostic. A file compiled by include runs inside the caller's
			// try/catch, where PHP raises a catchable ParseError instead.
			e.captureFailure()
			return value{}, withSource(&runtimeError{parseFailure.line, "ParseError", parseFailure.message}, filename)
		}
		return value{}, withSource(err, filename)
	}
	if failure := validateStatements(statements); failure != nil {
		e.captureIncludeCompileFailure()
		return value{}, withSource(&fatalDiagnostic{failure.line, failure.message}, filename)
	}
	if failure := duplicateConstants(statements); failure != nil {
		e.captureIncludeCompileFailure()
		return value{}, withSource(&fatalDiagnostic{failure.line, failure.message}, filename)
	}
	if failure := duplicateFunctions(statements, filename, e.version); failure != nil {
		e.captureIncludeCompileFailure()
		return value{}, withSource(&fatalDiagnostic{failure.line, failure.message}, filename)
	}
	if !once {
		e.loadedFiles[filename] = true
	}

	outerFilename := e.filename
	outerStrict := e.strictTypes
	e.filename = filename
	e.strictTypes = strictTypes
	e.includeDepth++
	defer func() { e.includeDepth--; e.filename = outerFilename; e.strictTypes = outerStrict }()
	e.callableDiagnostics(statements)
	e.controlDiagnostics(statements, nil)
	err = e.executeProgram(statements)
	if returned, ok := err.(*returned); ok {
		return returned.value, nil
	}
	if err != nil {
		if e.failureFrames == nil {
			e.captureFailure()
		}
		return value{}, withSource(err, filename)
	}
	return integer(1), nil
}

// resolveLocalInclude follows PHP's Unix local-file search: explicit paths
// bypass include_path; ordinary names search each entry, then the currently
// executing file's directory. Relative include_path entries use process CWD.
func (e *evaluator) resolveLocalInclude(path string) string {
	if path == "" || strings.Contains(path, "://") || strings.IndexByte(path, 0) >= 0 {
		return ""
	}
	resolve := func(candidate string) string {
		absolute, err := filepath.Abs(candidate)
		if err != nil {
			return ""
		}
		canonical, err := filepath.EvalSymlinks(absolute)
		if err != nil {
			return ""
		}
		info, err := os.Stat(canonical)
		if err != nil || info.IsDir() {
			return ""
		}
		return canonical
	}
	if filepath.IsAbs(path) || strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") || e.includePath == "" {
		return resolve(path)
	}
	entries := strings.Split(e.includePath, ":")
	for index, directory := range entries {
		// PHP stops when the scan pointer reaches the terminating byte: a
		// trailing separator does not add a final empty search directory.
		if directory == "" && index == len(entries)-1 {
			break
		}
		candidate := filepath.Join(directory, path)
		if directory == "" {
			candidate = "/" + path
		}
		if filename := resolve(candidate); filename != "" {
			return filename
		}
	}
	return resolve(filepath.Join(filepath.Dir(e.filename), path))
}
