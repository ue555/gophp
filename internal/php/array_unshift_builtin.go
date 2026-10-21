package php

func (e *evaluator) callArrayUnshift(node functionCall) (value, error) {
	if len(node.args) < 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_unshift() expects at least 2 arguments"}
	}
	target := node.args[0]
	if named, ok := target.(namedArgument); ok {
		if named.name != "array" {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + named.name}
		}
		target = named.operand
	}
	if !writable(target) {
		return value{}, &runtimeError{node.line, "Error", "array_unshift(): Argument #1 ($array) cannot be passed by reference"}
	}
	loc, err := e.locate(target)
	if err != nil {
		return value{}, err
	}
	array := loc.get()
	if array.kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_unshift", 1, "array", "array", array)
	}
	items, err := e.evalArguments(node.args[1:])
	if err != nil {
		return value{}, err
	}
	out := arrayValue()
	for _, item := range items {
		key, err := out.a.appendKey(node.line)
		if err != nil {
			return value{}, err
		}
		out.a = out.a.with(key, item.value)
	}
	for _, entry := range array.a.entries {
		key := entry.key
		if !key.isString {
			key, err = out.a.appendKey(node.line)
			if err != nil {
				return value{}, err
			}
		}
		out.a = out.a.with(key, entry.value)
	}
	if err := loc.set(out); err != nil {
		return value{}, err
	}
	return integer(int64(len(out.a.entries))), nil
}
