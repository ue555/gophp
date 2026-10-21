package php

func (e *evaluator) callArrayReplaceRecursive(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_replace_recursive() expects at least 1 argument"}
	}
	out := arrayValue()
	for i, arg := range args {
		if arg.named {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + arg.name}
		}
		if arg.value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "array_replace_recursive", i+1, "array", "array", arg.value)
		}
		out, err = replaceRecursive(out, arg.value, node.line, 0)
		if err != nil {
			return value{}, err
		}
	}
	return out, nil
}
func replaceRecursive(left, right value, line, depth int) (value, error) {
	if depth > 64 {
		return value{}, &runtimeError{line, "Error", "gophp: recursive array replace depth exceeded"}
	}
	out := arrayValue()
	for _, entry := range left.a.entries {
		out.a = out.a.with(entry.key, entry.value)
	}
	for _, entry := range right.a.entries {
		old, found := out.a.get(entry.key)
		newValue := dereference(entry.value)
		if found && old.kind == arrayKind && newValue.kind == arrayKind {
			merged, err := replaceRecursive(old, newValue, line, depth+1)
			if err != nil {
				return value{}, err
			}
			newValue = merged
		}
		out.a = out.a.with(entry.key, newValue)
	}
	return out, nil
}
