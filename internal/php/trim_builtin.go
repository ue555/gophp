package php

import (
	"fmt"
	"strings"
)

func (e *evaluator) callTrim(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "characters"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("%s() expects 1 to 2 arguments, %d given", name, len(args))}
	}
	stringValue, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, &runtimeError{node.line, "TypeError", fmt.Sprintf("%s(): Argument #1 ($string) must be of type string, %s given", name, bound[0].diagnosticName())}
	}
	chars := " \n\r\t\v\x00"
	if present[1] {
		charValue, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, &runtimeError{node.line, "TypeError", fmt.Sprintf("%s(): Argument #2 ($characters) must be of type string, %s given", name, bound[1].diagnosticName())}
		}
		chars = charValue.s
	}
	result := stringValue.s
	switch name {
	case "ltrim":
		result = strings.TrimLeft(result, chars)
	case "rtrim":
		result = strings.TrimRight(result, chars)
	default:
		result = strings.Trim(result, chars)
	}
	return value{kind: stringKind, s: result}, nil
}
