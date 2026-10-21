package php

func (e *evaluator) callArrayCombine(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "keys"}, {name: "values"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_combine() expects exactly 2 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_combine", 1, "keys", "array", bound[0])
	}
	if bound[1].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_combine", 2, "values", "array", bound[1])
	}
	if len(bound[0].a.entries) != len(bound[1].a.entries) {
		return value{}, &runtimeError{node.line, "ValueError", "array_combine(): Argument #1 ($keys) and argument #2 ($values) must have the same number of elements"}
	}
	out := arrayValue()
	for i, keyEntry := range bound[0].a.entries {
		key, err := e.arrayKey(dereference(keyEntry.value), node.line)
		if err != nil {
			return value{}, err
		}
		out.a = out.a.with(key, dereference(bound[1].a.entries[i].value))
	}
	return out, nil
}
