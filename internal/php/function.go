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
	byRefReturn := p.accept("&")
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
	decl := &methodDecl{name: qualify(p.namespace, name.text), line: name.line, source: sourceContext{p.filename, p.namespace, p.strictTypes}, version: p.version, returnsByReference: byRefReturn}
	previousCallable := p.currentCallable
	p.currentCallable = decl.name
	defer func() { p.currentCallable = previousCallable }()
	var err error
	decl.parameters, decl.body, decl.returnType, decl.returnNullable, err = p.callableBody(false)
	return decl, err
}

// Methods and functions deliberately share argument binding and return handling.
func (p *parser) callableBody(allowPromotion bool) ([]parameter, []statement, string, bool, error) {
	parameters, err := p.parameterList(allowPromotion)
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
	return p.typeDeclarationAt(true, "return")
}

// Keep each union arm and its intersection members distinct. Type matching
// uses the same structure; diagnostics use declaredTypeText.
func (p *parser) typeDeclaration(required bool) (string, bool, error) {
	return p.typeDeclarationAt(required, "other")
}

func (p *parser) typeDeclarationAt(required bool, position string) (string, bool, error) {
	nullable := p.accept("?")
	if !required && !nullable && (p.current().kind == "variable" || p.current().kind == "..." || p.current().kind == "&") {
		return "", false, nil
	}
	var arms []string
	seen := make(map[string]bool)
	for {
		grouped := p.accept("(")
		var parts []string
		intersectionSeen := make(map[string]bool)
		for {
			name, err := p.readName()
			if err != nil {
				return "", false, err
			}
			member := strings.ToLower(name.text)
			switch member {
			case "int", "float", "string", "bool", "array", "object", "mixed", "null", "false", "true":
			case "void", "never":
			case "callable", "iterable":
			case "self", "parent", "static":
				if p.classDepth == 0 {
					return "", false, &fatalDiagnostic{name.line, "Cannot use " + member + " when no class scope is active"}
				}
				if member == "parent" && p.currentClass != nil && p.currentClass.kind != "trait" && p.currentClass.parentName == "" {
					return "", false, &fatalDiagnostic{name.line, "Cannot use \"parent\" when current class scope has no parent"}
				}
				if member == "static" && position != "return" {
					return "", false, syntax(name.line, "static type can only be used as a return type")
				}
			default:
				member, _ = p.resolveName(name.text, "class")
			}
			if intersectionSeen[strings.ToLower(member)] {
				return "", false, &fatalDiagnostic{name.line, "Duplicate type " + member + " is redundant"}
			}
			intersectionSeen[strings.ToLower(member)] = true
			parts = append(parts, member)
			// In a parameter declaration, `T &$value` uses & for the
			// reference parameter. `T&U $value` remains an intersection.
			if position == "other" && p.current().kind == "&" && p.pos+1 < len(p.tokens) && (p.tokens[p.pos+1].kind == "variable" || p.tokens[p.pos+1].kind == "...") {
				break
			}
			if !p.accept("&") {
				break
			}
		}
		if grouped {
			if err := p.expect(")"); err != nil {
				return "", false, err
			}
			if len(parts) < 2 {
				return "", false, syntax(p.current().line, "intersection type requires at least two types")
			}
		}
		if len(parts) > 1 {
			if nullable {
				return "", false, syntax(p.current().line, "syntax error, unexpected token \"&\", expecting variable")
			}
			if len(arms) > 0 && !grouped {
				return "", false, syntax(p.current().line, "syntax error, unexpected token \"&\"")
			}
			for _, part := range parts {
				if builtinTypeName(part) || part == "static" {
					return "", false, &fatalDiagnostic{p.current().line, "Type " + part + " cannot be part of an intersection type"}
				}
			}
		}
		arm := strings.Join(parts, "&")
		if grouped {
			arm = "(" + arm + ")"
		}
		key := typeAlternative{types: parts}.key()
		if seen[key] {
			if len(parts) > 1 {
				for _, prior := range arms {
					other := parseDeclaredType(prior, false).alternatives[0]
					if other.key() == key {
						return "", false, &fatalDiagnostic{p.current().line, "Type " + strings.Join(parts, "&") + " is redundant with type " + strings.Join(other.types, "&")}
					}
				}
			}
			return "", false, &fatalDiagnostic{p.current().line, "Duplicate type " + arm + " is redundant"}
		}
		for _, prior := range arms {
			other := parseDeclaredType(prior, false).alternatives[0]
			current := typeAlternative{types: parts}
			if current.moreRestrictiveThan(other) {
				return "", false, &fatalDiagnostic{p.current().line, "Type " + strings.Join(parts, "&") + " is redundant as it is more restrictive than type " + strings.Join(other.types, "&")}
			}
			if other.moreRestrictiveThan(current) {
				return "", false, &fatalDiagnostic{p.current().line, "Type " + strings.Join(other.types, "&") + " is redundant as it is more restrictive than type " + strings.Join(parts, "&")}
			}
		}
		seen[key] = true
		arms = append(arms, arm)
		if nullable && p.current().kind == "|" {
			if required {
				return "", false, syntax(p.current().line, "syntax error, unexpected token \"|\"")
			}
			return "", false, syntax(p.current().line, "syntax error, unexpected token \"|\", expecting variable")
		}
		if p.current().kind == "|" && len(parts) > 1 && !grouped {
			if !required {
				return "", false, syntax(p.current().line, "syntax error, unexpected token \"|\", expecting variable")
			}
			return "", false, syntax(p.current().line, "syntax error, unexpected token \"|\"")
		}
		if !p.accept("|") {
			break
		}
	}
	if len(arms) > 1 {
		for _, m := range arms {
			switch m {
			case "mixed", "void", "never":
				label := m
				if m == "void" {
					label = "Void"
				} else if m == "mixed" {
					label = "Type mixed"
				}
				return "", false, &fatalDiagnostic{p.current().line, label + " can only be used as a standalone type"}
			}
		}
		if seen["bool"] && (seen["true"] || seen["false"]) || seen["true"] && seen["false"] {
			if seen["true"] && seen["false"] && !seen["bool"] {
				return "", false, &fatalDiagnostic{p.current().line, "Type contains both true and false, bool must be used instead"}
			}
			if seen["bool"] && seen["true"] {
				return "", false, &fatalDiagnostic{p.current().line, "Duplicate type true is redundant"}
			}
			return "", false, &fatalDiagnostic{p.current().line, "Duplicate type false is redundant"}
		}
		if seen["object"] {
			for _, m := range arms {
				if !builtinTypeName(m) {
					return "", false, &fatalDiagnostic{p.current().line, "Type " + strings.Join(arms, "|") + " contains both object and a class type, which is redundant"}
				}
			}
		}
	} else if nullable && (arms[0] == "mixed" || arms[0] == "null" || arms[0] == "void" || arms[0] == "never") {
		if arms[0] == "void" || arms[0] == "never" {
			label := arms[0]
			if label == "void" {
				label = "Void"
			}
			return "", false, &fatalDiagnostic{p.current().line, label + " can only be used as a standalone type"}
		}
		return "", false, syntax(p.current().line, "%s type cannot be nullable", arms[0])
	}
	if len(arms) == 1 && strings.HasPrefix(arms[0], "(") {
		return "", false, syntax(p.current().line, "parenthesized intersection requires a union type")
	}
	return strings.Join(arms, "|"), nullable, nil
}

