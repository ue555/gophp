package php

func (e *evaluator) callArrayMergeRecursive(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	out := arrayValue()
	for i, arg := range args {
		if arg.named {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + arg.name}
		}
		if arg.value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "array_merge_recursive", i+1, "arrays", "array", arg.value)
		}
		out, err = e.mergeRecursiveArrays(out, arg.value, node.line, 0)
		if err != nil {
			return value{}, err
		}
	}
	return out, nil
}

func (e *evaluator) mergeRecursiveArrays(left, right value, line, depth int) (value, error) {
	if depth > 64 {
		return value{}, &runtimeError{line, "Error", "gophp: recursive array merge depth exceeded"}
	}
	out := arrayValue()
	for _, entry := range left.a.entries {
		out.a = out.a.with(entry.key, entry.value)
	}
	for _, entry := range right.a.entries {
		if !entry.key.isString {
			key, err := out.a.appendKey(line)
			if err != nil {
				return value{}, err
			}
			out.a = out.a.with(key, entry.value)
			continue
		}
		old, ok := out.a.get(entry.key)
		if !ok {
			out.a = out.a.with(entry.key, entry.value)
			continue
		}
		old = dereference(old)
		incoming := dereference(entry.value)
		if old.kind != arrayKind {
			tmp := arrayValue()
			tmp.a = tmp.a.with(arrayKey{index: 0}, old)
			old = tmp
		}
		if incoming.kind != arrayKind {
			tmp := arrayValue()
			tmp.a = tmp.a.with(arrayKey{index: 0}, incoming)
			incoming = tmp
		}
		merged, err := e.mergeRecursiveArrays(old, incoming, line, depth+1)
		if err != nil {
			return value{}, err
		}
		out.a = out.a.with(entry.key, merged)
	}
	return out, nil
}
