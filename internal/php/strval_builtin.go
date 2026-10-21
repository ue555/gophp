package php

func (e *evaluator) callStrval(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strval() expects exactly 1 argument"}
	}
	s, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	return value{kind: stringKind, s: s}, nil
}
