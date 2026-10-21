package php

import "strings"

func (e *evaluator) callHTMLSpecialChars(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "flags"}, {name: "encoding"}, {name: "double_encode"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 4 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "htmlspecialchars() expects 1 to 4 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "htmlspecialchars", 1, "string", "string", bound[0])
	}
	flags := int64(11)
	if present[1] {
		flags = bound[1].i
	}
	double := !present[3] || bound[3].truthy()
	result := s.s
	if double {
		result = strings.ReplaceAll(result, "&", "&amp;")
	}
	result = strings.ReplaceAll(result, "<", "&lt;")
	result = strings.ReplaceAll(result, ">", "&gt;")
	if flags&2 != 0 {
		result = strings.ReplaceAll(result, "\"", "&quot;")
	}
	if flags&1 != 0 {
		result = strings.ReplaceAll(result, "'", "&#039;")
	}
	return value{kind: stringKind, s: result}, nil
}
