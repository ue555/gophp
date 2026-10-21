package php

import "strings"

func (e *evaluator) callStrPad(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "length"}, {name: "pad_string"}, {name: "pad_type"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || !present[1] || len(args) > 4 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "str_pad() expects at least 2 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "str_pad", 1, "string", "string", bound[0])
	}
	length, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "str_pad", 2, "length", "int", bound[1])
	}
	pad := " "
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "str_pad", 3, "pad_string", "string", bound[2])
		}
		pad = v.s
	}
	if pad == "" {
		return value{}, &runtimeError{node.line, "ValueError", "str_pad(): Argument #3 ($pad_string) must be a non-empty string"}
	}
	mode := int64(1)
	if present[3] {
		v, _, ok := e.coerceType(bound[3], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "str_pad", 4, "pad_type", "int", bound[3])
		}
		mode = v.i
	}
	if mode < 0 || mode > 2 {
		return value{}, &runtimeError{node.line, "ValueError", "str_pad(): Argument #4 ($pad_type) must be STR_PAD_LEFT, STR_PAD_RIGHT, or STR_PAD_BOTH"}
	}
	if length.i <= int64(len(s.s)) {
		return s, nil
	}
	if length.i > 1<<26 {
		return value{}, &runtimeError{node.line, "Error", "gophp: str_pad() result exceeds implementation size limit"}
	}
	n := int(length.i) - len(s.s)
	left := 0
	if mode == 0 {
		left = n
	} else if mode == 2 {
		left = n / 2
	}
	return value{kind: stringKind, s: padBytes(pad, left) + s.s + padBytes(pad, n-left)}, nil
}

func padBytes(pad string, n int) string {
	if n == 0 {
		return ""
	}
	return strings.Repeat(pad, n/len(pad)) + pad[:n%len(pad)]
}
