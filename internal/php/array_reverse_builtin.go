package php

func (e *evaluator) callArrayReverse(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "preserve_keys"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_reverse() expects 1 to 2 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_reverse", 1, "array", "array", bound[0])
	}
	preserve := present[1] && bound[1].truthy()
	out := arrayValue()
	next := int64(0)
	for i := len(bound[0].a.entries) - 1; i >= 0; i-- {
		entry := bound[0].a.entries[i]
		key := entry.key
		if !preserve && !key.isString {
			key = arrayKey{index: next}
			next++
		}
		out.a = out.a.with(key, entry.value)
	}
	return out, nil
}
