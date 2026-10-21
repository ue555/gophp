package php

import (
	"crypto/sha1"
	"fmt"
)

func (e *evaluator) callStringPieces(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if node.name == "sha1" {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "binary"}}, node.line)
		if err != nil {
			return value{}, err
		}
		if !present[0] || len(args) > 2 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "sha1() expects 1 to 2 arguments"}
		}
		s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "sha1", 1, "string", "string", bound[0])
		}
		hash := sha1.Sum([]byte(s.s))
		if present[1] && bound[1].truthy() {
			return value{kind: stringKind, s: string(hash[:])}, nil
		}
		return value{kind: stringKind, s: fmt.Sprintf("%x", hash)}, nil
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "length"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "str_split() expects 1 to 2 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "str_split", 1, "string", "string", bound[0])
	}
	length := int64(1)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "str_split", 2, "length", "int", bound[1])
		}
		length = v.i
	}
	if length < 1 {
		return value{}, &runtimeError{node.line, "ValueError", "str_split(): Argument #2 ($length) must be greater than 0"}
	}
	out := arrayValue()
	for start := 0; start < len(s.s); start += int(length) {
		end := start + int(length)
		if end > len(s.s) {
			end = len(s.s)
		}
		out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, value{kind: stringKind, s: s.s[start:end]})
	}
	return out, nil
}
