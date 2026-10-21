package php

import "math"

func (e *evaluator) callRound(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "num"}, {name: "precision"}, {name: "mode"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 3 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "round() expects 1 to 3 arguments"}
	}
	number, _, ok := e.coerceType(bound[0], "float", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "round", 1, "num", "int|float", bound[0])
	}
	precision := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "round", 2, "precision", "int", bound[1])
		}
		precision = v.i
	}
	if present[2] {
		return value{}, &runtimeError{node.line, "Error", "gophp: round() modes are not implemented"}
	}
	factor := math.Pow(10, float64(precision))
	return floating(math.Round(number.f*factor) / factor), nil
}
