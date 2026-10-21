package php

import "strings"

func (e *evaluator) callArrayWalk(node functionCall) (value, error) {
	recursive := strings.EqualFold(node.name, "array_walk_recursive")
	args, err := e.evalArgumentsCapturingArray(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "callback"}, {name: "arg"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || !present[1] || len(args) > 3 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects 2 to 3 arguments"}
	}
	arrayArgument := argumentForParameter(args, "array", 0)
	if arrayArgument == nil || arrayArgument.place == nil && arrayArgument.ref == nil {
		return value{}, &runtimeError{node.line, "Error", node.name + "(): Argument #1 ($array) cannot be passed by reference"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, node.name, 1, "array", "array", bound[0])
	}
	target, ok := e.resolveCallable(bound[1], node.line)
	if !ok {
		return value{}, builtinTypeError(node.line, node.name, 2, "callback", "callable", bound[1])
	}
	method := target.method
	if target.closure != nil {
		method = target.closure.method
	}
	var arrayPlace referencePlace
	if arrayArgument.place != nil {
		arrayPlace = *arrayArgument.place
	} else {
		cell := arrayArgument.ref
		arrayPlace = referencePlace{
			read:   func() value { return value{kind: refKind, ref: cell} },
			assign: func(v value) error { return e.writeReference(cell, v, node.line) },
		}
	}
	visiting := make(map[*phpArray]bool)
	var walk func(referencePlace, int) error
	walk = func(place referencePlace, depth int) error {
		// The language has no fixed 64-level limit. Keep a generous safety
		// bound for cycles whose copy-on-write path changes array identity.
		if depth > 1024 {
			return &runtimeError{node.line, "Error", "gophp: array_walk_recursive() nesting limit exceeded"}
		}
		current := dereference(place.read())
		if current.kind != arrayKind {
			return nil
		}
		if visiting[current.a] {
			return &runtimeError{node.line, "Error", "Recursion detected"}
		}
		visiting[current.a] = true
		defer delete(visiting, current.a)
		for _, entry := range current.a.entries {
			key := value{kind: stringKind, s: entry.key.text}
			if !entry.key.isString {
				key = integer(entry.key.index)
			}
			child, err := e.offsetReferencePlace(place, offsetExpr{key: literal{value: key, line: node.line}, line: node.line})
			if err != nil {
				return err
			}
			item := dereference(child.read())
			if recursive && item.kind == arrayKind {
				if err := walk(child, depth+1); err != nil {
					return err
				}
				continue
			}
			params := []callArgument{{value: item, line: node.line}, {value: key, line: node.line}}
			if method != nil && len(method.parameters) > 0 && method.parameters[0].byRef {
				params[0].place = &child
			}
			if present[2] {
				params = append(params, callArgument{value: bound[2], line: node.line})
			}
			if _, err := e.invokeInternalCallback(target, params, node.line); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk(arrayPlace, 0); err != nil {
		return value{}, err
	}
	return boolean(true), nil
}
