package php

import "strings"

func (e *evaluator) callPregQuote(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "str"}, {name: "delimiter"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "preg_quote() expects 1 to 2 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "preg_quote", 1, "str", "string", bound[0])
	}
	delimiter := ""
	if present[1] && bound[1].kind != nullKind {
		d, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_quote", 2, "delimiter", "?string", bound[1])
		}
		delimiter = d.s
	}
	var out strings.Builder
	for _, c := range s.s {
		if strings.ContainsRune(`.\\+*?[^]$(){}=!<>|:-#`, c) || (delimiter != "" && strings.ContainsRune(delimiter, c)) {
			out.WriteByte('\\')
		}
		out.WriteRune(c)
	}
	return value{kind: stringKind, s: out.String()}, nil
}
