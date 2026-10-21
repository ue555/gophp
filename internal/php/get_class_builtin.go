package php

import "strings"

func (e *evaluator) callGetClass(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "get_class() expects exactly 1 argument"}
	}
	if bound[0].kind != objectKind {
		return value{}, builtinTypeError(node.line, "get_class", 1, "object", "object", bound[0])
	}
	return value{kind: stringKind, s: bound[0].o.className()}, nil
}

func (e *evaluator) callPropertyExists(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}, {name: "property"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "property_exists() expects exactly 2 arguments"}
	}
	nameValue, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "property_exists", 2, "property", "string", bound[1])
	}
	name := nameValue.s
	var class *classDecl
	if bound[0].kind == objectKind {
		class = bound[0].o.class
		if bound[0].o.properties != nil && bound[0].o.properties.find(arrayKey{text: name, isString: true}) >= 0 {
			return boolean(true), nil
		}
	} else if bound[0].kind == stringKind {
		e.initClasses()
		class = e.classes[strings.ToLower(bound[0].s)]
		if class == nil {
			if err := e.tryAutoload(bound[0].s, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[strings.ToLower(bound[0].s)]
		}
	} else {
		return value{}, builtinTypeError(node.line, "property_exists", 1, "object_or_class", "object|string", bound[0])
	}
	for c := class; c != nil; c = c.parent {
		for _, property := range c.properties {
			if property.name == name {
				return boolean(true), nil
			}
		}
	}
	return boolean(false), nil
}

func (e *evaluator) callMethodExists(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}, {name: "method"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "method_exists() expects exactly 2 arguments"}
	}
	name, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "method_exists", 2, "method", "string", bound[1])
	}
	var class *classDecl
	if bound[0].kind == objectKind {
		class = bound[0].o.class
	} else if bound[0].kind == stringKind {
		e.initClasses()
		class = e.classes[strings.ToLower(bound[0].s)]
		if class == nil {
			if err := e.tryAutoload(bound[0].s, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[strings.ToLower(bound[0].s)]
		}
	} else {
		return value{}, builtinTypeError(node.line, "method_exists", 1, "object_or_class", "object|string", bound[0])
	}
	if class == nil {
		return boolean(false), nil
	}
	return boolean(e.lookupMethod(class, name.s) != nil), nil
}

func (e *evaluator) callIsA(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}, {name: "class"}, {name: "allow_string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "is_a() expects 2 to 3 arguments"}
	}
	if bound[1].kind != stringKind {
		return value{}, builtinTypeError(node.line, "is_a", 2, "class", "string", bound[1])
	}
	targetName := strings.TrimPrefix(bound[1].s, "\\")
	var class *classDecl
	if bound[0].kind == objectKind {
		class = bound[0].o.class
	} else if bound[0].kind == stringKind && present[2] && bound[2].truthy() {
		e.initClasses()
		class = e.classes[strings.ToLower(strings.TrimPrefix(bound[0].s, "\\"))]
		if class == nil {
			if err := e.tryAutoload(bound[0].s, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[strings.ToLower(strings.TrimPrefix(bound[0].s, "\\"))]
		}
	} else {
		return boolean(false), nil
	}
	if class == nil {
		return boolean(false), nil
	}
	for c := class; c != nil; c = c.parent {
		if strings.EqualFold(c.name, targetName) {
			return boolean(true), nil
		}
		if classImplements(c, targetName) {
			return boolean(true), nil
		}
	}
	return boolean(false), nil
}
