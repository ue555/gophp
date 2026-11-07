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
		message := "class_implements(): Class " + bound[0].s + " does not exist"
		if !present[1] || bound[1].truthy() {
			message += " and could not be loaded"
		}
		if err := e.warn(node.line, message); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	out := arrayValue()
	for _, iface := range implementedInterfaces(class, make(map[*classDecl][]*classDecl)) {
		out.a = out.a.with(arrayKey{isString: true, text: iface.name}, value{kind: stringKind, s: iface.name})
	}
	return out, nil
}

// implementedInterfaces lists interfaces in the order PHP's class table holds
// them for user-declared classes: the parent's list first, then the interfaces
// declared on the class in source order, then each declared interface's own
// list appended in reverse (zend_do_inherit_interfaces walks it backwards),
// skipping interfaces already present. Engine-defined classes follow the same
// rule, which may differ from PHP's internal registration order for them.
func implementedInterfaces(decl *classDecl, memo map[*classDecl][]*classDecl) []*classDecl {
	if cached, ok := memo[decl]; ok {
		return cached
	}
	var out []*classDecl
	seen := make(map[*classDecl]bool)
	add := func(iface *classDecl) {
		if iface == nil || seen[iface] {
			return
		}
		seen[iface] = true
		out = append(out, iface)
	}
	if decl.parent != nil {
		for _, iface := range implementedInterfaces(decl.parent, memo) {
			add(iface)
		}
	}
	for _, iface := range decl.interfaces {
		add(iface)
	}
	for _, iface := range decl.interfaces {
		inherited := implementedInterfaces(iface, memo)
		for i := len(inherited) - 1; i >= 0; i-- {
			add(inherited[i])
		}
	}
	memo[decl] = out
	return out
}
