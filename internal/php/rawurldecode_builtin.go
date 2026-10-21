package php

import (
	"net/url"
	"strings"
)

func (e *evaluator) callRawURLEncode(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "rawurlencode() expects exactly 1 argument"}
	}
	v, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "rawurlencode", 1, "string", "string", bound[0])
	}
	s := url.QueryEscape(v.s)
	s = strings.ReplaceAll(s, "+", "%20")
	return value{kind: stringKind, s: s}, nil
}

func (e *evaluator) callRawURLDecode(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "rawurldecode() expects exactly 1 argument"}
	}
	v, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "rawurldecode", 1, "string", "string", bound[0])
	}
	var out strings.Builder
	for i := 0; i < len(v.s); i++ {
		if v.s[i] == '%' && i+2 < len(v.s) {
			a, b := hexNibble(v.s[i+1]), hexNibble(v.s[i+2])
			if a >= 0 && b >= 0 {
				out.WriteByte(byte(a<<4 | b))
				i += 2
				continue
			}
		}
		out.WriteByte(v.s[i])
	}
	return value{kind: stringKind, s: out.String()}, nil
}

func hexNibble(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	}
	return -1
}
