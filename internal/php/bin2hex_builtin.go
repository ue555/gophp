package php

import hexenc "encoding/hex"

func (e *evaluator) callBin2hex(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "bin2hex() expects exactly 1 argument"}
	}
	input, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "bin2hex", 1, "string", "string", bound[0])
	}
	return value{kind: stringKind, s: hexenc.EncodeToString([]byte(input.s))}, nil
}
