package php

func (e *evaluator) callArraySearch(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "needle"}, {name: "haystack"}, {name: "strict"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_search() expects 2 to 3 arguments"}
	}
	if bound[1].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_search", 2, "haystack", "array", bound[1])
	}
	strict := present[2] && bound[2].truthy()
	for _, entry := range bound[1].a.entries {
		candidate := dereference(entry.value)
		matched := false
		if strict {
			matched = bound[0].identical(candidate)
		} else {
			result, err := e.comparisonValue("==", bound[0], candidate, node.line)
			if err != nil {
				return value{}, err
			}
			matched = result.truthy()
		}
		if matched {
			if entry.key.isString {
				return value{kind: stringKind, s: entry.key.text}, nil
			}
			return integer(entry.key.index), nil
		}
	}
	return boolean(false), nil
}
