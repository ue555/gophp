package php

type phpClosure struct {
	method    *methodDecl
	captures  variableScope
	boundThis value
	called    *classDecl
}

func (p *parser) anonymousFunction(line int, static bool) (expr, error) {
	if p.accept("&") {
		return nil, syntax(line, "unsupported reference return from closure")
	}
	parameters, err := p.parameterList()
	if err != nil {
		return nil, err
	}
	var captures []string
	if p.acceptKeyword("use") {
		if err := p.expect("("); err != nil {
			return nil, err
		}
		seen := make(map[string]bool)
		for _, parameter := range parameters {
			seen[parameter.name] = true
		}
		for p.current().kind != ")" {
			name := p.current()
			if err := p.expect("variable"); err != nil {
				return nil, err
			}
			if name.text == "this" || seen[name.text] {
				return nil, syntax(name.line, "invalid closure capture $%s", name.text)
			}
			seen[name.text] = true
			captures = append(captures, name.text)
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
	method := &methodDecl{name: "{closure}", parameters: parameters, returnType: returnType, returnNullable: nullable, body: body, line: line, source: sourceContext{p.filename, p.namespace, p.strictTypes}}
	return closureExpr{method: method, captures: captures, static: static, line: line}, nil
}

func (p *parser) arrowFunction(line int, static bool, end ...string) (expr, error) {
	parameters, err := p.parameterList()
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
	method := &methodDecl{name: "{closure}", parameters: parameters, returnType: returnType, returnNullable: nullable, body: []statement{{kind: "return", expressions: []expr{bodyExpr}, line: expressionLine(bodyExpr)}}, line: line, source: sourceContext{p.filename, p.namespace, p.strictTypes}}
	return closureExpr{method: method, auto: true, static: static, line: line}, nil
}

func (e *evaluator) evalClosure(node closureExpr) value {
	e.initClasses()
	method := *node.method
	method.owner = e.scope
	closure := &phpClosure{method: &method, captures: make(variableScope), called: e.lateClass}
	if !node.static {
		if this, ok := e.readVariable("this"); ok && this.kind == objectKind {
			closure.boundThis = this
		}
	}
	if node.auto {
		for name, v := range e.variables {
			if name != "this" {
				closure.captures[name] = v
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
		for _, name := range node.captures {
			if v, ok := e.readVariable(name); ok {
				closure.captures[name] = v
			} else {
				e.warn(node.line, "Undefined variable $"+name)
				closure.captures[name] = value{}
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
	called := closure.called
	if closure.boundThis.kind == objectKind {
		called = closure.boundThis.o.class
	}
	return e.invokeWithCapture(closure.boundThis, called, closure.method, values, line, closure)
}
