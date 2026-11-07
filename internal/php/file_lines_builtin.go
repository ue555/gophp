package php

import (
	"errors"
	"os"
	"strings"
	"syscall"
)

// callFileLines implements file(): the whole file split into lines that keep
// their "\n" unless FILE_IGNORE_NEW_LINES (which also drops a preceding "\r"),
// with FILE_SKIP_EMPTY_LINES removing lines that are empty after that step.
// Symfony TypeInfo reads class files this way to collect `use` statements.
func (e *evaluator) callFileLines(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}, {name: "flags"}, {name: "context"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 3 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "file() expects 1 to 3 arguments"}
	}
	path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "file", 1, "filename", "string", bound[0])
	}
	flags := int64(0)
	if present[1] {
		converted, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "file", 2, "flags", "int", bound[1])
		}
		flags = converted.i
	}
	if flags&^7 != 0 {
		return value{}, &runtimeError{node.line, "ValueError", "file(): Argument #2 ($flags) must be a valid flag value"}
	}
	if present[2] && bound[2].kind != nullKind {
		return value{}, &runtimeError{node.line, "Error", "gophp: file() stream contexts are not implemented"}
	}
	filename := path.s
	if flags&1 != 0 {
		if found := e.resolveLocalInclude(filename); found != "" {
			filename = found
		}
	}
	data, readErr := os.ReadFile(filename)
	if readErr != nil {
		if err := e.warn(node.line, "file("+path.s+"): Failed to open stream: "+streamErrorText(readErr)); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	out := arrayValue()
	text := string(data)
	for len(text) > 0 {
		line := text
		if i := strings.IndexByte(text, '\n'); i >= 0 {
			line, text = text[:i+1], text[i+1:]
		} else {
			text = ""
		}
		if flags&2 != 0 {
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")
		}
		if flags&4 != 0 && line == "" {
			continue
		}
		out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, value{kind: stringKind, s: line})
	}
	return out, nil
}

// streamErrorText renders an OS error the way PHP's stream wrapper reports it.
func streamErrorText(err error) string {
	switch {
	case errors.Is(err, syscall.ENOENT):
		return "No such file or directory"
	case errors.Is(err, syscall.EACCES):
		return "Permission denied"
	case errors.Is(err, syscall.ENOTDIR):
		return "Not a directory"
	case errors.Is(err, syscall.EISDIR):
		return "Is a directory"
	}
	return err.Error()
}
