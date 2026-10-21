package php

import "strings"

func (e *evaluator) callBasename(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "path"}, {name: "suffix"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "basename() expects 1 or 2 arguments"}
	}
	path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "basename", 1, "path", "string", bound[0])
	}
	text := strings.TrimRight(path.s, "/")
	if slash := strings.LastIndexByte(text, '/'); slash >= 0 {
		text = text[slash+1:]
	}
	if present[1] {
		suffix, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "basename", 2, "suffix", "string", bound[1])
		}
		if suffix.s != "" && strings.HasSuffix(text, suffix.s) {
			text = strings.TrimSuffix(text, suffix.s)
		}
	}
	return value{kind: stringKind, s: text}, nil
}
