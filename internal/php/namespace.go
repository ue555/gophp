package php

import (
	"fmt"
	"math"
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
	name       string
	initial    expr
	line       int
	attributes []attributeDecl
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
	// `new \ClassName` is a keyword followed by a fully qualified class
	// name. A contiguous `new\name` is instead one qualified name token.
	if strings.EqualFold(t.text, "new") && p.current().kind == "\\" && p.current().start != t.start+len(t.text) {
		return t, nil
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
	var attributes []attributeDecl
	if p.current().kind == "#[" {
		start := p.pos
		var err error
		attributes, err = p.attributes()
		if err != nil {
			return nil, err
		}
		if !p.keyword("const") {
			p.pos = start
		} else if p.version != PHP85 {
			return nil, syntax(p.current().line, "syntax error, unexpected token \"const\"")
		}
	}
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
			result = append(result, statement{kind: "const", constant: &constantDecl{name: qualify(p.namespace, t.text), initial: x, line: t.line, attributes: attributes}, line: t.line})
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
		e.constants = map[string]value{
			"COUNT_NORMAL": integer(0), "COUNT_RECURSIVE": integer(1),
			"ARRAY_FILTER_USE_BOTH": integer(1), "ARRAY_FILTER_USE_KEY": integer(2),
			"CASE_LOWER": integer(0), "CASE_UPPER": integer(1),
			"ENT_COMPAT": integer(2), "ENT_QUOTES": integer(3), "ENT_NOQUOTES": integer(0), "ENT_SUBSTITUTE": integer(8), "ENT_HTML5": integer(48),
			"XML_PI_NODE": integer(7), "XML_COMMENT_NODE": integer(8), "LIBXML_NOBLANKS": integer(256), "LIBXML_NONET": integer(2048),
			"LIBXML_NOERROR": integer(32), "LIBXML_COMPACT": integer(65536), "LIBXML_HTML_NODEFDTD": integer(4), "LIBXML_NOXMLDECL": integer(8192),
			"PHP_QUERY_RFC1738": integer(1), "PHP_QUERY_RFC3986": integer(2),
			"DEBUG_BACKTRACE_PROVIDE_OBJECT": integer(1), "DEBUG_BACKTRACE_IGNORE_ARGS": integer(2),
			"E_ERROR": integer(1), "E_WARNING": integer(2), "E_PARSE": integer(4), "E_NOTICE": integer(8),
			"E_CORE_ERROR": integer(16), "E_CORE_WARNING": integer(32), "E_COMPILE_ERROR": integer(64), "E_COMPILE_WARNING": integer(128),
			"E_USER_ERROR": integer(256), "E_USER_WARNING": integer(512), "E_USER_NOTICE": integer(1024),
			"E_STRICT": integer(2048), "E_RECOVERABLE_ERROR": integer(4096), "E_DEPRECATED": integer(8192),
			"E_USER_DEPRECATED": integer(16384), "E_ALL": integer(allErrorLevels(e.version)),
			"DIRECTORY_SEPARATOR": {kind: stringKind, s: "/"},
			"PHP_SAPI":            {kind: stringKind, s: "cli"},
			"PHP_OS":              {kind: stringKind, s: "Linux"},
			"PHP_OS_FAMILY":       {kind: stringKind, s: "Linux"},
			"PHP_EOL":             {kind: stringKind, s: "\n"},
			"PHP_INT_MAX":         integer(9223372036854775807), "PHP_INT_MIN": integer(-9223372036854775808), "PHP_INT_SIZE": integer(8),
			"PHP_MAXPATHLEN": integer(4096),
			"INF":            floating(math.Inf(1)), "NAN": floating(math.NaN()),
			"JSON_HEX_TAG": integer(1), "JSON_HEX_AMP": integer(2), "JSON_HEX_APOS": integer(4), "JSON_HEX_QUOT": integer(8),
			"JSON_UNESCAPED_SLASHES": integer(64), "JSON_PRETTY_PRINT": integer(128), "JSON_UNESCAPED_UNICODE": integer(256), "JSON_INVALID_UTF8_IGNORE": integer(1048576),
			"JSON_PRESERVE_ZERO_FRACTION": integer(1024), "JSON_INVALID_UTF8_SUBSTITUTE": integer(2097152),
			"JSON_PARTIAL_OUTPUT_ON_ERROR": integer(512), "JSON_THROW_ON_ERROR": integer(4194304),
			"JSON_ERROR_NONE": integer(0), "JSON_ERROR_DEPTH": integer(1), "JSON_ERROR_STATE_MISMATCH": integer(2),
			"JSON_ERROR_CTRL_CHAR": integer(3), "JSON_ERROR_SYNTAX": integer(4), "JSON_ERROR_UTF8": integer(5),
			"PREG_PATTERN_ORDER": integer(1), "PREG_SET_ORDER": integer(2),
			"PREG_OFFSET_CAPTURE": integer(256), "PREG_UNMATCHED_AS_NULL": integer(512),
			"PREG_SPLIT_NO_EMPTY": integer(1), "PREG_SPLIT_DELIM_CAPTURE": integer(2), "PREG_SPLIT_OFFSET_CAPTURE": integer(4),
			"STR_PAD_LEFT": integer(0), "STR_PAD_RIGHT": integer(1), "STR_PAD_BOTH": integer(2),
			"SORT_REGULAR": integer(0), "SORT_NUMERIC": integer(1), "SORT_STRING": integer(2), "SORT_NATURAL": integer(6), "SORT_FLAG_CASE": integer(8),
			"LOCK_SH": integer(1), "LOCK_EX": integer(2), "LOCK_UN": integer(3), "LOCK_NB": integer(4), "FILE_APPEND": integer(8),
			"SEEK_SET": integer(0), "SEEK_CUR": integer(1), "SEEK_END": integer(2),
			"GLOB_ONLYDIR": integer(8192), "GLOB_NOSORT": integer(4), "GLOB_MARK": integer(2),
			"PHP_URL_SCHEME": integer(0), "PHP_URL_HOST": integer(1), "PHP_URL_PORT": integer(2), "PHP_URL_USER": integer(3), "PHP_URL_PASS": integer(4), "PHP_URL_PATH": integer(5), "PHP_URL_QUERY": integer(6), "PHP_URL_FRAGMENT": integer(7),
			"LOG_USER":       integer(8),
			"EXTR_OVERWRITE": integer(0), "EXTR_SKIP": integer(1), "EXTR_PREFIX_SAME": integer(2), "EXTR_REFS": integer(256),
			"FILTER_DEFAULT": integer(516), "FILTER_CALLBACK": integer(1024), "FILTER_VALIDATE_INT": integer(257), "FILTER_VALIDATE_BOOL": integer(258), "FILTER_VALIDATE_BOOLEAN": integer(258), "FILTER_VALIDATE_URL": integer(273), "FILTER_VALIDATE_IP": integer(275), "FILTER_FLAG_IPV4": integer(1048576), "FILTER_FLAG_IPV6": integer(2097152), "FILTER_REQUIRE_ARRAY": integer(16777216), "FILTER_REQUIRE_SCALAR": integer(33554432), "FILTER_NULL_ON_FAILURE": integer(134217728),
		}
		// The comparison targets pin patch versions. Keep this runtime identity
		// in one place so version checks in Composer and vendor code use the
		// same target as diagnostic compatibility rules.
		switch e.version {
		case PHP83:
			e.constants["PHP_VERSION_ID"] = integer(80335)
			e.constants["PHP_VERSION"] = value{kind: stringKind, s: "8.3.35"}
		case PHP85:
			e.constants["PHP_VERSION_ID"] = integer(80511)
			e.constants["PHP_VERSION"] = value{kind: stringKind, s: "8.5.11"}
		default:
			e.constants["PHP_VERSION_ID"] = integer(80426)
			e.constants["PHP_VERSION"] = value{kind: stringKind, s: "8.4.26"}
		}
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
		if err := e.constantDeprecation(n.name, n.line); err != nil {
			return value{}, err
		}
		return v, nil
	}
	if n.fallback != "" {
		if v, ok := e.constants[n.fallback]; ok {
			if err := e.constantDeprecation(n.fallback, n.line); err != nil {
				return value{}, err
			}
			return v, nil
		}
	}
	return value{}, &runtimeError{n.line, "Error", fmt.Sprintf("Undefined constant %q", n.name)}
}

func (e *evaluator) constantDeprecation(name string, line int) error {
	for _, attribute := range e.constantAttributes[name] {
		if strings.EqualFold(attribute.name, "Deprecated") {
			message := "Constant " + name + " is deprecated"
			return e.emitError(phpEUserDeprecated, "Deprecated", line, message)
		}
	}
	return nil
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
