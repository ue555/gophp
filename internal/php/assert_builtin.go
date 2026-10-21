package php

func (e *evaluator) callAssert(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "assertion"}, {name: "description"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "assert() expects 1 to 2 arguments"}
	}
	if bound[0].truthy() {
		return boolean(true), nil
	}
	message := "assertion failed"
	if present[1] {
		s, err := e.toText(bound[1], node.line)
		if err != nil {
			return value{}, err
		}
		message = s
	}
	return value{}, &runtimeError{node.line, "AssertionError", message}
}
