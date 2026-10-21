package php

import "strings"

func (e *evaluator) callStrRepeat(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "times"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "str_repeat() expects exactly 2 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "str_repeat", 1, "string", "string", bound[0])
	}
	times, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "str_repeat", 2, "times", "int", bound[1])
	}
	if times.i < 0 {
		return value{}, &runtimeError{node.line, "ValueError", "str_repeat(): Argument #2 ($times) must be greater than or equal to 0"}
	}
	if times.i > 0 && int64(len(s.s)) > 0 && times.i > int64(1<<30)/int64(len(s.s)) {
		return value{}, &runtimeError{node.line, "Error", "gophp: str_repeat() result exceeds implementation size limit"}
	}
	return value{kind: stringKind, s: strings.Repeat(s.s, int(times.i))}, nil
}
