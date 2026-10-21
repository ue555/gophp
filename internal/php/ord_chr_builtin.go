package php

import "strings"

func (e *evaluator) callOrdChr(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	name := strings.ToLower(node.name)
	param := "character"
	if name == "chr" {
		param = "codepoint"
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: param}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects exactly 1 argument"}
	}
	if name == "ord" {
		s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "ord", 1, "character", "string", bound[0])
		}
		if s.s == "" {
			return value{}, &runtimeError{node.line, "ValueError", "ord(): Argument #1 ($character) must be a character"}
		}
		return integer(int64(s.s[0])), nil
	}
	n, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "chr", 1, "codepoint", "int", bound[0])
	}
	return value{kind: stringKind, s: string([]byte{byte(n.i)})}, nil
}
