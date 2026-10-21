package php

func (e *evaluator) callArrayPad(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "length"}, {name: "value"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 3 || !present[0] || !present[1] || !present[2] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_pad() expects exactly 3 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_pad", 1, "array", "array", bound[0])
	}
	length, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "array_pad", 2, "length", "int", bound[1])
	}
	target := length.i
	if target < 0 {
		target = -target
	}
	if target > 1048576 {
		return value{}, &runtimeError{node.line, "ValueError", "array_pad(): Argument #2 ($length) is too large"}
	}
	count := int64(len(bound[0].a.entries))
	if target <= count {
		return bound[0], nil
	}
	out := arrayValue()
	if length.i < 0 {
		for i := count; i < target; i++ {
			out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, bound[2])
		}
	}
	for _, entry := range bound[0].a.entries {
		if length.i < 0 && !entry.key.isString {
			out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, entry.value)
		} else {
			out.a = out.a.with(entry.key, entry.value)
		}
	}
	if length.i > 0 {
		for i := count; i < target; i++ {
			key, err := out.a.appendKey(node.line)
			if err != nil {
				return value{}, err
			}
			out.a = out.a.with(key, bound[2])
		}
	}
	return out, nil
}
