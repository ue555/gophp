package php

func (e *evaluator) callUnpack(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "format"}, {name: "string"}, {name: "offset"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "unpack() expects 2 or 3 arguments"}
	}
	format, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "unpack", 1, "format", "string", bound[0])
	}
	data, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "unpack", 2, "string", "string", bound[1])
	}
	offset := int64(0)
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "unpack", 3, "offset", "int", bound[2])
		}
		offset = v.i
	}
	if offset < 0 || offset > int64(len(data.s)) {
		return value{}, &runtimeError{node.line, "ValueError", "unpack(): Argument #3 ($offset) must be contained in argument #2 ($string)"}
	}
	if format.s != "C*" {
		return value{}, &runtimeError{node.line, "Error", "gophp: unpack format is not implemented: " + format.s}
	}
	out := arrayValue()
	for i, b := range []byte(data.s[offset:]) {
		out.a = out.a.with(arrayKey{index: int64(i + 1)}, integer(int64(b)))
	}
	return out, nil
}
