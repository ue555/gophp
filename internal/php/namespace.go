package php

import (
	"fmt"
	"strings"
)

type sourceContext struct {
	filename, namespace string
	strictTypes         bool
}
type constantExpr struct {
	name, fallback string
	line           int
}
type constantDecl struct {
	name    string
	initial expr
	line    int
}

func (constantExpr) expression() {}

func symbolKey(name string) string { return strings.ToLower(name) }
func qualify(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return namespace + "\\" + name
}

// Imports are lexical, resolved into each AST name at its source position.
// Only unqualified functions/constants may fall back at runtime.
func (p *parser) resolveName(raw, kind string) (string, string) {
	if strings.HasPrefix(raw, "\\") {
		return raw[1:], ""
	}
	if strings.HasPrefix(strings.ToLower(raw), "namespace\\") {
		return qualify(p.namespace, raw[len("namespace\\"):]), ""
	}
	parts := strings.SplitN(raw, "\\", 2)
	lookupKind, key := kind, parts[0]
	if len(parts) > 1 {
		lookupKind = "class"
	}
	if lookupKind != "const" {
		key = symbolKey(key)
	}
	if imported := p.imports[lookupKind][key]; imported != "" {
		if len(parts) > 1 {
			imported += "\\" + parts[1]
		}
		return imported, ""
	}
	name := qualify(p.namespace, raw)
	if kind != "class" && len(parts) == 1 && p.namespace != "" {
		return name, raw
	}
	return name, ""
}

func (p *parser) readName() (token, error) {
	t := p.advance()
	return p.nameAfter(t)
}

func (p *parser) nameAfter(t token) (token, error) {
	if t.kind == "\\" {
		n := p.current()
		if err := p.expect("name"); err != nil {
			return t, err
		}
		t.kind, t.text = "name", "\\"+n.text
	} else if t.kind != "name" {
		return t, unexpected(t, "name")
	}
	for p.current().kind == "\\" && p.tokens[p.pos+1].kind == "name" {
		p.advance()
		t.text += "\\" + p.advance().text
	}
	return t, nil
}

func (p *parser) resetNamespace(name string) {
	p.namespace = name
	p.imports = map[string]map[string]string{"class": {}, "function": {}, "const": {}}
	p.declaredNames = map[string]map[string]bool{"class": {}, "function": {}, "const": {}}
}

func (p *parser) program() ([]statement, error) {
	p.resetNamespace("")
	var result []statement
	style, outside := "", false
	for p.current().kind != "eof" {
		if p.keyword("declare") {
			if outside || len(result) != 0 || style != "" || p.seenStrictTypes {
				return nil, &fatalDiagnostic{p.current().line, "strict_types declaration must be the very first statement in the script"}
			}
			line := p.advance().line
			if err := p.expect("("); err != nil {
				return nil, err
			}
			if !p.acceptKeyword("strict_types") {
				return nil, syntax(line, "unsupported declare directive")
			}
			if err := p.expect("="); err != nil {
				return nil, err
			}
			if p.current().kind != "number" || (p.current().text != "0" && p.current().text != "1") {
				return nil, syntax(p.current().line, "unsupported strict_types value")
			}
			strict := p.current().text == "1"
			p.advance()
			if err := p.expect(")"); err != nil {
				return nil, err
			}
			if err := p.expect(";"); err != nil {
				return nil, err
			}
			p.strictTypes = strict
			p.seenStrictTypes = true
			continue
		}
		if p.keyword("namespace") && p.tokens[p.pos+1].kind != "\\" {
			line := p.advance().line
			if outside && style == "" {
				return nil, &fatalDiagnostic{line, "Namespace declaration statement has to be the very first statement or after any declare call in the script"}
			}
			name := ""
			if p.current().kind != "{" {
				t, err := p.readName()
				if err != nil {
					return nil, err
				}
				if strings.HasPrefix(t.text, "\\") {
					return nil, unexpected(t)
				}
				name = t.text
			}
			form := "semicolon"
			if p.current().kind == "{" {
				form = "bracketed"
			}
			if style != "" && form != style {
				return nil, &fatalDiagnostic{line, "Cannot mix bracketed namespace declarations with unbracketed namespace declarations"}
			}
			style = form
			p.resetNamespace(name)
			if form == "semicolon" {
				if err := p.expect(";"); err != nil {
					return nil, err
				}
				continue
			}
			p.advance()
			for p.current().kind != "}" {
				stmts, err := p.topStatement()
				if err != nil {
					return nil, err
				}
				result = append(result, stmts...)
			}
			p.advance()
			p.resetNamespace("")
			continue
		}
		if style == "bracketed" {
			return nil, &fatalDiagnostic{p.current().line, "No code may exist outside of namespace {}"}
		}
		outside = true
		stmts, err := p.topStatement()
		if err != nil {
			return nil, err
		}
		result = append(result, stmts...)
	}
	return result, nil
}

