package php

func (e *evaluator) callArrayIntersect(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_intersect() expects at least 2 arguments"}
	}
	for i, arg := range args {
		if arg.value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "array_intersect", i+1, "array", "array", arg.value)
		}
	}
	sets := make([]map[string]bool, 0, len(args)-1)
	for _, arg := range args[1:] {
		set := make(map[string]bool)
		for _, entry := range arg.value.a.entries {
			s, err := e.toText(dereference(entry.value), node.line)
			if err != nil {
				return value{}, err
			}
			set[s] = true
		}
		sets = append(sets, set)
	}
	out := arrayValue()
	for _, entry := range args[0].value.a.entries {
		s, err := e.toText(dereference(entry.value), node.line)
		if err != nil {
			return value{}, err
		}
		keep := true
		for _, set := range sets {
			if !set[s] {
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
