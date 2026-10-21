package php

import "strings"

func (e *evaluator) callClassImplements(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object_or_class"}, {name: "autoload"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "class_implements() expects 1 or 2 arguments"}
	}
	e.initClasses()
	var class *classDecl
	if bound[0].kind == objectKind && bound[0].o != nil {
		class = bound[0].o.class
	} else if bound[0].kind == stringKind {
		name := strings.TrimPrefix(bound[0].s, "\\")
		class = e.classes[symbolKey(name)]
		if class == nil && (!present[1] || bound[1].truthy()) {
			if err := e.tryAutoload(name, node.line); err != nil {
				return value{}, err
			}
			class = e.classes[symbolKey(name)]
		}
	} else {
		return value{}, builtinTypeError(node.line, "class_implements", 1, "object_or_class", "object|string", bound[0])
	}
	if class == nil {
		if err := e.warn(node.line, "class_implements(): Class "+bound[0].s+" does not exist and could not be loaded"); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	out := arrayValue()
	seen := make(map[*classDecl]bool)
	var visitInterface func(*classDecl)
	visitInterface = func(iface *classDecl) {
		if iface == nil || seen[iface] {
			return
		}
		seen[iface] = true
		out.a = out.a.with(arrayKey{isString: true, text: iface.name}, value{kind: stringKind, s: iface.name})
		for _, parent := range iface.interfaces {
			visitInterface(parent)
		}
	}
	for current := class; current != nil; current = current.parent {
		for _, iface := range current.interfaces {
			visitInterface(iface)
		}
	}
	return out, nil
}
