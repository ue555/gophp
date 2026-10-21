package php

func (e *evaluator) callArrayFlip(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_flip() expects exactly 1 argument"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_flip", 1, "array", "array", bound[0])
	}
	out := arrayValue()
	for _, entry := range bound[0].a.entries {
		v := dereference(entry.value)
		if v.kind != intKind && v.kind != stringKind {
			if err := e.warn(node.line, "Can only flip string and integer values, entry skipped"); err != nil {
				return value{}, err
			}
			continue
		}
		key, err := e.arrayKey(v, node.line)
		if err != nil {
			return value{}, err
		}
		oldKey := value{kind: stringKind, s: entry.key.text}
		if !entry.key.isString {
			oldKey = integer(entry.key.index)
		}
		out.a = out.a.with(key, oldKey)
	}
	return out, nil
}
