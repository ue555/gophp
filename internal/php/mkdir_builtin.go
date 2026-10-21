package php

import "os"

func (e *evaluator) callMkdir(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "directory"}, {name: "permissions"}, {name: "recursive"}, {name: "context"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 4 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "mkdir() expects 1 to 4 arguments"}
	}
	path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "mkdir", 1, "directory", "string", bound[0])
	}
	mode := int64(0777)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "mkdir", 2, "permissions", "int", bound[1])
		}
		mode = v.i
	}
	if present[3] && bound[3].kind != nullKind {
		return value{}, &runtimeError{node.line, "Error", "gophp: mkdir stream context is not implemented"}
	}
	var createErr error
	if present[2] && bound[2].truthy() {
		createErr = os.MkdirAll(path.s, os.FileMode(mode))
	} else {
		createErr = os.Mkdir(path.s, os.FileMode(mode))
	}
	if createErr != nil {
		if err := e.warn(node.line, "mkdir(): "+createErr.Error()); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	return boolean(true), nil
}
