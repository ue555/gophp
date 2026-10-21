package php

import "strings"

func (e *evaluator) callStrpbrk(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "characters"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strpbrk() expects exactly 2 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "strpbrk", 1, "string", "string", bound[0])
	}
	chars, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "strpbrk", 2, "characters", "string", bound[1])
	}
	if chars.s == "" {
		return value{}, &runtimeError{node.line, "ValueError", "strpbrk(): Argument #2 ($characters) must be a non-empty string"}
	}
	index := strings.IndexAny(s.s, chars.s)
	if index < 0 {
		return boolean(false), nil
	}
	return value{kind: stringKind, s: s.s[index:]}, nil
}
