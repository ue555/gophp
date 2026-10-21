package php

import "strings"

func (e *evaluator) callStrtok(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strtok() expects 1 to 2 arguments"}
	}
	var text, token value
	if len(args) == 2 {
		text, _, _ = e.coerceType(args[0].value, "string", false, e.strictTypes)
		token, _, _ = e.coerceType(args[1].value, "string", false, e.strictTypes)
		e.strtokRest = text.s
	} else {
		token, _, _ = e.coerceType(args[0].value, "string", false, e.strictTypes)
	}
	if token.s == "" {
		return value{}, &runtimeError{node.line, "ValueError", "strtok(): Argument #2 ($token) must not be empty"}
	}
	remaining := strings.TrimLeft(e.strtokRest, token.s)
	if remaining == "" {
		e.strtokRest = ""
		return boolean(false), nil
	}
	end := strings.IndexAny(remaining, token.s)
	if end < 0 {
		e.strtokRest = ""
		return value{kind: stringKind, s: remaining}, nil
	}
	e.strtokRest = remaining[end:]
	return value{kind: stringKind, s: remaining[:end]}, nil
}
