package php

// array_udiff_assoc compares corresponding keys with user-supplied value
// comparison. Keep the first array's insertion order and original values.
func (e *evaluator) callArrayUDiffAssoc(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 3 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_udiff_assoc() expects at least 3 arguments"}
	}
	for i, arg := range args[:len(args)-1] {
		if arg.value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "array_udiff_assoc", i+1, "array", "array", arg.value)
		}
	}
	target, ok := e.resolveCallable(args[len(args)-1].value, node.line)
	if !ok {
		return value{}, builtinTypeError(node.line, "array_udiff_assoc", len(args), "value_compare_func", "callable", args[len(args)-1].value)
	}
	out := arrayValue()
	for _, entry := range args[0].value.a.entries {
		matched := false
		for _, other := range args[1 : len(args)-1] {
			candidate, found := other.value.a.getRaw(entry.key)
			if !found {
				continue
			}
			comparison, callErr := e.invokeInternalCallback(target, []callArgument{{value: dereference(entry.value), line: node.line}, {value: dereference(candidate), line: node.line}}, node.line)
			if callErr != nil {
				return value{}, callErr
			}
			if comparison.i == 0 && comparison.kind == intKind {
				matched = true
				break
			}
		}
		if !matched {
			out.a = out.a.with(entry.key, entry.value)
		}
	}
	return out, nil
}
