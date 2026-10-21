package php

import (
	"fmt"
	"strings"
)

func urlCodecBuiltin(name string) bool {
	switch strings.ToLower(name) {
	case "urldecode", "rawurldecode", "urlencode", "rawurlencode":
		return true
	}
	return false
}

func urlHex(b byte) int {
	switch {
	case b >= '0' && b <= '9':
		return int(b - '0')
	case b >= 'A' && b <= 'F':
		return int(b - 'A' + 10)
	case b >= 'a' && b <= 'f':
		return int(b - 'a' + 10)
	}
	return -1
}

func decodeURLString(source string, plusSpace bool) string {
	var out strings.Builder
	for i := 0; i < len(source); i++ {
		if source[i] == '%' && i+2 < len(source) {
			high, low := urlHex(source[i+1]), urlHex(source[i+2])
			if high >= 0 && low >= 0 {
				out.WriteByte(byte(high<<4 | low))
				i += 2
				continue
			}
		}
		if source[i] == '+' && plusSpace {
			out.WriteByte(' ')
		} else {
			out.WriteByte(source[i])
		}
	}
	return out.String()
}

func (e *evaluator) callURLCodec(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects exactly 1 argument"}
	}
	v, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, node.name, 1, "string", "string", bound[0])
	}
	name := strings.ToLower(node.name)
	var out strings.Builder
	if name == "urldecode" || name == "rawurldecode" {
		return value{kind: stringKind, s: decodeURLString(v.s, name == "urldecode")}, nil
	} else {
		for i := 0; i < len(v.s); i++ {
			b := v.s[i]
			if b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b == '_' || b == '.' || b == '~' && name == "rawurlencode" {
				out.WriteByte(b)
			} else if b == ' ' && name == "urlencode" {
				out.WriteByte('+')
			} else {
				out.WriteString(fmt.Sprintf("%%%02X", b))
			}
		}
	}
	return value{kind: stringKind, s: out.String()}, nil
}
