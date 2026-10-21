package php

import "fmt"

func (e *evaluator) callExit(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if e.version == PHP83 {
		for _, arg := range args {
			if arg.named {
				return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + arg.name}
			}
		}
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "status"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) > 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "exit() expects at most 1 argument"}
	}
	if !present[0] {
		return value{}, &exitSignal{code: 0}
	}
	status := bound[0]
	if status.kind == intKind {
		return value{}, &exitSignal{code: int(status.i) & 255}
	}
	if e.version != PHP83 {
		if status.kind == boolKind || status.kind == floatKind {
			if e.strictTypes {
				return value{}, builtinTypeError(node.line, "exit", 1, "status", "string|int", status)
			}
			return value{}, &exitSignal{code: int(status.asFloat()) & 255}
		}
		if status.kind == arrayKind || status.kind == objectKind {
			return value{}, builtinTypeError(node.line, "exit", 1, "status", "string|int", status)
		}
	}
	message, err := e.toText(status, node.line)
	if err != nil {
		return value{}, err
	}
	fmt.Fprint(e.stdout, message)
	return value{}, &exitSignal{code: 0}
}
