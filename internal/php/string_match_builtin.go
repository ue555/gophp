package php

import "strings"

func (e *evaluator) callStrMatch(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "haystack"}, {name: "needle"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects exactly 2 arguments"}
	}
	for i, name := range []string{"haystack", "needle"} {
		v, _, ok := e.coerceType(bound[i], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, node.name, i+1, name, "string", bound[i])
		}
		bound[i] = v
	}
	switch strings.ToLower(node.name) {
	case "str_contains":
		return boolean(strings.Contains(bound[0].s, bound[1].s)), nil
	case "str_starts_with":
		return boolean(strings.HasPrefix(bound[0].s, bound[1].s)), nil
	default:
		return boolean(strings.HasSuffix(bound[0].s, bound[1].s)), nil
	}
}
