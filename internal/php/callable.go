package php

import (
	"fmt"
	"strings"
)

// resolvedCallable is an executable target captured at the point where a
// Closure is made. In particular, method visibility is checked once, while
// the scope which was allowed to access it still exists.
type resolvedCallable struct {
	function string
	method   *methodDecl
	receiver value
	called   *classDecl
	closure  *phpClosure
}

func (e *evaluator) resolveCallable(v value, line int) (resolvedCallable, bool) {
	if v.kind == objectKind && v.o != nil && v.o.closure != nil {
		return resolvedCallable{closure: v.o.closure}, true
	}
	if v.kind == objectKind && v.o != nil {
		return e.resolveCallableMethod(v, "__invoke", line)
	}
	if v.kind == stringKind {
		name := strings.TrimPrefix(v.s, "\\")
		if parts := strings.SplitN(name, "::", 2); len(parts) == 2 {
			return e.resolveCallableMethod(value{kind: stringKind, s: parts[0]}, parts[1], line)
		}
		if e.functions[symbolKey(name)] != nil || builtinNameForVersion(name, e.version) {
			return resolvedCallable{function: name}, true
		}
		return resolvedCallable{}, false
	}
	if v.kind != arrayKind || v.a == nil || len(v.a.entries) != 2 {
		return resolvedCallable{}, false
	}
	first, second := v.a.entries[0], v.a.entries[1]
	if first.key.isString || first.key.index != 0 || second.key.isString || second.key.index != 1 || second.value.kind != stringKind {
		return resolvedCallable{}, false
	}
	return e.resolveCallableMethod(first.value, second.value.s, line)
}

func (e *evaluator) resolveCallableMethod(receiver value, name string, line int) (resolvedCallable, bool) {
	var class *classDecl
	if receiver.kind == objectKind && receiver.o != nil {
		class = receiver.o.class
	} else if receiver.kind == stringKind {
		e.initClasses()
		className := strings.TrimPrefix(receiver.s, "\\")
		switch strings.ToLower(className) {
		case "self":
			class = e.scope
		case "parent":
			if e.scope != nil {
				class = e.scope.parent
			}
		case "static":
			class = e.lateClass
		default:
			class = e.classes[symbolKey(className)]
		}
	}
	if class == nil {
		return resolvedCallable{}, false
	}
	method := e.lookupMethod(class, name)
	if method == nil || e.methodAccess(class, method, line) != nil {
		return resolvedCallable{}, false
	}
	if receiver.kind == stringKind && !method.static {
		return resolvedCallable{}, false
	}
	called := class
	if receiver.kind == stringKind && (strings.EqualFold(receiver.s, "self") || strings.EqualFold(receiver.s, "parent") || strings.EqualFold(receiver.s, "static")) && e.lateClass != nil {
		called = e.lateClass
	}
	if method.static {
		receiver = value{}
	}
	return resolvedCallable{method: method, receiver: receiver, called: called}, true
}

func (e *evaluator) invokeResolved(target resolvedCallable, args []callArgument, line int) (value, error) {
	if target.closure != nil {
		return e.invokeClosureValues(target.closure, args, line)
	}
	if target.method != nil {
		return e.invokeWithClass(target.receiver, target.called, target.method, args, line)
	}
	if builtinNameForVersion(target.function, e.version) {
		expressions := make([]expr, len(args))
		for i, arg := range args {
			item := expr(literal{value: arg.value, line: line})
			if arg.named {
				item = namedArgument{name: arg.name, operand: item, line: line}
			}
			expressions[i] = item
		}
		return e.callBuiltin(functionCall{name: target.function, args: expressions, line: line})
	}
	return e.invoke(value{}, e.functions[symbolKey(target.function)], args, line)
}

// Calls made by a PHP internal function use coercive argument checks even
// when the userland call to that internal function came from a strict file.
func (e *evaluator) invokeInternalCallback(target resolvedCallable, args []callArgument, line int) (value, error) {
	callerStrict := e.strictTypes
	e.strictTypes = false
	defer func() { e.strictTypes = callerStrict }()
	return e.invokeResolved(target, args, line)
}

func (e *evaluator) callableClosure(target resolvedCallable) value {
	e.initClasses()
	return value{kind: objectKind, o: &phpObject{class: e.classes["closure"], properties: &phpArray{}, closure: &phpClosure{target: &target}}}
}

