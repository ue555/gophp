package php

func (e *evaluator) callErrorReporting(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "error_level"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) > 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "error_reporting() expects at most 1 argument"}
	}
	previous := e.errorReporting
	if e.suppressedDiagnostics > 0 {
		previous &= 1 | 4 | 16 | 64 | 256 | 4096
	}
	if present[0] && bound[0].kind != nullKind {
		level, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "error_reporting", 1, "error_level", "?int", bound[0])
		}
		e.errorReporting = level.i
	}
	return integer(previous), nil
}