func (p *parser) parameterList(allowPromotion bool) ([]parameter, error) {
	if err := p.expect("("); err != nil {
		return nil, err
	}
	var parameters []parameter
	seen := make(map[string]bool)
	for p.current().kind != ")" {
		attributes, err := p.attributes()
		if err != nil {
			return nil, err
		}
		visibility := ""
		readonly := false
		for {
			if p.keyword("public") || p.keyword("protected") || p.keyword("private") {
				if visibility != "" {
					return nil, unexpected(p.current())
				}
				visibility = strings.ToLower(p.advance().text)
				continue
			}
			if p.acceptKeyword("readonly") {
				if readonly {
					return nil, unexpected(p.tokens[p.pos-1])
				}
				readonly = true
				continue
			}
			break
		}
		promoted := visibility != "" || readonly
		if promoted && !allowPromotion {
			return nil, &fatalDiagnostic{p.current().line, "Cannot declare promoted property outside a constructor"}
		}
		if readonly && visibility == "" {
			visibility = "public"
		}
		typeName, nullable, err := p.typeDeclaration(false)
		if err != nil {
			return nil, err
		}
		if promoted {
			if err := p.validMemberType(typeName, "property"); err != nil {
				return nil, err
			}
		}
		if typeName == "void" || typeName == "never" {
			return nil, &fatalDiagnostic{p.current().line, typeName + " cannot be used as a parameter type"}
		}
		byRef := p.accept("&")
		variadic := p.accept("...")
		if variadic && promoted {
			return nil, &fatalDiagnostic{p.current().line, "Cannot declare variadic promoted property"}
		}
		param := p.current()
		if err := p.expect("variable"); err != nil {
			return nil, err
		}
		if param.text == "this" || seen[param.text] {
			return nil, syntax(param.line, "invalid or duplicate parameter $%s", param.text)
		}
		seen[param.text] = true
		parameter := parameter{name: param.text, line: param.line, typeName: typeName, nullable: nullable, variadic: variadic, byRef: byRef, promoted: promoted, visibility: visibility, readonly: readonly, attributes: attributes}
		if promoted && readonly && typeName == "" {
			return nil, &fatalDiagnostic{param.line, "Readonly property must have type"}
		}
		if variadic && p.current().kind == "=" {
			return nil, syntax(param.line, "variadic parameter cannot have a default value")
		}
		if p.accept("=") {
			var err error
			parameter.initial, err = p.expression(0)
			if err != nil {
				return nil, err
			}
			if initial, ok := parameter.initial.(literal); ok && initial.value.kind == nullKind && promoted && typeName != "" && !nullable && !typeIncludesNull(typeName) {
				return nil, &fatalDiagnostic{param.line, fmt.Sprintf("Cannot use null as default value for parameter $%s of type %s", param.text, declaredTypeText(typeName, false))}
			}
			if initial, ok := parameter.initial.(literal); ok && initial.value.kind == nullKind && typeName != "" && !nullable && !typeIncludesNull(typeName) {
				parameter.nullable = true
				parameter.implicitNullable = true
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
		if builtinNameForVersion(key, version) {
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
	if builtinNameForVersion(key, e.version) {
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
		if p.implicitNullable && e.version != PHP83 {
			e.diagnostic("Deprecated", decl.line, fmt.Sprintf("%s(): Implicitly marking parameter $%s as nullable is deprecated, the explicit nullable type must be used instead", name, p.name))
		}
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
	if strings.EqualFold(node.name, "eval") {
		return e.callEval(node)
	}
	if node.fallback != "" && e.functions[symbolKey(node.name)] == nil {
		if e.functions[symbolKey(node.fallback)] != nil || builtinNameForVersion(node.fallback, e.version) {
			node.name = node.fallback
		}
	}
	if builtinNameForVersion(node.name, e.version) {
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
