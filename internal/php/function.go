package php

import (
	"fmt"
	"strings"
)

func reservedName(name string) bool {
	switch strings.ToLower(name) {
	case "if", "else", "elseif", "while", "do", "for", "foreach", "switch", "case", "default", "break", "continue", "return", "class", "function", "fn", "new", "echo", "print", "array", "true", "false", "null", "and", "or", "xor", "try", "catch", "finally", "throw", "global", "static", "unset", "isset", "empty", "include", "require", "match", "clone", "yield", "exit", "die", "list", "eval", "declare", "namespace", "use":
		return true
	}
	return false
}

func (p *parser) functionDeclaration() (*methodDecl, error) {
	name := p.current()
	if err := p.expect("name"); err != nil {
		return nil, err
	}
	if reservedName(name.text) {
		return nil, syntax(name.line, "unsupported function name %q", name.text)
	}
	decl := &methodDecl{name: name.text, line: name.line}
	var err error
	decl.parameters, decl.body, err = p.callableBody()
	return decl, err
}

// Methods and functions deliberately share argument binding and return handling.
func (p *parser) callableBody() ([]parameter, []statement, error) {
	if err := p.expect("("); err != nil {
		return nil, nil, err
	}
	var parameters []parameter
	seen := make(map[string]bool)
	for p.current().kind != ")" {
		param := p.current()
		if err := p.expect("variable"); err != nil {
			return nil, nil, err
		}
		if param.text == "this" || seen[param.text] {
			return nil, nil, syntax(param.line, "invalid or duplicate parameter $%s", param.text)
		}
		seen[param.text] = true
		parameter := parameter{name: param.text, line: param.line}
		if p.accept("=") {
			var err error
			parameter.initial, err = p.expression(0)
			if err != nil {
				return nil, nil, err
			}
		}
		parameters = append(parameters, parameter)
		if !p.accept(",") {
			break
		}
	}
	if err := p.expect(")"); err != nil {
		return nil, nil, err
	}
	body, err := p.block()
	return parameters, body, err
}

type fatalDiagnostic struct {
	line    int
	message string
}

func (f *fatalDiagnostic) Error() string { return f.message }

func redeclaration(decl, old *methodDecl, filename string, version Version) string {
	prefix := "Cannot redeclare "
	if version == PHP84 || version == PHP85 {
		prefix += "function "
	}
	return fmt.Sprintf("%s%s() (previously declared in %s:%d)", prefix, decl.name, filename, old.line)
}

func duplicateFunctions(statements []statement, filename string, version Version) *runtimeError {
	seen := make(map[string]*methodDecl)
	for _, s := range statements {
		if s.kind != "function" {
			continue
		}
		key := strings.ToLower(s.function.name)
		if builtinName(key) {
			return &runtimeError{s.function.line, "", builtinRedeclaration(s.function.name, version)}
		}
		if old := seen[key]; old != nil {
			return &runtimeError{s.function.line, "", redeclaration(s.function, old, filename, version)}
		}
		seen[key] = s.function
	}
	return nil
}

func (e *evaluator) declareFunction(decl *methodDecl) error {
	if e.functions == nil {
		e.functions = make(map[string]*methodDecl)
	}
	key := strings.ToLower(decl.name)
	if builtinName(key) {
		return &fatalDiagnostic{decl.line, builtinRedeclaration(decl.name, e.version)}
	}
	if old := e.functions[key]; old != nil {
		return &fatalDiagnostic{decl.line, redeclaration(decl, old, e.filename, e.version)}
	}
	e.functions[key] = decl
	return nil
}

func requiredParameters(decl *methodDecl) int {
	last := 0
	for i, p := range decl.parameters {
		if p.initial == nil {
			last = i + 1
		}
	}
	return last
}

// These are compilation diagnostics, including declarations in skipped bodies.
func (e *evaluator) callableDiagnostics(statements []statement) {
	for _, s := range statements {
		if s.function != nil {
			e.parameterDiagnostics(s.function, s.function.name)
		}
		if s.class != nil {
			for _, m := range s.class.methodOrder {
				e.parameterDiagnostics(m, s.class.name+"::"+m.name)
			}
		}
		e.callableDiagnostics(s.thenBranch)
		e.callableDiagnostics(s.elseBranch)
	}
}

func (e *evaluator) parameterDiagnostics(decl *methodDecl, name string) {
	last := requiredParameters(decl)
	for i, p := range decl.parameters {
		if p.initial != nil && i < last {
			message := fmt.Sprintf("Optional parameter $%s declared before required parameter $%s is implicitly treated as a required parameter", p.name, decl.parameters[last-1].name)
			if e.version != PHP83 {
				message = name + "(): " + message
			}
			e.diagnostic("Deprecated", decl.line, message)
		}
	}
	e.callableDiagnostics(decl.body)
}

func (e *evaluator) callFunction(node functionCall) (value, error) {
	if builtinName(node.name) {
		return e.callBuiltin(node)
	}
	decl := e.functions[strings.ToLower(node.name)]
	if decl == nil {
		return value{}, &runtimeError{node.line, "Error", "Call to undefined function " + node.name + "()"}
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	return e.invoke(value{}, decl, args, node.line)
}
