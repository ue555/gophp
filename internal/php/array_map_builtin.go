package php

func (e *evaluator) callArrayMap(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_map() expects at least 2 arguments"}
	}
	callback := args[0].value
	var target resolvedCallable
	if callback.kind != nullKind {
		var ok bool
		target, ok = e.resolveCallable(callback, node.line)
		if !ok {
			return value{}, builtinTypeError(node.line, "array_map", 1, "callback", "?callable", callback)
		}
	}
	arrays := make([]*phpArray, len(args)-1)
	max := 0
	for i, arg := range args[1:] {
		if arg.value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "array_map", i+2, "array", "array", arg.value)
		}
		arrays[i] = arg.value.a
		if len(arrays[i].entries) > max {
			max = len(arrays[i].entries)
		}
	}
	out := arrayValue()
	for i := 0; i < max; i++ {
		params := make([]callArgument, len(arrays))
		for j, a := range arrays {
			v := value{kind: nullKind}
			if i < len(a.entries) {
				v = dereference(a.entries[i].value)
			}
			params[j] = callArgument{value: v, line: node.line}
		}
		v := value{kind: nullKind}
		if callback.kind == nullKind {
			if len(arrays) == 1 {
				v = params[0].value
			} else {
				v = arrayValue()
				for j, param := range params {
					v.a = v.a.with(arrayKey{index: int64(j)}, param.value)
				}
			}
		} else {
			v, err = e.invokeInternalCallback(target, params, node.line)
			if err != nil {
				return value{}, err
			}
		}
		key := arrayKey{index: int64(i)}
		if len(arrays) == 1 {
			key = arrays[0].entries[i].key
		}
		out.a = out.a.with(key, v)
	}
	return out, nil
}
