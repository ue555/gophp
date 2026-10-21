package php

import "time"

func (e *evaluator) callUSleep(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "microseconds"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "usleep() expects exactly 1 argument"}
	}
	v, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "usleep", 1, "microseconds", "int", bound[0])
	}
	if v.i < 0 {
		return value{}, &runtimeError{node.line, "ValueError", "usleep(): Argument #1 ($microseconds) must be greater than or equal to 0"}
	}
	time.Sleep(time.Duration(v.i) * time.Microsecond)
	return value{kind: nullKind}, nil
}