func (p *parser) topStatement() ([]statement, error) {
	if p.acceptKeyword("use") {
		return nil, p.importDeclaration()
	}
	if p.acceptKeyword("const") {
		var result []statement
		for {
			t := p.current()
			if err := p.expect("name"); err != nil {
				return nil, err
			}
			if err := p.registerDeclaration(t, "const"); err != nil {
				return nil, err
			}
			if err := p.expect("="); err != nil {
				return nil, err
			}
			x, err := p.expression(0, ",", ";")
			if err != nil {
				return nil, err
			}
			result = append(result, statement{kind: "const", constant: &constantDecl{qualify(p.namespace, t.text), x, t.line}, line: t.line})
			if !p.accept(",") {
				break
			}
		}
		return result, p.expect(";")
	}
	s, err := p.statement()
	return []statement{s}, err
}

func (p *parser) registerDeclaration(t token, kind string) error {
	key := t.text
	if kind != "const" {
		key = symbolKey(key)
	}
	if kind == "const" && p.declaredNames[kind][key] {
		return &fatalDiagnostic{t.line, "Cannot redeclare constant " + qualify(p.namespace, t.text)}
	}
	if imported := p.imports[kind][key]; imported != "" && imported != qualify(p.namespace, t.text) {
		label := kind
		if kind == "const" {
			label = "const"
		}
		return &fatalDiagnostic{t.line, fmt.Sprintf("Cannot declare %s %s because the name is already in use", label, qualify(p.namespace, t.text))}
	}
	p.declaredNames[kind][key] = true
	return nil
}

func (p *parser) importDeclaration() error {
	kind := "class"
	if p.acceptKeyword("function") {
		kind = "function"
	} else if p.acceptKeyword("const") {
		kind = "const"
	}
	for {
		t, err := p.readName()
		if err != nil {
			return err
		}
		if p.accept("\\") {
			if err = p.expect("{"); err != nil {
				return err
			}
			for {
				memberKind := kind
				if p.keyword("function") || p.keyword("const") {
					if kind != "class" {
						return unexpected(p.current())
					}
					memberKind = strings.ToLower(p.advance().text)
				}
				member, err := p.readName()
				if err != nil {
					return err
				}
				member.text = strings.TrimPrefix(t.text, "\\") + "\\" + member.text
				if err = p.addImport(member, memberKind); err != nil {
					return err
				}
				if !p.accept(",") || p.current().kind == "}" {
					break
				}
			}
			if err = p.expect("}"); err != nil {
				return err
			}
		} else if err = p.addImport(t, kind); err != nil {
			return err
		}
		if !p.accept(",") {
			break
		}
	}
	return p.expect(";")
}

func (p *parser) addImport(t token, kind string) error {
	name := strings.TrimPrefix(t.text, "\\")
	parts := strings.Split(name, "\\")
	alias := parts[len(parts)-1]
	if p.acceptKeyword("as") {
		a := p.current()
		if err := p.expect("name"); err != nil {
			return err
		}
		alias = a.text
	}
	key := alias
	if kind != "const" {
		key = symbolKey(alias)
	}
	if reservedName(alias) {
		return syntax(t.line, "unsupported import alias %q", alias)
	}
	if p.imports[kind][key] != "" || p.declaredNames[kind][key] {
		prefix := ""
		if kind != "class" {
			prefix = kind + " "
		}
		return &fatalDiagnostic{t.line, fmt.Sprintf("Cannot use %s%s as %s because the name is already in use", prefix, name, alias)}
	}
	p.imports[kind][key] = name
	return nil
}

func (e *evaluator) initConstants() {
	if e.constants == nil {
		e.constants = map[string]value{"COUNT_NORMAL": integer(0), "COUNT_RECURSIVE": integer(1)}
	}
}

func duplicateConstants(statements []statement) *runtimeError {
	seen := make(map[string]bool)
	for _, s := range statements {
		if s.constant == nil {
			continue
		}
		if seen[s.constant.name] {
			return &runtimeError{s.constant.line, "", "Cannot redeclare constant " + s.constant.name}
		}
		seen[s.constant.name] = true
	}
	return nil
}
func (e *evaluator) constant(n constantExpr) (value, error) {
	e.initConstants()
	if v, ok := e.constants[n.name]; ok {
		return v, nil
	}
	if n.fallback != "" {
		if v, ok := e.constants[n.fallback]; ok {
			return v, nil
		}
	}
	return value{}, &runtimeError{n.line, "Error", fmt.Sprintf("Undefined constant %q", n.name)}
}
func (e *evaluator) defineConstant(name string, v value, line int) bool {
	e.initConstants()
	if _, ok := e.constants[name]; ok {
		e.warn(line, "Constant "+name+" already defined")
		return false
	}
	e.constants[name] = v
	return true
}
