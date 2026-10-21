package php

import "os"

func (e *evaluator) callGetenv(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "name"}, {name: "local_only"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "getenv() expects at most 2 arguments"}
	}
	if !present[0] {
		out := arrayValue()
		for _, item := range os.Environ() {
			for i := 0; i < len(item); i++ {
				if item[i] == '=' {
					out.a = out.a.with(arrayKey{isString: true, text: item[:i]}, value{kind: stringKind, s: item[i+1:]})
					break
				}
			}
		}
		return out, nil
	}
	name, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "getenv", 1, "name", "?string", bound[0])
	}
	v, found := os.LookupEnv(name.s)
	if !found {
		return boolean(false), nil
	}
	return value{kind: stringKind, s: v}, nil
}