func (e *evaluator) evalFirstClass(node firstClassCallable) (value, error) {
	var target resolvedCallable
	var ok bool
	switch callee := node.callee.(type) {
	case functionCall:
		name := callee.name
		if callee.fallback != "" && e.functions[symbolKey(name)] == nil && !builtinNameForVersion(name, e.version) {
			name = callee.fallback
		}
		target, ok = e.resolveCallable(value{kind: stringKind, s: name}, node.line)
	case methodCall:
		receiver, err := e.eval(callee.receiver)
		if err != nil {
			return value{}, err
		}
		target, ok = e.resolveCallableMethod(receiver, callee.name, node.line)
		if !ok && receiver.kind == objectKind && receiver.o != nil {
			if method := e.lookupMethod(receiver.o.class, callee.name); method != nil {
				if err := e.methodAccess(receiver.o.class, method, node.line); err != nil {
					return value{}, err
				}
			}
		}
	case staticCallExpr:
		class, err := e.resolveStaticClass(callee.className, callee.classExpr, node.line)
		if err != nil {
			return value{}, err
		}
		name, err := e.staticMemberName(callee.name, callee.nameExpr, node.line)
		if err != nil {
			return value{}, err
		}
		target, ok = e.resolveCallableMethod(value{kind: stringKind, s: class.name}, name, node.line)
		if ok && (strings.EqualFold(callee.className, "self") || strings.EqualFold(callee.className, "parent") || strings.EqualFold(callee.className, "static")) && e.lateClass != nil {
			target.called = e.lateClass
		}
	default:
		v, err := e.eval(node.callee)
		if err != nil {
			return value{}, err
		}
		target, ok = e.resolveCallable(v, node.line)
	}
	if !ok {
		return value{}, &runtimeError{node.line, "Error", "Failed to create closure from callable"}
	}
	return e.callableClosure(target), nil
}

func (e *evaluator) closureFromCallable(expressions []expr, line int) (value, error) {
	args, err := e.evalArguments(expressions)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "callback"}}, line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{line, "ArgumentCountError", fmt.Sprintf("Closure::fromCallable() expects exactly 1 argument, %d given", len(args))}
	}
	if bound[0].kind == objectKind && bound[0].o != nil && bound[0].o.closure != nil {
		return bound[0], nil
	}
	target, ok := e.resolveCallable(bound[0], line)
	if !ok {
		return value{}, &runtimeError{line, "TypeError", "Closure::fromCallable(): Argument #1 ($callback) must be a valid callback"}
	}
	return e.callableClosure(target), nil
}

func (e *evaluator) callCallableBuiltin(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	name := strings.ToLower(node.name)
	switch name {
	case "is_callable":
		if len(args) > 2 {
			return value{}, &runtimeError{node.line, "Error", "is_callable(): Argument #3 ($callable_name) requires a reference"}
		}
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}, {name: "syntax_only"}}, node.line)
		if err != nil {
			return value{}, err
		}
		if !present[0] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "is_callable() expects at least 1 argument, 0 given"}
		}
		if present[1] && bound[1].truthy() {
			return boolean(callableSyntax(bound[0])), nil
		}
		_, ok := e.resolveCallable(bound[0], node.line)
		return boolean(ok), nil
	case "call_user_func", "call_user_func_array":
		parameters := []parameter{{name: "callback"}}
		if name == "call_user_func" {
			parameters = append(parameters, parameter{name: "args", variadic: true})
		} else {
			parameters = append(parameters, parameter{name: "args"})
		}
		bound, present, extra, err := bindCallArguments(args, parameters, node.line)
		if err != nil {
			return value{}, err
		}
		if !present[0] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("%s() expects at least 1 argument, %d given", name, len(args))}
		}
		if name == "call_user_func_array" && (len(args) != 2 || !present[1]) {
			return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("call_user_func_array() expects exactly 2 arguments, %d given", len(args))}
		}
		target, ok := e.resolveCallable(bound[0], node.line)
		if !ok {
			return value{}, &runtimeError{node.line, "TypeError", name + "(): Argument #1 ($callback) must be a valid callback"}
		}
		if name == "call_user_func" {
			return e.invokeResolved(target, callArgumentsFromArray(extra.a, node.line), node.line)
		}
		if bound[1].kind != arrayKind {
			return value{}, &runtimeError{node.line, "TypeError", "call_user_func_array(): Argument #2 ($args) must be of type array"}
		}
		return e.invokeResolved(target, callArgumentsFromArray(bound[1].a, node.line), node.line)
	}
	return value{}, &runtimeError{node.line, "Error", "Unsupported callable builtin"}
}

func callArgumentsFromArray(array *phpArray, line int) []callArgument {
	args := make([]callArgument, 0, len(array.entries))
	for _, entry := range array.entries {
		args = append(args, callArgument{name: entry.key.text, named: entry.key.isString, value: dereference(entry.value), line: line})
	}
	return args
}

func callableSyntax(v value) bool {
	if v.kind == stringKind {
		return v.s != ""
	}
	if v.kind == objectKind && v.o != nil && v.o.closure != nil {
		return true
	}
	if v.kind != arrayKind || v.a == nil || len(v.a.entries) != 2 {
		return false
	}
	a, b := v.a.entries[0], v.a.entries[1]
	return !a.key.isString && a.key.index == 0 && !b.key.isString && b.key.index == 1 &&
		(a.value.kind == stringKind || a.value.kind == objectKind) && b.value.kind == stringKind
}
