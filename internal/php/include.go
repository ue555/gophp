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

	// Local files are the initial supported domain. Relative paths are tried
	// from the working directory, then from the including file's directory.
	// Stream wrappers and configurable include_path are not yet implemented.
	candidates := []string{path}
	if !filepath.IsAbs(path) && !strings.HasPrefix(path, "."+string(filepath.Separator)) && !strings.HasPrefix(path, ".."+string(filepath.Separator)) {
		candidates = append(candidates, filepath.Join(filepath.Dir(e.filename), path))
	}
	var filename string
	for _, candidate := range candidates {
		absolute, absErr := filepath.Abs(candidate)
		if absErr != nil {
			continue
		}
		if canonical, canonicalErr := filepath.EvalSymlinks(absolute); canonicalErr == nil {
			absolute = canonical
		}
		info, statErr := os.Stat(absolute)
		if statErr == nil && !info.IsDir() {
			filename = absolute
			break
		}
	}
	if filename == "" {
		e.warn(node.line, fmt.Sprintf("%s(%s): Failed to open stream: No such file or directory", node.kind, path))
		if strings.HasPrefix(node.kind, "require") {
			return value{}, &runtimeError{node.line, "Error", fmt.Sprintf("Failed opening required '%s'", path)}
		}
		e.warn(node.line, fmt.Sprintf("%s(): Failed opening '%s' for inclusion", node.kind, path))
		return boolean(false), nil
	}
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
	tokens, err := lex(string(source))
	var statements []statement
	strictTypes := false
	if err == nil {
		p := parser{tokens: tokens, version: e.version, filename: filename}
		statements, err = p.program()
		strictTypes = p.strictTypes
	}
	if err != nil {
		return value{}, withSource(err, filename)
	}
	if failure := validateStatements(statements); failure != nil {
		return value{}, withSource(failure, filename)
	}
	if failure := duplicateConstants(statements); failure != nil {
		return value{}, withSource(failure, filename)
	}
	if failure := duplicateFunctions(statements, filename, e.version); failure != nil {
		return value{}, withSource(failure, filename)
	}

	// Mark before evaluation: a recursive require_once of the current file is
	// skipped, and a later retry after a runtime failure remains skipped.
	e.loadedFiles[filename] = true
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
		return value{}, withSource(err, filename)
	}
	return integer(1), nil
}
