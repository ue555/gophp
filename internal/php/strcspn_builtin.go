package php

import "strings"

func (e *evaluator) callStrcspn(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "characters"}, {name: "offset"}, {name: "length"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || !present[1] || len(args) > 4 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strcspn() expects 2 to 4 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "strcspn", 1, "string", "string", bound[0])
	}
	chars, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "strcspn", 2, "characters", "string", bound[1])
	}
	start, end := int64(0), int64(len(s.s))
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "strcspn", 3, "offset", "int", bound[2])
		}
		start = v.i
	}
	if start < 0 {
		start += int64(len(s.s))
	}
	if start < 0 || start > int64(len(s.s)) {
		return value{}, &runtimeError{node.line, "ValueError", "strcspn(): Argument #3 ($offset) must be contained in argument #1 ($string)"}
	}
	if present[3] && bound[3].kind != nullKind {
		v, _, ok := e.coerceType(bound[3], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "strcspn", 4, "length", "?int", bound[3])
		}
		if v.i >= 0 {
			end = start + v.i
		} else {
			end += v.i
		}
	}
	if end > int64(len(s.s)) {
		end = int64(len(s.s))
	}
	if end < start {
		end = start
	}
	count := int64(0)
	for i := start; i < end; i++ {
		if strings.IndexByte(chars.s, s.s[i]) >= 0 {
			break
		}
		count++
	}
	return integer(count), nil
}
