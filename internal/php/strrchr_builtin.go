package php

import "strings"

func (e *evaluator) callStrrchr(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "haystack"}, {name: "needle"}, {name: "before_needle"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strrchr() expects 2 to 3 arguments"}
	}
	h, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "strrchr", 1, "haystack", "string", bound[0])
	}
	n, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "strrchr", 2, "needle", "string", bound[1])
	}
	index := len(h.s)
	if n.s != "" {
		index = strings.LastIndexByte(h.s, n.s[0])
	}
	if index < 0 {
		return boolean(false), nil
	}
	if present[2] && bound[2].truthy() {
		return value{kind: stringKind, s: h.s[:index]}, nil
	}
	return value{kind: stringKind, s: h.s[index:]}, nil
}
