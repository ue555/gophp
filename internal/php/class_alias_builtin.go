package php

import "strings"

func (e *evaluator) callClassAlias(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "class"}, {name: "alias"}, {name: "autoload"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "class_alias() expects 2 to 3 arguments"}
	}
	original, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "class_alias", 1, "class", "string", bound[0])
	}
	alias, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "class_alias", 2, "alias", "string", bound[1])
	}
	e.initClasses()
	key := strings.ToLower(strings.TrimPrefix(original.s, "\\"))
	class := e.classes[key]
	if class == nil && (!present[2] || bound[2].truthy()) {
		class, _ = e.resolveClassName(original.s, node.line)
	}
	if class == nil || class.kind == "trait" {
		return boolean(false), nil
	}
	aliasKey := strings.ToLower(strings.TrimPrefix(alias.s, "\\"))
	if e.classes[aliasKey] != nil {
		if err := e.warn(node.line, "Cannot declare class "+alias.s+", because the name is already in use"); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	e.classes[aliasKey] = class
	return boolean(true), nil
}
