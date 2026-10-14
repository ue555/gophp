package php

import (
	"fmt"
	"strings"
)

func equalName(a, b string) bool { return strings.EqualFold(a, b) }

func reservedName(name string) bool {
	switch strings.ToLower(name) {
	case "if", "else", "elseif", "while", "do", "for", "foreach", "switch", "case", "default", "break", "continue", "return", "class", "function", "fn", "new", "echo", "print", "array", "true", "false", "null", "and", "or", "xor", "try", "catch", "finally", "throw", "global", "static", "unset", "isset", "empty", "include", "require", "include_once", "require_once", "instanceof", "match", "clone", "yield", "exit", "die", "list", "eval", "declare", "namespace", "use", "const":
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
	if err := p.registerDeclaration(name, "function"); err != nil {
		return nil, err
	}
	decl := &methodDecl{name: qualify(p.namespace, name.text), line: name.line, source: sourceContext{p.filename, p.namespace, p.strictTypes}}
	var err error
	decl.parameters, decl.body, decl.returnType, decl.returnNullable, err = p.callableBody()
	return decl, err
}

// Methods and functions deliberately share argument binding and return handling.
func (p *parser) callableBody() ([]parameter, []statement, string, bool, error) {
	parameters, err := p.parameterList()
	if err != nil {
		return nil, nil, "", false, err
	}
	returnType, nullable, err := p.returnDeclaration()
	if err != nil {
		return nil, nil, "", false, err
	}
	body, err := p.block()
	return parameters, body, returnType, nullable, err
}

func (p *parser) returnDeclaration() (string, bool, error) {
	if !p.accept(":") {
		return "", false, nil
	}
	return p.typeDeclaration(true)
}

// The first implementation unit supports simple scalar, array, object and
// class declarations. Compound forms are rejected by the surrounding parser
// instead of being silently discarded.
func (p *parser) typeDeclaration(required bool) (string, bool, error) {
	nullable := p.accept("?")
	if !required && !nullable && (p.current().kind == "variable" || p.current().kind == "...") {
		return "", false, nil
	}
	name, err := p.readName()
	if err != nil {
		return "", false, err
	}
	typeName := strings.ToLower(name.text)
	switch typeName {
	case "int", "float", "string", "bool", "array", "object", "mixed":
		if nullable && typeName == "mixed" {
			return "", false, syntax(name.line, "mixed type cannot be nullable")
		}
	case "void", "callable", "iterable", "never", "null", "false", "true", "self", "parent", "static":
		return "", false, syntax(name.line, "unsupported type declaration %s", name.text)
	default:
		typeName, _ = p.resolveName(name.text, "class")
	}
	return typeName, nullable, nil
}

func (p *parser) parameterList() ([]parameter, error) {
	if err := p.expect("("); err != nil {
		return nil, err
	}
	var parameters []parameter
	seen := make(map[string]bool)
	for p.current().kind != ")" {
		typeName, nullable, err := p.typeDeclaration(false)
		if err != nil {
			return nil, err
		}
		variadic := p.accept("...")
		param := p.current()
		if err := p.expect("variable"); err != nil {
			return nil, err
		}
		if param.text == "this" || seen[param.text] {
			return nil, syntax(param.line, "invalid or duplicate parameter $%s", param.text)
		}
		seen[param.text] = true
		parameter := parameter{name: param.text, line: param.line, typeName: typeName, nullable: nullable, variadic: variadic}
		if variadic && p.current().kind == "=" {
			return nil, syntax(param.line, "variadic parameter cannot have a default value")
		}
		if p.accept("=") {
			var err error
			parameter.initial, err = p.expression(0)
			if err != nil {
				return nil, err
			}
		}
		parameters = append(parameters, parameter)
		if !p.accept(",") {
			break
		}
		if variadic && p.current().kind != ")" {
			return nil, syntax(p.current().line, "variadic parameter must be the last parameter")
		}
	}
	if err := p.expect(")"); err != nil {
		return nil, err
	}
	return parameters, nil
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
		if p.initial == nil && !p.variadic {
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
		for _, c := range s.cases {
			e.callableDiagnostics(c.body)
		}
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
	if node.fallback != "" && e.functions[symbolKey(node.name)] == nil {
		if e.functions[symbolKey(node.fallback)] != nil || builtinName(node.fallback) {
			node.name = node.fallback
		}
	}
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
