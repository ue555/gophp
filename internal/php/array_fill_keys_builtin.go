package php

func (e *evaluator) callArrayFillKeys(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "keys"}, {name: "value"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_fill_keys() expects exactly 2 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_fill_keys", 1, "keys", "array", bound[0])
	}
	out := arrayValue()
	for _, entry := range bound[0].a.entries {
		key, err := e.arrayKey(dereference(entry.value), node.line)
		if err != nil {
			return value{}, err
		}
		out.a = out.a.with(key, bound[1])
	}
	return out, nil
}
