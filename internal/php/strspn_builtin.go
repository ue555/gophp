package php

import "strings"

func (e *evaluator) callStrspn(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "characters"}, {name: "offset"}, {name: "length"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 4 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strspn() expects 2 to 4 arguments"}
	}
	fields := make([]string, 2)
	for i, name := range []string{"string", "characters"} {
		v, _, ok := e.coerceType(bound[i], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "strspn", i+1, name, "string", bound[i])
		}
		fields[i] = v.s
	}
	start, end := int64(0), int64(len(fields[0]))
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "strspn", 3, "offset", "int", bound[2])
		}
		start = v.i
	}
	if start < 0 {
		start += end
	}
	if start < 0 {
		start = 0
	}
	if start > end {
		start = end
	}
	if present[3] && bound[3].kind != nullKind {
		v, _, ok := e.coerceType(bound[3], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "strspn", 4, "length", "?int", bound[3])
		}
		if v.i >= 0 {
			end = start + v.i
		} else {
			end += v.i
		}
	}
	if end < start {
		end = start
	}
	if end > int64(len(fields[0])) {
		end = int64(len(fields[0]))
	}
	count := int64(0)
	for i := start; i < end; i++ {
		if !strings.ContainsRune(fields[1], rune(fields[0][i])) {
			break
		}
		count++
	}
	return integer(count), nil
}
