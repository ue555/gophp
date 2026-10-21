package php

import (
	"os"
	"path/filepath"
)

func (e *evaluator) callRealpath(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "path"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "realpath() expects exactly 1 argument"}
	}
	path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "realpath", 1, "path", "string", bound[0])
	}
	name := path.s
	if name == "" {
		name = "."
	}
	abs, err := filepath.Abs(name)
	if err != nil {
		return boolean(false), nil
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return boolean(false), nil
	}
	if _, err = os.Stat(resolved); err != nil {
		return boolean(false), nil
	}
	return value{kind: stringKind, s: resolved}, nil
}
