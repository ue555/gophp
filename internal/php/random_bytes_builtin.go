package php

import "crypto/rand"

func (e *evaluator) callRandomBytes(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "length"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "random_bytes() expects exactly 1 argument"}
	}
	n, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "random_bytes", 1, "length", "int", bound[0])
	}
	if n.i < 1 {
		return value{}, &runtimeError{node.line, "ValueError", "random_bytes(): Argument #1 ($length) must be greater than 0"}
	}
	if n.i > 1<<24 {
		return value{}, &runtimeError{node.line, "Error", "gophp: random_bytes() length exceeds implementation limit"}
	}
	bytes := make([]byte, n.i)
	if _, err := rand.Read(bytes); err != nil {
		return value{}, &runtimeError{node.line, "Random\\RandomException", "Random byte generation failed"}
	}
	return value{kind: stringKind, s: string(bytes)}, nil
}
