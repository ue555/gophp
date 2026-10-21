package php

import "strings"

func sameAutoloadTarget(a, b resolvedCallable) bool {
	if a.method != nil || b.method != nil {
		return a.method == b.method && a.receiver == b.receiver && a.called == b.called
	}
	if a.closure != nil || b.closure != nil {
		return a.closure == b.closure
	}
	return strings.EqualFold(a.function, b.function)
}

func (e *evaluator) callAutoloadBuiltin(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	name := strings.ToLower(node.name)
	params := []parameter{{name: "callback"}}
	if name == "spl_autoload_register" {
		params = append(params, parameter{name: "throw"}, parameter{name: "prepend"})
	}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > len(params) {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "(): Invalid argument count"}
	}
	target, ok := e.resolveCallable(bound[0], node.line)
	if !ok {
		return value{}, &runtimeError{node.line, "TypeError", name + "(): Argument #1 ($callback) must be a valid callback"}
	}
	if name == "spl_autoload_unregister" {
		for i, registered := range e.autoloaders {
			if sameAutoloadTarget(registered, target) {
				e.autoloaders = append(e.autoloaders[:i], e.autoloaders[i+1:]...)
				return boolean(true), nil
			}
		}
		return boolean(false), nil
	}
	for _, registered := range e.autoloaders {
		if sameAutoloadTarget(registered, target) {
			return boolean(true), nil
		}
	}
	if present[2] && bound[2].truthy() {
		e.autoloaders = append([]resolvedCallable{target}, e.autoloaders...)
	} else {
		e.autoloaders = append(e.autoloaders, target)
	}
	return boolean(true), nil
}

func (e *evaluator) tryAutoload(name string, line int) error {
	key := strings.ToLower(name)
	if len(e.autoloaders) == 0 || e.autoloading[key] {
		return nil
	}
	if e.autoloading == nil {
		e.autoloading = make(map[string]bool)
	}
	e.autoloading[key] = true
	defer delete(e.autoloading, key)
	for _, target := range append([]resolvedCallable(nil), e.autoloaders...) {
		_, err := e.invokeResolved(target, []callArgument{{value: value{kind: stringKind, s: name}, line: line}}, line)
		if err != nil {
			return err
		}
		if e.classes[key] != nil {
			break
		}
	}
	return nil
}
