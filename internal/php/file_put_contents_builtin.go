package php

import (
	"os"
	"path/filepath"
)

func (e *evaluator) callFilePutContents(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}, {name: "data"}, {name: "flags"}, {name: "context"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 4 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "file_put_contents() expects between 2 and 4 arguments"}
	}
	path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "file_put_contents", 1, "filename", "string", bound[0])
	}
	data, err := e.toText(bound[1], node.line)
	if err != nil {
		return value{}, err
	}
	flags := int64(0)
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "file_put_contents", 3, "flags", "int", bound[2])
		}
		flags = v.i
	}
	if present[3] {
		return value{}, &runtimeError{node.line, "Error", "gophp: file_put_contents() stream contexts are not implemented"}
	}
	filename := path.s
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(filepath.Dir(e.filename), filename)
	}
	mode := os.O_CREATE | os.O_WRONLY
	if flags&8 != 0 {
		mode |= os.O_APPEND
	} else {
		mode |= os.O_TRUNC
	}
	f, err := os.OpenFile(filename, mode, 0666)
	if err != nil {
		if warning := e.warn(node.line, "file_put_contents("+path.s+"): Failed to open stream: "+err.Error()); warning != nil {
			return value{}, warning
		}
		return boolean(false), nil
	}
	count, writeErr := f.WriteString(data)
	closeErr := f.Close()
	if writeErr != nil || closeErr != nil {
		return boolean(false), nil
	}
	return integer(int64(count)), nil
}
