package php

import "strconv"

func (e *evaluator) callSetExceptionHandler(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "callback"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "set_exception_handler() expects exactly 1 argument"}
	}
	if bound[0].kind != nullKind {
		if _, ok := e.resolveCallable(bound[0], node.line); !ok {
			return value{}, builtinTypeError(node.line, "set_exception_handler", 1, "callback", "?callable", bound[0])
		}
	}
	previous := e.exceptionHandler
	e.exceptionHandlerStack = append(e.exceptionHandlerStack, previous)
	e.exceptionHandler = bound[0]
	if previous.kind == 0 {
		return value{kind: nullKind}, nil
	}
	return previous, nil
}

func (e *evaluator) callRestoreExceptionHandler(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) != 0 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "restore_exception_handler() expects exactly 0 arguments, " + strconv.Itoa(len(args)) + " given"}
	}
	if len(e.exceptionHandlerStack) == 0 {
		return boolean(true), nil
	}
	last := len(e.exceptionHandlerStack) - 1
	e.exceptionHandler = e.exceptionHandlerStack[last]
	e.exceptionHandlerStack = e.exceptionHandlerStack[:last]
	return boolean(true), nil
}
