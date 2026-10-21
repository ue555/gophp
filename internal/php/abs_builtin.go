package php

import "math"

func (e *evaluator) callAbs(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "num"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "abs() expects exactly 1 argument"}
	}
	v := bound[0]
	if v.kind == intKind {
		if v.i == math.MinInt64 {
			return floating(math.Abs(float64(v.i))), nil
		}
		if v.i < 0 {
			return integer(-v.i), nil
		}
		return v, nil
	}
	if v.kind == floatKind {
		return floating(math.Abs(v.f)), nil
	}
	converted, _, ok := e.coerceType(v, "float", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "abs", 1, "num", "int|float", v)
	}
	return floating(math.Abs(converted.f)), nil
}
