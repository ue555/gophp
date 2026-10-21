package php

import "strings"

func (e *evaluator) callIsSubclassOf(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}, {name: "class"}, {name: "allow_string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "is_subclass_of() expects 2 to 3 arguments"}
	}
	if bound[1].kind != stringKind {
		return value{}, builtinTypeError(node.line, "is_subclass_of", 2, "class", "string", bound[1])
	}
	targetName := strings.TrimPrefix(bound[1].s, "\\")
	var class *classDecl
	if bound[0].kind == objectKind && bound[0].o != nil {
		class = bound[0].o.class
	} else if bound[0].kind == stringKind && (!present[2] || bound[2].truthy()) {
		e.initClasses()
		class = e.classes[strings.ToLower(strings.TrimPrefix(bound[0].s, "\\"))]
		if class == nil {
			if err := e.tryAutoload(bound[0].s, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[strings.ToLower(strings.TrimPrefix(bound[0].s, "\\"))]
		}
	}
	if class == nil || strings.EqualFold(class.name, targetName) {
		return boolean(false), nil
	}
	for parent := class.parent; parent != nil; parent = parent.parent {
		if strings.EqualFold(parent.name, targetName) {
			return boolean(true), nil
		}
	}
	return boolean(classImplements(class, targetName)), nil
}
