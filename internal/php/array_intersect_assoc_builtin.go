package php

func (e *evaluator) callArrayIntersectAssoc(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_intersect_assoc() expects at least 2 arguments"}
	}
	for i, arg := range args {
		if arg.value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "array_intersect_assoc", i+1, "array", "array", arg.value)
		}
	}
	out := arrayValue()
	for _, entry := range args[0].value.a.entries {
		left, err := e.toText(dereference(entry.value), node.line)
		if err != nil {
			return value{}, err
		}
		keep := true
		for _, arg := range args[1:] {
			other, found := arg.value.a.getRaw(entry.key)
			if !found {
				keep = false
				break
			}
			right, err := e.toText(dereference(other), node.line)
			if err != nil {
				return value{}, err
			}
			if left != right {
				keep = false
				break
			}
		}
		if keep {
			out.a = out.a.with(entry.key, entry.value)
		}
	}
	return out, nil
}
