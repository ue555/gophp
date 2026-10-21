package php

func (e *evaluator) callArraySplice(node functionCall) (value, error) {
	if len(node.args) < 2 || len(node.args) > 4 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_splice() expects 2 to 4 arguments"}
	}
	if !writable(node.args[0]) {
		return value{}, &runtimeError{node.line, "Error", "array_splice(): Argument #1 ($array) cannot be passed by reference"}
	}
	loc, err := e.locate(node.args[0])
	if err != nil {
		return value{}, err
	}
	array := loc.get()
	if array.kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_splice", 1, "array", "array", array)
	}
	args, err := e.evalArguments(node.args[1:])
	if err != nil {
		return value{}, err
	}
	offset, _, ok := e.coerceType(args[0].value, "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "array_splice", 2, "offset", "int", args[0].value)
	}
	n := int64(len(array.a.entries))
	start := offset.i
	if start < 0 {
		start += n
	}
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}
	end := n
	if len(args) >= 2 && args[1].value.kind != nullKind {
		length, _, ok := e.coerceType(args[1].value, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "array_splice", 3, "length", "?int", args[1].value)
		}
		if length.i >= 0 {
			end = start + length.i
		} else {
			end = n + length.i
		}
		if end < start {
			end = start
		}
		if end > n {
			end = n
		}
	}
	replacement := []value{}
	if len(args) == 3 {
		v := args[2].value
		if v.kind == arrayKind {
			for _, entry := range v.a.entries {
				replacement = append(replacement, dereference(entry.value))
			}
		} else if v.kind != nullKind {
			replacement = append(replacement, v)
		}
	}
	removed := arrayValue()
	for _, entry := range array.a.entries[start:end] {
		key := entry.key
		if !key.isString {
			key = arrayKey{index: int64(len(removed.a.entries))}
		}
		removed.a = removed.a.with(key, dereference(entry.value))
	}
	out := arrayValue()
	appendEntry := func(entry arrayEntry) error {
		key := entry.key
		if !key.isString {
			var keyErr error
			key, keyErr = out.a.appendKey(node.line)
			if keyErr != nil {
				return keyErr
			}
		}
		out.a = out.a.with(key, dereference(entry.value))
		return nil
	}
	for _, entry := range array.a.entries[:start] {
		if err := appendEntry(entry); err != nil {
			return value{}, err
		}
	}
	for _, v := range replacement {
		if err := appendEntry(arrayEntry{key: arrayKey{index: 0}, value: v}); err != nil {
			return value{}, err
		}
	}
	for _, entry := range array.a.entries[end:] {
		if err := appendEntry(entry); err != nil {
			return value{}, err
		}
	}
	if err := loc.set(out); err != nil {
		return value{}, err
	}
	return removed, nil
}
