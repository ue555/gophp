package php

import "strings"

func (e *evaluator) callArrayPredicate(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "callback"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects exactly 2 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, name, 1, "array", "array", bound[0])
	}
	target, ok := e.resolveCallable(bound[1], node.line)
	if !ok {
		return value{}, builtinTypeError(node.line, name, 2, "callback", "callable", bound[1])
	}
	for _, entry := range bound[0].a.entries {
		key := integer(entry.key.index)
		if entry.key.isString {
			key = value{kind: stringKind, s: entry.key.text}
		}
		result, err := e.invokeInternalCallback(target, []callArgument{{value: dereference(entry.value)}, {value: key}}, node.line)
		if err != nil {
			return value{}, err
		}
		if name == "array_any" && result.truthy() {
			return boolean(true), nil
		}
		if name == "array_all" && !result.truthy() {
			return boolean(false), nil
		}
	}
	return boolean(name == "array_all"), nil
}
