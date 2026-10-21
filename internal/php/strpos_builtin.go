package php

import "strings"

func (e *evaluator) callStrpos(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "haystack"}, {name: "needle"}, {name: "offset"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strpos() expects 2 to 3 arguments"}
	}
	for i, name := range []string{"haystack", "needle"} {
		v, _, ok := e.coerceType(bound[i], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "strpos", i+1, name, "string", bound[i])
		}
		bound[i] = v
	}
	start := int64(0)
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "strpos", 3, "offset", "int", bound[2])
		}
		start = v.i
	}
	if start < 0 {
		start += int64(len(bound[0].s))
	}
	if start < 0 || start > int64(len(bound[0].s)) {
		return value{}, &runtimeError{node.line, "ValueError", "strpos(): Argument #3 ($offset) must be contained in argument #1 ($haystack)"}
	}
	index := strings.Index(bound[0].s[start:], bound[1].s)
	if index < 0 {
		return boolean(false), nil
	}
	return integer(start + int64(index)), nil
}
