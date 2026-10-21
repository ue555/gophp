package php

func (e *evaluator) callStrrev(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strrev() expects exactly 1 argument"}
	}
	input, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "strrev", 1, "string", "string", bound[0])
	}
	bytes := []byte(input.s)
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return value{kind: stringKind, s: string(bytes)}, nil
}
