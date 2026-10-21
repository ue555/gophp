package php

import "strings"

func (e *evaluator) callGetClassMethods(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "get_class_methods() expects exactly 1 argument"}
	}
	var class *classDecl
	if bound[0].kind == objectKind && bound[0].o != nil {
		class = bound[0].o.class
	} else if bound[0].kind == stringKind {
		e.initClasses()
		class = e.classes[strings.ToLower(strings.TrimPrefix(bound[0].s, "\\"))]
		if class == nil {
			if err := e.tryAutoload(bound[0].s, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[strings.ToLower(strings.TrimPrefix(bound[0].s, "\\"))]
		}
	} else {
		return value{}, builtinTypeError(node.line, "get_class_methods", 1, "object_or_class", "object|string", bound[0])
	}
	if class == nil {
		return boolean(false), nil
	}
	out := arrayValue()
	seen := map[string]bool{}
	for c := class; c != nil; c = c.parent {
		for _, method := range c.methodOrder {
			key := strings.ToLower(method.name)
			if seen[key] {
				continue
			}
			seen[key] = true
			if method.visibility == "private" && e.scope != c {
				continue
			}
			if method.visibility == "protected" && e.scope == nil {
				continue
			}
			out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, value{kind: stringKind, s: method.name})
		}
	}
	return out, nil
}
