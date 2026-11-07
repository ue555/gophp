package php

func (e *evaluator) callPHPVersion(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "extension"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) > 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "phpversion() expects at most 1 argument"}
	}
	if present[0] && bound[0].kind != nullKind {
		if _, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes); !ok {
			return value{}, builtinTypeError(node.line, "phpversion", 1, "extension", "?string", bound[0])
		}
		return boolean(false), nil
	}
	return value{kind: stringKind, s: e.version.Patch()}, nil
}
