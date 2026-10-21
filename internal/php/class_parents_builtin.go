package php

import "strings"

func (e *evaluator) callClassParents(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}, {name: "autoload"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "class_parents() expects 1 or 2 arguments"}
	}
	var class *classDecl
	if bound[0].kind == objectKind && bound[0].o != nil {
		class = bound[0].o.class
	} else if bound[0].kind == stringKind {
		e.initClasses()
		name := strings.TrimPrefix(bound[0].s, "\\")
		class = e.classes[strings.ToLower(name)]
		if class == nil && (!present[1] || bound[1].truthy()) {
			if err := e.tryAutoload(name, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[strings.ToLower(name)]
		}
	} else {
		return value{}, builtinTypeError(node.line, "class_parents", 1, "object_or_class", "object|string", bound[0])
	}
	if class == nil {
		return boolean(false), nil
	}
	out := arrayValue()
	for p := class.parent; p != nil; p = p.parent {
		v := value{kind: stringKind, s: p.name}
		out.a = out.a.with(arrayKey{text: p.name, isString: true}, v)
	}
	return out, nil
}

func (e *evaluator) callClassUses(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}, {name: "autoload"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "class_uses() expects 1 or 2 arguments"}
	}
	var class *classDecl
	if bound[0].kind == objectKind && bound[0].o != nil {
		class = bound[0].o.class
	} else if bound[0].kind == stringKind {
		e.initClasses()
		name := strings.TrimPrefix(bound[0].s, "\\")
		class = e.classes[strings.ToLower(name)]
		if class == nil && (!present[1] || bound[1].truthy()) {
			if err := e.tryAutoload(name, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[strings.ToLower(name)]
		}
	} else {
		return value{}, builtinTypeError(node.line, "class_uses", 1, "object_or_class", "object|string", bound[0])
	}
	if class == nil {
		return boolean(false), nil
	}
	out := arrayValue()
	for _, use := range class.traitUses {
		for _, name := range use.names {
			trait := e.classes[strings.ToLower(name)]
			if trait == nil {
				continue
			}
			v := value{kind: stringKind, s: trait.name}
			out.a = out.a.with(arrayKey{text: trait.name, isString: true}, v)
		}
	}
	return out, nil
}
