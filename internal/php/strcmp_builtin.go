package php

import "strings"

func (e *evaluator) callStringCompare(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string1"}, {name: "string2"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects exactly 2 arguments"}
	}
	values := make([]string, 2)
	for i, name := range []string{"string1", "string2"} {
		v, _, ok := e.coerceType(bound[i], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, node.name, i+1, name, "string", bound[i])
		}
		values[i] = v.s
	}
	if strings.EqualFold(node.name, "strcasecmp") {
		values[0], values[1] = strings.ToLower(values[0]), strings.ToLower(values[1])
	}
	return integer(int64(strings.Compare(values[0], values[1]))), nil
}
