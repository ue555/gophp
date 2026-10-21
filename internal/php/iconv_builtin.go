package php

import (
	"strings"
	"unicode/utf8"
)

func normalizedEncoding(name string) string {
	name = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(name, "-", ""), "_", ""))
	switch name {
	case "UTF8":
		return "UTF8"
	case "ASCII", "USASCII":
		return "ASCII"
	case "ISO88591", "LATIN1":
		return "LATIN1"
	}
	return ""
}

func (e *evaluator) callIconv(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "from_encoding"}, {name: "to_encoding"}, {name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 3 || !present[0] || !present[1] || !present[2] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "iconv() expects exactly 3 arguments"}
	}
	converted := make([]value, 3)
	for i, name := range []string{"from_encoding", "to_encoding", "string"} {
		v, _, ok := e.coerceType(bound[i], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "iconv", i+1, name, "string", bound[i])
		}
		converted[i] = v
	}
	from := normalizedEncoding(converted[0].s)
	target := converted[1].s
	ignoreInvalid := strings.HasSuffix(strings.ToUpper(target), "//IGNORE")
	if ignoreInvalid {
		target = target[:len(target)-len("//IGNORE")]
	}
	to := normalizedEncoding(target)
	if from == "" || to == "" {
		return value{}, &runtimeError{node.line, "Error", "gophp: iconv encoding is not implemented: " + converted[0].s + " to " + converted[1].s}
	}
	input := converted[2].s
	if input == "" {
		return value{kind: stringKind, s: ""}, nil
	}
	var runes []rune
	switch from {
	case "UTF8":
		if ignoreInvalid {
			input = strings.ToValidUTF8(input, "")
		}
		if !utf8.ValidString(input) {
			if err := e.warn(node.line, "iconv(): Detected an illegal character in input string"); err != nil {
				return value{}, err
			}
			return boolean(false), nil
		}
		runes = []rune(input)
	case "ASCII", "LATIN1":
		for i := 0; i < len(input); i++ {
			if from == "ASCII" && input[i] > 127 {
				if err := e.warn(node.line, "iconv(): Detected an illegal character in input string"); err != nil {
					return value{}, err
				}
				return boolean(false), nil
			}
			runes = append(runes, rune(input[i]))
		}
	}
	if to == "UTF8" {
		return value{kind: stringKind, s: string(runes)}, nil
	}
	var out strings.Builder
	maximum := rune(255)
	if to == "ASCII" {
		maximum = 127
	}
	for _, r := range runes {
		if r > maximum {
			if ignoreInvalid {
				continue
			}
			if err := e.warn(node.line, "iconv(): Detected an illegal character in input string"); err != nil {
				return value{}, err
			}
			return boolean(false), nil
		}
		out.WriteByte(byte(r))
	}
	return value{kind: stringKind, s: out.String()}, nil
}
