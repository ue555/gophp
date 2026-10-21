package php

import "strings"

func (e *evaluator) callMinMax(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) == 0 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects at least 1 argument"}
	}
	values := make([]value, 0, len(args))
	if len(args) == 1 {
		if args[0].value.kind != arrayKind {
			return value{}, &runtimeError{node.line, "TypeError", node.name + "(): Argument #1 ($value) must be of type array when passed as a single argument"}
		}
		for _, entry := range args[0].value.a.entries {
			values = append(values, dereference(entry.value))
		}
	} else {
		for _, arg := range args {
			values = append(values, arg.value)
		}
	}
	if len(values) == 0 {
		return value{}, &runtimeError{node.line, "ValueError", node.name + "(): Argument #1 ($value) must contain at least one element"}
	}
	best := values[0]
	for _, candidate := range values[1:] {
		order, err := e.valueOrder(candidate, best, node.line, make(map[*phpArray]bool))
		if err != nil {
			return value{}, err
		}
		if e.pendingDiagnostic != nil {
			return value{}, e.pendingDiagnostic
		}
		if strings.EqualFold(node.name, "min") && order < 0 || strings.EqualFold(node.name, "max") && order > 0 {
			best = candidate
		}
	}
	return best, nil
}
