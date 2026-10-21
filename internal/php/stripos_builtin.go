package php

import "strings"

func (e *evaluator) callStripos(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "haystack"}, {name: "needle"}, {name: "offset"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "stripos() expects 2 to 3 arguments"}
	}
	hay, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "stripos", 1, "haystack", "string", bound[0])
	}
	needle, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "stripos", 2, "needle", "string", bound[1])
	}
	offset := int64(0)
	if present[2] {
		number, _, valid := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !valid {
			return value{}, builtinTypeError(node.line, "stripos", 3, "offset", "int", bound[2])
		}
		offset = number.i
	}
	if offset < 0 {
		offset += int64(len(hay.s))
	}
	if offset < 0 || offset > int64(len(hay.s)) {
		return value{}, &runtimeError{node.line, "ValueError", "stripos(): Argument #3 ($offset) must be contained in argument #1 ($haystack)"}
	}
	index := strings.Index(strings.ToLower(hay.s[offset:]), strings.ToLower(needle.s))
	if index < 0 {
		return boolean(false), nil
	}
	return integer(int64(index) + offset), nil
}
