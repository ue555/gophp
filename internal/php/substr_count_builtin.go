package php

import "strings"

func (e *evaluator) callSubstrCount(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	b, p, _, err := bindCallArguments(args, []parameter{{name: "haystack"}, {name: "needle"}, {name: "offset"}, {name: "length"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 4 || !p[0] || !p[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "substr_count() expects 2 to 4 arguments"}
	}
	for i, n := range []string{"haystack", "needle"} {
		v, _, ok := e.coerceType(b[i], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "substr_count", i+1, n, "string", b[i])
		}
		b[i] = v
	}
	start := int64(0)
	length := int64(len(b[0].s))
	specified := p[3] && b[3].kind != nullKind
	for i, n := range []string{"offset", "length"} {
		j := i + 2
		if !p[j] || j == 3 && !specified {
			continue
		}
		v, _, ok := e.coerceType(b[j], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "substr_count", j+1, n, "int", b[j])
		}
		if j == 2 {
			start = v.i
		} else {
			length = v.i
		}
	}
	if b[1].s == "" {
		word := "must not"
		if e.version == PHP83 {
			word = "cannot"
		}
		return value{}, &runtimeError{node.line, "ValueError", "substr_count(): Argument #2 ($needle) " + word + " be empty"}
	}
	size := int64(len(b[0].s))
	if start < 0 {
		start += size
	}
	if start < 0 || start > size {
		return value{}, &runtimeError{node.line, "ValueError", "substr_count(): Argument #3 ($offset) must be contained in argument #1 ($haystack)"}
	}
	end := size
	if specified {
		if length < 0 {
			end = size + length
		} else {
			if length > size-start {
				return value{}, &runtimeError{node.line, "ValueError", "substr_count(): Argument #4 ($length) must be contained in argument #1 ($haystack)"}
			}
			end = start + length
		}
		if end < start {
			return value{}, &runtimeError{node.line, "ValueError", "substr_count(): Argument #4 ($length) must be contained in argument #1 ($haystack)"}
		}
	}
	return integer(int64(strings.Count(b[0].s[start:end], b[1].s))), nil
}
