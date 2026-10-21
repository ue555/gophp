package php

import "strings"

type phpClosure struct {
	method    *methodDecl
	target    *resolvedCallable
	captures  variableScope
	boundThis value
	called    *classDecl
	static    bool
}

func (p *parser) anonymousFunction(line int, static bool) (expr, error) {
	byRefReturn := p.accept("&")
	parameters, err := p.parameterList(false)
	if err != nil {
		return nil, err
	}
	var captures []closureCapture
	if p.acceptKeyword("use") {
		if err := p.expect("("); err != nil {
			return nil, err
		}
		seen := make(map[string]bool)
		for _, parameter := range parameters {
			seen[parameter.name] = true
		}
		for p.current().kind != ")" {
			byRef := p.accept("&")
			name := p.current()
			if err := p.expect("variable"); err != nil {
				return nil, err
			}
			if name.text == "this" || seen[name.text] {
				return nil, syntax(name.line, "invalid closure capture $%s", name.text)
			}
			seen[name.text] = true
			captures = append(captures, closureCapture{name: name.text, byRef: byRef})
			if !p.accept(",") {
				break
			}
		}
		if err := p.expect(")"); err != nil {
			return nil, err
		}
	}
	returnType, nullable, err := p.returnDeclaration()
	if err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	method := &methodDecl{name: "{closure}", parameters: parameters, returnType: returnType, returnNullable: nullable, body: body, line: line, source: sourceContext{p.filename, p.namespace, p.strictTypes}, version: p.version, returnsByReference: byRefReturn}
	return closureExpr{method: method, captures: captures, static: static, line: line}, nil
}

func (p *parser) arrowFunction(line int, static bool, end ...string) (expr, error) {
	byRefReturn := p.accept("&")
	parameters, err := p.parameterList(false)
	if err != nil {
		return nil, err
	}
	returnType, nullable, err := p.returnDeclaration()
	if err != nil {
		return nil, err
	}
	if err := p.expect("=>"); err != nil {
		return nil, err
	}
	bodyExpr, err := p.expression(0, end...)
	if err != nil {
		return nil, err
	}
	method := &methodDecl{name: "{closure}", parameters: parameters, returnType: returnType, returnNullable: nullable, arrow: true, body: []statement{{kind: "return", expressions: []expr{bodyExpr}, line: expressionLine(bodyExpr)}}, line: line, source: sourceContext{p.filename, p.namespace, p.strictTypes}, version: p.version, returnsByReference: byRefReturn}
	return closureExpr{method: method, auto: true, static: static, line: line}, nil
}

func (e *evaluator) evalClosure(node closureExpr) value {
	e.initClasses()
	method := *node.method
	method.owner = e.scope
	closure := &phpClosure{method: &method, captures: make(variableScope), called: e.lateClass, static: node.static}
	if !node.static {
		if this, ok := e.readVariable("this"); ok && this.kind == objectKind {
			closure.boundThis = this
		}
	}
	if node.auto {
		for name, v := range e.variables {
			if name != "this" {
				closure.captures[name] = dereference(v)
			}
		}
		for name := range e.bindings {
			if name != "this" {
				if v, ok := e.readVariable(name); ok {
					closure.captures[name] = v
				}
			}
		}
	} else {
		for _, capture := range node.captures {
			if capture.byRef {
				closure.captures[capture.name] = value{kind: refKind, ref: e.ensureVariableReference(capture.name)}
			} else if v, ok := e.readVariable(capture.name); ok {
				closure.captures[capture.name] = v
			} else {
				e.warn(node.line, "Undefined variable $"+capture.name)
				closure.captures[capture.name] = value{}
			}
		}
	}
	return value{kind: objectKind, o: &phpObject{class: e.classes["closure"], properties: &phpArray{}, closure: closure}}
}

func (e *evaluator) invokeClosure(closure *phpClosure, args []expr, line int) (value, error) {
	values, err := e.evalArguments(args)
	if err != nil {
		return value{}, err
	}
	return e.invokeClosureValues(closure, values, line)
}

func (e *evaluator) invokeClosureValues(closure *phpClosure, values []callArgument, line int) (value, error) {
	if closure.target != nil {
		return e.invokeResolved(*closure.target, values, line)
	}
	called := closure.called
	if closure.boundThis.kind == objectKind {
		called = closure.boundThis.o.class
	}
	return e.invokeWithCapture(closure.boundThis, called, closure.method, values, line, closure)
}

func (e *evaluator) closureBind(expressions []expr, line int) (value, error) {
	args, err := e.evalArguments(expressions)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "closure"}, {name: "newThis"}, {name: "newScope"}}, line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || !present[1] || len(args) > 3 {
		return value{}, &runtimeError{line, "ArgumentCountError", "Closure::bind() expects between 2 and 3 arguments"}
	}
	if bound[0].kind != objectKind || bound[0].o == nil || bound[0].o.closure == nil {
		return value{}, builtinTypeError(line, "Closure::bind", 1, "closure", "Closure", bound[0])
	}
	if bound[1].kind != nullKind && bound[1].kind != objectKind {
		return value{}, builtinTypeError(line, "Closure::bind", 2, "newThis", "?object", bound[1])
	}
	original := bound[0].o.closure
	if original.target != nil {
		return value{}, &runtimeError{line, "Error", "gophp: rebinding a first-class callable is not supported"}
	}
	if original.static && bound[1].kind == objectKind {
		if err := e.warn(line, "Cannot bind an instance to a static closure"); err != nil {
			return value{}, err
		}
		return value{kind: nullKind}, nil
	}
	owner := original.method.owner
	if present[2] {
		switch bound[2].kind {
		case nullKind:
			owner = nil
		case objectKind:
			owner = bound[2].o.class
		case stringKind:
			if !strings.EqualFold(bound[2].s, "static") {
				owner, err = e.resolveClassName(bound[2].s, line)
				if err != nil {
					return value{}, err
				}
			}
		default:
			return value{}, builtinTypeError(line, "Closure::bind", 3, "newScope", "object|string|null", bound[2])
		}
	}
	method := *original.method
	method.owner = owner
	captures := make(variableScope, len(original.captures))
	for name, captured := range original.captures {
		captures[name] = captured
	}
	copy := &phpClosure{method: &method, captures: captures, boundThis: bound[1], static: original.static, called: owner}
	if bound[1].kind == objectKind {
		copy.called = bound[1].o.class
	}
	return value{kind: objectKind, o: &phpObject{class: e.classes["closure"], properties: &phpArray{}, closure: copy}}, nil
}
