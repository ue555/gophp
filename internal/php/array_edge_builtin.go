package php

import "strings"

// PHP 8.5's array_first()/array_last() return the value at an insertion-order
// edge without changing the array's internal iterator or sharing its binding.
func (e *evaluator) callArrayEdgeValue(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects exactly 1 argument"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, name, 1, "array", "array", bound[0])
	}
	if bound[0].a == nil || len(bound[0].a.entries) == 0 {
		return value{kind: nullKind}, nil
	}
	index := 0
	if name == "array_last" {
		index = len(bound[0].a.entries) - 1
	}
	return dereference(bound[0].a.entries[index].value), nil
}
