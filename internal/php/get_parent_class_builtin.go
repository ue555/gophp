package php

import "strings"

func (e *evaluator) callGetParentClass(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) > 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "get_parent_class() expects at most 1 argument"}
	}
	var class *classDecl
	if !present[0] {
		class = e.scope
	} else if bound[0].kind == objectKind && bound[0].o != nil {
		class = bound[0].o.class
	} else if bound[0].kind == stringKind {
		e.initClasses()
		name := strings.TrimPrefix(bound[0].s, "\\")
		class = e.classes[strings.ToLower(name)]
		if class == nil {
			if err := e.tryAutoload(name, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[strings.ToLower(name)]
		}
	} else {
		return value{}, builtinTypeError(node.line, "get_parent_class", 1, "object_or_class", "object|string", bound[0])
	}
	if class == nil || class.parent == nil {
		return boolean(false), nil
	}
	return value{kind: stringKind, s: class.parent.name}, nil
}
