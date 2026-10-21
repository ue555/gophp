package php

type shutdownCall struct {
	target resolvedCallable
	args   []callArgument
}

func (e *evaluator) callRegisterShutdownFunction(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "register_shutdown_function() expects at least 1 argument"}
	}
	target, ok := e.resolveCallable(args[0].value, node.line)
	if !ok {
		return value{}, builtinTypeError(node.line, "register_shutdown_function", 1, "callback", "callable", args[0].value)
	}
	e.shutdownFunctions = append(e.shutdownFunctions, shutdownCall{target: target, args: append([]callArgument(nil), args[1:]...)})
	return value{kind: nullKind}, nil
}

func (e *evaluator) runShutdownFunctions() error {
	for i := 0; i < len(e.shutdownFunctions); i++ {
		call := e.shutdownFunctions[i]
		if _, err := e.invokeResolved(call.target, call.args, 0); err != nil {
			return err
		}
	}
	return nil
}
