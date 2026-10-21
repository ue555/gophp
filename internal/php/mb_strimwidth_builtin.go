package php

import "strings"

func mbRuneWidth(r rune) int {
	if r >= 0x1100 && (r <= 0x115f || r >= 0x2e80 && r <= 0xa4cf || r >= 0xac00 && r <= 0xd7a3 || r >= 0xf900 && r <= 0xfaff || r >= 0xfe10 && r <= 0xfe6b || r >= 0xff01 && r <= 0xff60 || r >= 0xffe0 && r <= 0xffe6 || r >= 0x1f200 && r <= 0x3fffd) {
		return 2
	}
	return 1
}

func (e *evaluator) callMBStrimwidth(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "string"}, {name: "start"}, {name: "width"}, {name: "trim_marker"}, {name: "encoding"}}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 3 || len(args) > 5 || !present[0] || !present[1] || !present[2] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "mb_strimwidth() expects 3 to 5 arguments"}
	}
	if bound[0].kind != stringKind {
		return value{}, builtinTypeError(node.line, "mb_strimwidth", 1, "string", "string", bound[0])
	}
	start, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "mb_strimwidth", 2, "start", "int", bound[1])
	}
	width, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "mb_strimwidth", 3, "width", "int", bound[2])
	}
	marker := ""
	if present[3] {
		if bound[3].kind != stringKind {
			return value{}, builtinTypeError(node.line, "mb_strimwidth", 4, "trim_marker", "string", bound[3])
		}
		marker = bound[3].s
	}
	if present[4] && bound[4].kind != nullKind && (bound[4].kind != stringKind || !strings.EqualFold(bound[4].s, "UTF-8")) {
		return value{}, &runtimeError{node.line, "ValueError", "mb_strimwidth(): Argument #5 ($encoding) must be a valid encoding"}
	}
	if width.i < 0 {
		if err := e.diagnostic("Deprecated", node.line, "mb_strimwidth(): Passing a negative width is deprecated"); err != nil {
			return value{}, err
		}
	}
	runes := []rune(bound[0].s)
	i := int(start.i)
	if i < 0 {
		i += len(runes)
	}
	if i < 0 {
		i = 0
	}
	if i > len(runes) {
		i = len(runes)
	}
	runes = runes[i:]
	w := int(width.i)
	if w < 0 {
		w = len(runes) + w
	}
	if w < 0 {
		w = 0
	}
	total := 0
	for _, r := range runes {
		total += mbRuneWidth(r)
	}
	if total <= w {
		return value{kind: stringKind, s: string(runes)}, nil
	}
	for _, r := range marker {
		w -= mbRuneWidth(r)
	}
	if w < 0 {
		w = 0
	}
	out := make([]rune, 0, len(runes))
	for _, r := range runes {
		size := mbRuneWidth(r)
		if size > w {
			break
		}
		out = append(out, r)
		w -= size
	}
	return value{kind: stringKind, s: string(out) + marker}, nil
}
