package php

func (e *evaluator) callIsCountable(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "is_countable() expects exactly 1 argument"}
	}
	v := bound[0]
	return boolean(v.kind == arrayKind || v.kind == objectKind && v.o != nil && classImplements(v.o.class, "Countable")), nil
}
