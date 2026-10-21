package php

func (e *evaluator) callUcfirst(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "ucfirst() expects exactly 1 argument"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "ucfirst", 1, "string", "string", bound[0])
	}
	if len(s.s) > 0 && s.s[0] >= 'a' && s.s[0] <= 'z' {
		s.s = string(s.s[0]-'a'+'A') + s.s[1:]
	}
	return s, nil
}
