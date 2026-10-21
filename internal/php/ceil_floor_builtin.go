package php

import (
	"math"
	"strings"
)

func (e *evaluator) callCeilFloor(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "num"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects exactly 1 argument"}
	}
	num, _, ok := e.coerceType(bound[0], "float", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, name, 1, "num", "int|float", bound[0])
	}
	if name == "ceil" {
		return floating(math.Ceil(num.f)), nil
	}
	return floating(math.Floor(num.f)), nil
}
