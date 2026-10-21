package php

import "regexp"

func (e *evaluator) callMbSplit(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "pattern"}, {name: "string"}, {name: "limit"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "mb_split() expects 2 to 3 arguments"}
	}
	pattern, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "mb_split", 1, "pattern", "string", bound[0])
	}
	s, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "mb_split", 2, "string", "string", bound[1])
	}
	limit := -1
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "mb_split", 3, "limit", "int", bound[2])
		}
		limit = int(v.i)
		if limit <= 0 {
			limit = -1
		}
	}
	re, err := regexp.Compile(pattern.s)
	if err != nil {
		return boolean(false), nil
	}
	parts := re.Split(s.s, limit)
	out := arrayValue()
	for _, part := range parts {
		out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, value{kind: stringKind, s: part})
	}
	return out, nil
}
