package php

import "strings"

func (e *evaluator) callCaseString(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects exactly 1 argument"}
	}
	input, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, name, 1, "string", "string", bound[0])
	}
	bytes := []byte(input.s)
	for i, c := range bytes {
		if name == "strtolower" && c >= 'A' && c <= 'Z' {
			bytes[i] = c + ('a' - 'A')
		} else if name == "strtoupper" && c >= 'a' && c <= 'z' {
			bytes[i] = c - ('a' - 'A')
		}
	}
	return value{kind: stringKind, s: string(bytes)}, nil
}
