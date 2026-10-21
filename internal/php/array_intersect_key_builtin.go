package php

func (e *evaluator) callArrayIntersectKey(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_intersect_key() expects at least 2 arguments"}
	}
	for i, arg := range args {
		if arg.value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "array_intersect_key", i+1, "array", "array", arg.value)
		}
	}
	out := arrayValue()
	for _, entry := range args[0].value.a.entries {
		found := true
		for _, arg := range args[1:] {
			if arg.value.a.find(entry.key) < 0 {
				found = false
				break
			}
		}
		if found {
			out.a = out.a.with(entry.key, entry.value)
		}
	}
	return out, nil
}
