package php

func (e *evaluator) callCompact(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	out := arrayValue()
	var collect func(value, int) error
	collect = func(v value, depth int) error {
		if depth > 64 {
			return &runtimeError{node.line, "Error", "gophp: compact() nesting limit exceeded"}
		}
		v = dereference(v)
		if v.kind == arrayKind {
			for _, entry := range v.a.entries {
				if err := collect(entry.value, depth+1); err != nil {
					return err
				}
			}
			return nil
		}
		if v.kind != stringKind {
			return builtinTypeError(node.line, "compact", 1, "var_names", "array|string", v)
		}
		owner := e.variableOwner(v.s)
		if stored, ok := owner[v.s]; ok {
			out.a = out.a.with(arrayKey{text: v.s, isString: true}, dereference(stored))
		}
		return nil
	}
	for _, arg := range args {
		if arg.named {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + arg.name}
		}
		if err := collect(arg.value, 0); err != nil {
			return value{}, err
		}
	}
	return out, nil
}
