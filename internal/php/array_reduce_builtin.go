package php

func (e *evaluator) callArrayReduce(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "callback"}, {name: "initial"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_reduce() expects 2 to 3 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_reduce", 1, "array", "array", bound[0])
	}
	target, ok := e.resolveCallable(bound[1], node.line)
	if !ok {
		return value{}, builtinTypeError(node.line, "array_reduce", 2, "callback", "callable", bound[1])
	}
	carry := value{kind: nullKind}
	if present[2] {
		carry = bound[2]
	}
	for _, entry := range bound[0].a.entries {
		carry, err = e.invokeInternalCallback(target, []callArgument{{value: carry}, {value: dereference(entry.value)}}, node.line)
		if err != nil {
			return value{}, err
		}
	}
	return carry, nil
}
