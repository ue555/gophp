package php

func (e *evaluator) callGetCfgVar(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "option"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "get_cfg_var() expects exactly 1 argument"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "get_cfg_var", 1, "option", "string", bound[0])
	}
	switch s.s {
	case "include_path":
		return value{kind: stringKind, s: defaultIncludePath}, nil
	case "memory_limit":
		return value{kind: stringKind, s: "128M"}, nil
	case "default_charset":
		return value{kind: stringKind, s: "UTF-8"}, nil
	}
	return boolean(false), nil
}
