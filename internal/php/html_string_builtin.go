package php

import (
	"html"
	"strings"
)

func (e *evaluator) callHTMLStringBuiltin(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "string"}, {name: "allowed_tags"}}
	if strings.EqualFold(node.name, "html_entity_decode") {
		params[1].name = "flags"
		params = append(params, parameter{name: "encoding"})
	}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects at least 1 argument"}
	}
	if len(args) > len(params) {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects at most " + string(rune('0'+len(params))) + " arguments"}
	}
	v, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, node.name, 1, "string", "string", bound[0])
	}
	if strings.EqualFold(node.name, "html_entity_decode") {
		return value{kind: stringKind, s: html.UnescapeString(v.s)}, nil
	}
	allowed := map[string]bool{}
	if present[1] && bound[1].kind == stringKind {
		for _, token := range domTokens.FindAllString(bound[1].s, -1) {
			if strings.HasPrefix(token, "<") {
				allowed[strings.ToLower(strings.Trim(token, "<>/ "))] = true
			}
		}
	}
	var out strings.Builder
	last := 0
	for _, index := range domTokens.FindAllStringIndex(v.s, -1) {
		out.WriteString(v.s[last:index[0]])
		token := v.s[index[0]:index[1]]
		if !strings.HasPrefix(token, "<") {
			out.WriteString(token)
		} else if parts := strings.Fields(strings.Trim(token, "<>/ ")); len(parts) > 0 && allowed[strings.ToLower(parts[0])] {
			out.WriteString(token)
		}
		last = index[1]
	}
	out.WriteString(v.s[last:])
	return value{kind: stringKind, s: out.String()}, nil
}
