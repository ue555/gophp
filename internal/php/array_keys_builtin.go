package php

func (e *evaluator) callArrayKeys(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "filter_value"}, {name: "strict"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 3 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_keys() expects 1 to 3 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_keys", 1, "array", "array", bound[0])
	}
	strict := present[2] && bound[2].truthy()
	out := arrayValue()
	for _, entry := range bound[0].a.entries {
		if present[1] {
			match := false
			if strict {
				match = dereference(entry.value).identical(bound[1])
			} else {
				v, err := e.comparisonValue("==", dereference(entry.value), bound[1], node.line)
				if err != nil {
					return value{}, err
				}
				match = v.truthy()
			}
			if !match {
				continue
			}
		}
		key := value{kind: stringKind, s: entry.key.text}
		if !entry.key.isString {
			key = integer(entry.key.index)
		}
		out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, key)
	}
	return out, nil
}
