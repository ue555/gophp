package php

import "strings"

func (e *evaluator) callExplode(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "separator"}, {name: "string"}, {name: "limit"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "explode() expects 2 to 3 arguments"}
	}
	for i, name := range []string{"separator", "string"} {
		v, _, ok := e.coerceType(bound[i], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "explode", i+1, name, "string", bound[i])
		}
		bound[i] = v
	}
	if bound[0].s == "" {
		return value{}, &runtimeError{node.line, "ValueError", "explode(): Argument #1 ($separator) cannot be empty"}
	}
	parts := strings.Split(bound[1].s, bound[0].s)
	if present[2] {
		limit, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "explode", 3, "limit", "int", bound[2])
		}
		n := int(limit.i)
		if n > 0 && len(parts) > n {
			parts = append(parts[:n-1], strings.Join(parts[n-1:], bound[0].s))
		} else if n < 0 {
			parts = parts[:max(0, len(parts)+n)]
		} else if n == 0 {
			parts = []string{bound[1].s}
		}
	}
	out := arrayValue()
	for i, part := range parts {
		out.a = out.a.with(arrayKey{index: int64(i)}, value{kind: stringKind, s: part})
	}
	return out, nil
}
