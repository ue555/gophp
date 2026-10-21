package php

func (e *evaluator) callArrayReplace(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_replace() expects at least 1 argument"}
	}
	out := arrayValue()
	for i, arg := range args {
		if arg.value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "array_replace", i+1, "array", "array", arg.value)
		}
		for _, entry := range arg.value.a.entries {
			out.a = out.a.with(entry.key, dereference(entry.value))
		}
	}
	return out, nil
}
