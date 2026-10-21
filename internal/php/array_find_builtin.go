package php

func (e *evaluator) callArrayFindKey(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "callback"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_find_key() expects exactly 2 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_find_key", 1, "array", "array", bound[0])
	}
	callback, ok := e.resolveCallable(bound[1], node.line)
	if !ok {
		return value{}, &runtimeError{node.line, "TypeError", "array_find_key(): Argument #2 ($callback) must be a valid callback"}
	}
	for _, entry := range bound[0].a.entries {
		key := integer(entry.key.index)
		if entry.key.isString {
			key = value{kind: stringKind, s: entry.key.text}
		}
		result, err := e.invokeInternalCallback(callback, []callArgument{{value: dereference(entry.value)}, {value: key}}, node.line)
		if err != nil {
			return value{}, err
		}
		if result.truthy() {
			return key, nil
		}
	}
	return value{kind: nullKind}, nil
}
