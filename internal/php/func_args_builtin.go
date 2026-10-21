package php

import "strings"

func (e *evaluator) callFuncArgs(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	max := 0
	if name == "func_get_arg" {
		max = 1
	}
	if len(args) != max {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects exactly " + string(rune('0'+max)) + " arguments"}
	}
	if len(e.frames) == 0 {
		return value{}, &runtimeError{node.line, "Error", name + "() must be called from a function context"}
	}
	frame := e.frames[len(e.frames)-1]
	if name == "func_num_args" {
		return integer(int64(len(frame.callArgs))), nil
	}
	if name == "func_get_args" {
		out := arrayValue()
		for i, arg := range frame.callArgs {
			out.a = out.a.with(arrayKey{index: int64(i)}, arg.value)
		}
		return out, nil
	}
	index, _, ok := e.coerceType(args[0].value, "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, name, 1, "position", "int", args[0].value)
	}
	if index.i < 0 || index.i >= int64(len(frame.callArgs)) {
		return value{}, &runtimeError{node.line, "ValueError", "func_get_arg(): Argument #1 ($position) must be less than the number of arguments passed to the currently executed function"}
	}
	return frame.callArgs[index.i].value, nil
}
