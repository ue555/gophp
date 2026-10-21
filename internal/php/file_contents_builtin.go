package php

import "os"

func (e *evaluator) callFileGetContents(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "filename"}, {name: "use_include_path"}, {name: "context"}, {name: "offset"}, {name: "length"}}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 5 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "file_get_contents() expects 1 to 5 arguments"}
	}
	path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "file_get_contents", 1, "filename", "string", bound[0])
	}
	if present[2] && bound[2].kind != nullKind {
		return value{}, &runtimeError{node.line, "Error", "gophp: file_get_contents stream contexts are not implemented"}
	}
	if path.s == "php://input" {
		return value{kind: stringKind, s: e.requestBody}, nil
	}
	filename := path.s
	if present[1] && bound[1].truthy() {
		if found := e.resolveLocalInclude(filename); found != "" {
			filename = found
		}
	}
	data, readErr := os.ReadFile(filename)
	if readErr != nil {
		if err := e.warn(node.line, "file_get_contents("+path.s+"): Failed to open stream: "+readErr.Error()); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	start := int64(0)
	if present[3] {
		offset, _, ok := e.coerceType(bound[3], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "file_get_contents", 4, "offset", "int", bound[3])
		}
		start = offset.i
	}
	if start < 0 {
		start += int64(len(data))
	}
	if start < 0 || start > int64(len(data)) {
		return boolean(false), nil
	}
	data = data[start:]
	if present[4] {
		length, _, ok := e.coerceType(bound[4], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "file_get_contents", 5, "length", "?int", bound[4])
		}
		if length.i >= 0 && length.i < int64(len(data)) {
			data = data[:length.i]
		}
	}
	return value{kind: stringKind, s: string(data)}, nil
}
