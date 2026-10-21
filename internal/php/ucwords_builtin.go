package php

import "strings"

func (e *evaluator) callUcwords(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "separators"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "ucwords() expects 1 to 2 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "ucwords", 1, "string", "string", bound[0])
	}
	separators := " \t\r\n\f\v"
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "ucwords", 2, "separators", "string", bound[1])
		}
		separators = v.s
	}
	bytes := []byte(s.s)
	start := true
	for i, c := range bytes {
		if start && c >= 'a' && c <= 'z' {
			bytes[i] = c - 'a' + 'A'
		}
		start = strings.IndexByte(separators, c) >= 0
	}
	return value{kind: stringKind, s: string(bytes)}, nil
}
