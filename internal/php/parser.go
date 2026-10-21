package php

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type expr interface{ expression() }
type literal struct {
	value value
	line  int
}
type variable struct{ token token }
type unary struct {
	op      token
	operand expr
}
type binary struct {
	op          token
	left, right expr
}
type assignment struct {
	target, source expr
	line           int
}
type referenceAssignment struct {
	target, source expr
	line           int
}
type compoundAssignment struct {
	target, source expr
	op             token
}
type increment struct {
	target expr
	op     token
	post   bool
}
type functionCall struct {
	name     string
	args     []expr
	line     int
	fallback string
}
type languageExpr struct {
	name string
	args []expr
	line int
}
type castExpr struct {
	name    string
	operand expr
	line    int
}
type arrayElement struct {
	key, value expr
	byRef      bool
}
type arrayLiteral struct {
	elements []arrayElement
	line     int
}
type offsetExpr struct {
	base, key expr
	line      int
}
type propertyExpr struct {
	base     expr
	name     string
	line     int
	nameExpr expr
	grouped  bool
}
type newObject struct {
	name               string
	anonymous          *classDecl
	classExpr          expr
	line               int
	arguments, grouped bool
	args               []expr
}
type methodCall struct {
	receiver expr
	name     string
	nameExpr expr
	args     []expr
	line     int
}
type includeExpr struct {
	kind string
	path expr
	line int
}
type yieldExpr struct {
	key, item expr
	from      bool
	line      int
}
type cloneExpr struct {
	operand expr
	line    int
}
type dynamicVariable struct {
	nameExpr expr
	line     int
}
type instanceofExpr struct {
	operand   expr
	className string
	classExpr expr
	line      int
}
type dynamicCall struct {
	callee expr
	args   []expr
	line   int
}
type firstClassCallable struct {
	callee expr
	line   int
}
type unpackArgument struct {
	operand expr
	line    int
}
type namedArgument struct {
	name    string
	operand expr
	line    int
}
type throwExpr struct {
	operand expr
	line    int
}
type closureExpr struct {
	method   *methodDecl
	captures []closureCapture
	auto     bool
	static   bool
	line     int
}

type closureCapture struct {
	name  string
	byRef bool
}
type classNameExpr struct {
	name    string
	line    int
	dynamic expr
}
type classConstantExpr struct {
	className, name string
	line            int
}
type staticPropertyExpr struct {
	className, name     string
	line                int
	classExpr, nameExpr expr
}
type staticCallExpr struct {
	className, name     string
	args                []expr
	line                int
	classExpr, nameExpr expr
}

func (literal) expression()             {}
func (variable) expression()            {}
func (unary) expression()               {}
func (binary) expression()              {}
func (assignment) expression()          {}
func (referenceAssignment) expression() {}
func (compoundAssignment) expression()  {}
func (increment) expression()           {}
func (functionCall) expression()        {}
func (languageExpr) expression()        {}
func (castExpr) expression()            {}
func (arrayLiteral) expression()        {}
func (offsetExpr) expression()          {}
func (propertyExpr) expression()        {}
func (newObject) expression()           {}
func (yieldExpr) expression()           {}
func (cloneExpr) expression()           {}
func (dynamicVariable) expression()     {}
func (methodCall) expression()          {}
func (includeExpr) expression()         {}
func (instanceofExpr) expression()      {}
func (dynamicCall) expression()         {}
func (firstClassCallable) expression()  {}
func (unpackArgument) expression()      {}
func (namedArgument) expression()       {}
func (throwExpr) expression()           {}
func (closureExpr) expression()         {}
func (classNameExpr) expression()       {}
func (classConstantExpr) expression()   {}
func (staticPropertyExpr) expression()  {}
func (staticCallExpr) expression()      {}

// Boolean conversion opcodes in Zend use the line left by compiling their
// operand. For this expression subset, that is the last operand's token line.
func expressionLine(expression expr) int {
	switch node := expression.(type) {
	case literal:
		return node.line
	case unpackArgument:
		return node.line
	case namedArgument:
		return node.line
	case throwExpr:
		return node.line
	case variable:
		return node.token.line
	case unary:
		return expressionLine(node.operand)
	case binary:
		return expressionLine(node.right)
	case assignment:
		return expressionLine(node.source)
	case referenceAssignment:
		return node.line
	case compoundAssignment:
		return expressionLine(node.source)
	case increment:
		return node.op.line
	case functionCall:
		return node.line
	case languageExpr:
		return node.line
	case castExpr:
		return expressionLine(node.operand)
	case arrayLiteral:
		return node.line
	case offsetExpr:
		return node.line
	case propertyExpr:
		return node.line
	case newObject:
		return node.line
	case methodCall:
		return node.line
	case includeExpr:
		return node.line
	case yieldExpr:
		return node.line
	case cloneExpr:
		return node.line
	case dynamicVariable:
		return node.line
	case instanceofExpr:
		return node.line
	case dynamicCall:
		return node.line
	case firstClassCallable:
		return node.line
	case closureExpr:
		return node.line
	case ternaryExpr:
		return expressionLine(node.no)
	case coalesceAssignment:
		return expressionLine(node.source)
	case nullsafeExpr:
		return node.line
	case matchExpr:
		return node.line
	case constantExpr:
		return node.line
	case classNameExpr:
		return node.line
	case classConstantExpr:
		return node.line
	case staticPropertyExpr:
		return node.line
	case staticCallExpr:
		return node.line
	case interpolatedString:
		return node.line
	}
	panic("invalid AST expression")
}

type statement struct {
	// kind is echo, expression, block, if, class, or return.
	kind                   string
	expressions            []expr
	condition              expr
	conditionLine          int
	thenBranch             []statement
	elseBranch             []statement
	class                  *classDecl
	function               *methodDecl
	initial, update        []expr
	keyTarget, valueTarget expr
	foreachByRef           bool
	line                   int
	cases                  []switchCase
	constant               *constantDecl
	names                  []string
	catches                []catchClause
	finally                []statement
}

type catchClause struct {
	types []string
	name  string
	body  []statement
}

type parser struct {
	tokens              []token
	pos                 int
	version             Version
	filename, namespace string
	strictTypes         bool
	seenStrictTypes     bool
	classDepth          int
	currentClass        *classDecl
	currentCallable     string
	anonymousClasses    int
	imports             map[string]map[string]string
	declaredNames       map[string]map[string]bool
}

func (p *parser) current() token { return p.tokens[p.pos] }
func (p *parser) advance() token {
	t := p.current()
	if t.kind != "eof" {
		p.pos++
	}
	return t
}
func (p *parser) accept(kind string) bool {
	if p.current().kind != kind {
		return false
	}
	p.advance()
	return true
}

func (p *parser) acceptKeyword(keyword string) bool {
	if p.current().kind != "name" || !strings.EqualFold(p.current().text, keyword) {
		return false
	}
	p.advance()
	return true
}
func (p *parser) expect(kind string, alternatives ...string) error {
	if p.accept(kind) {
		return nil
	}
	return unexpected(p.current(), append(alternatives, kind)...)
}

func unexpected(t token, expected ...string) error {
	message := "syntax error, unexpected token " + strconv.Quote(t.kind)
	if t.kind == "eof" {
		message = "syntax error, unexpected end of file"
	}
	if len(expected) != 0 {
		quoted := make([]string, len(expected))
		for i, kind := range expected {
			quoted[i] = strconv.Quote(kind)
		}
		message += ", expecting " + strings.Join(quoted, " or ")
	}
	return syntax(t.line, "%s", message)
}

func parse(tokens []token) ([]statement, error) {
	return parseWithVersion(tokens, PHP84)
}

func parseWithVersion(tokens []token, version Version) ([]statement, error) {
	p := parser{tokens: tokens, version: version}
	return p.program()
}

func (p *parser) statements(end string) ([]statement, error) {
	var statements []statement
	for p.current().kind != end {
		if p.current().kind == "eof" {
			return nil, p.expect(end)
		}
		if p.accept(";") {
			continue
		}
		s, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, s)
	}
	return statements, nil
}

func (p *parser) statement() (statement, error) {
	if p.current().kind == "#[" {
		attributes, err := p.attributes()
		if err != nil {
			return statement{}, err
		}
		decl, err := p.statement()
		if err != nil {
			return decl, err
		}
		switch decl.kind {
		case "class":
			decl.class.attributes = attributes
		case "function":
			decl.function.attributes = attributes
		default:
			return statement{}, syntax(p.current().line, "Attributes can only be applied to declarations")
		}
		return decl, nil
	}
	if p.current().kind == "name" && p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].kind == ":" {
		label := p.advance()
		p.advance()
		return statement{kind: "label", names: []string{label.text}, line: label.line}, nil
	}
	if t := p.current(); p.acceptKeyword("goto") {
		label := p.current()
		if err := p.expect("name"); err != nil {
			return statement{}, err
		}
		return statement{kind: "goto", names: []string{label.text}, line: t.line}, p.expect(";")
	}
	if t := p.current(); p.acceptKeyword("try") {
		body, err := p.block()
		if err != nil {
			return statement{}, err
		}
		s := statement{kind: "try", line: t.line, thenBranch: body}
		for p.acceptKeyword("catch") {
			if err := p.expect("("); err != nil {
				return s, err
			}
			clause := catchClause{}
			for {
				typeToken, err := p.readName()
				if err != nil {
					return s, err
				}
				resolved, _ := p.resolveName(typeToken.text, "class")
				clause.types = append(clause.types, resolved)
				if !p.accept("|") {
					break
				}
			}
			if p.current().kind == "variable" {
				clause.name = p.advance().text
			}
			if err := p.expect(")"); err != nil {
				return s, err
			}
			clause.body, err = p.block()
			if err != nil {
				return s, err
			}
			s.catches = append(s.catches, clause)
		}
		hasFinally := p.acceptKeyword("finally")
		if hasFinally {
			s.finally, err = p.block()
			if err != nil {
				return s, err
			}
		}
		if len(s.catches) == 0 && !hasFinally {
			return s, syntax(t.line, "try must have a catch or finally block")
		}
		return s, nil
	}
	if p.keyword("namespace") && p.tokens[p.pos+1].kind != "\\" {
		return statement{}, &fatalDiagnostic{p.current().line, "Namespace declarations cannot be nested"}
	}
	if p.keyword("use") || p.keyword("const") {
		return statement{}, unexpected(p.current())
	}
	if t := p.current(); p.acceptKeyword("switch") {
		return p.switchStatement(t.line)
	}
	if t := p.current(); p.acceptKeyword("unset") {
		args, err := p.languageArguments("unset")
		if err != nil {
			return statement{}, err
		}
		return statement{kind: "unset", expressions: args, line: t.line}, p.expect(";")
	}
	if t := p.current(); p.acceptKeyword("foreach") {
		return p.foreachStatement(t.line)
	}
	if p.acceptKeyword("function") {
		decl, err := p.functionDeclaration()
		return statement{kind: "function", function: decl}, err
	}
	for _, kind := range []string{"while", "for", "do", "break", "continue"} {
		t := p.current()
		if p.acceptKeyword(kind) {
			return p.loopStatement(kind, t.line)
		}
	}
	if p.acceptKeyword("class") {
		decl, err := p.classDeclaration("class")
		return statement{kind: "class", class: decl}, err
	}
	if p.acceptKeyword("interface") {
		decl, err := p.classDeclaration("interface")
		return statement{kind: "class", class: decl}, err
	}
	if p.acceptKeyword("trait") {
		decl, err := p.classDeclaration("trait")
		return statement{kind: "class", class: decl}, err
	}
	if p.acceptKeyword("enum") {
		decl, err := p.classDeclaration("enum")
		return statement{kind: "class", class: decl}, err
	}
	if p.acceptKeyword("final") {
		readonly := p.acceptKeyword("readonly")
		if !p.acceptKeyword("class") {
			return statement{}, unexpected(p.current())
		}
		decl, err := p.classDeclaration("class")
		if decl != nil {
			decl.final = true
			decl.readonly = readonly
			if readonly && err == nil {
				err = validateReadonlyClassDeclaration(decl)
			}
		}
		return statement{kind: "class", class: decl}, err
	}
	if p.acceptKeyword("readonly") {
		final := p.acceptKeyword("final")
		abstract := false
		if !final {
			abstract = p.acceptKeyword("abstract")
		}
		if !p.acceptKeyword("class") {
			return statement{}, unexpected(p.current())
		}
		decl, err := p.classDeclaration("class")
		if decl != nil {
			decl.readonly = true
			decl.final = final
			decl.abstract = abstract
			if err == nil {
				err = validateReadonlyClassDeclaration(decl)
			}
		}
		return statement{kind: "class", class: decl}, err
	}
	if p.acceptKeyword("abstract") {
		readonly := p.acceptKeyword("readonly")
		if !p.acceptKeyword("class") {
			return statement{}, unexpected(p.current())
		}
		decl, err := p.classDeclaration("class")
		if decl != nil {
			decl.abstract = true
			decl.readonly = readonly
			if readonly && err == nil {
				err = validateReadonlyClassDeclaration(decl)
			}
		}
		return statement{kind: "class", class: decl}, err
	}
	if p.acceptKeyword("return") {
		s := statement{kind: "return", line: p.tokens[p.pos-1].line}
		if p.accept(";") {
			return s, nil
		}
		x, err := p.expression(0, ";")
		if err != nil {
			return s, err
		}
		s.expressions = []expr{x}
		return s, p.expect(";")
	}
	if p.acceptKeyword("global") {
		s := statement{kind: "global"}
		for {
			t := p.current()
			if err := p.expect("variable"); err != nil {
				return s, err
			}
			s.names = append(s.names, t.text)
			if !p.accept(",") {
				break
			}
		}
		return s, p.expect(";")
	}
	if p.keyword("static") && p.tokens[p.pos+1].kind == "variable" {
		p.advance()
		s := statement{kind: "static", line: p.tokens[p.pos-1].line}
		for {
			name := p.current()
			if err := p.expect("variable"); err != nil {
				return s, err
			}
			s.names = append(s.names, name.text)
			var initial expr
			if p.accept("=") {
				var err error
				initial, err = p.expression(0, ",", ";")
				if err != nil {
					return s, err
				}
			}
			s.expressions = append(s.expressions, initial)
			if !p.accept(",") {
				break
			}
		}
		return s, p.expect(";")
	}
	if p.accept(";") {
		return statement{kind: "expression"}, nil
	}
	if p.current().kind == "{" {
		body, err := p.block()
		return statement{kind: "block", thenBranch: body}, err
	}
	if p.acceptKeyword("if") {
		return p.ifStatement()
	}
	s := statement{kind: "expression"}
	if p.acceptKeyword("echo") {
		s.kind = "echo"
	}
	end := []string{";"}
	if s.kind == "echo" {
		end = []string{",", ";"}
	}
	for {
		e, err := p.expression(0, end...)
		if err != nil {
			return s, err
		}
		s.expressions = append(s.expressions, e)
		if s.kind != "echo" || !p.accept(",") {
			break
		}
	}
	if s.kind == "echo" {
		// An echo expression can be followed by another comma-separated
		// expression or by the statement terminator.
		return s, p.expect(";", ",")
	}
	return s, p.expect(";")
}

func (p *parser) ifStatement() (statement, error) {
	s := statement{kind: "if"}
	if err := p.expect("("); err != nil {
		return s, err
	}
	var err error
	s.condition, err = p.expression(0, ")")
	if err != nil {
		return s, err
	}
	s.conditionLine = expressionLine(s.condition)
	if err := p.expect(")"); err != nil {
		return s, err
	}
	s.thenBranch, err = p.body()
	if err != nil {
		return s, err
	}
	// Both spellings become an if statement in the preceding else branch.
	// The existing evaluator therefore visits later conditions only as needed.
	if p.acceptKeyword("elseif") {
		var next statement
		next, err = p.ifStatement()
		s.elseBranch = []statement{next}
	} else if p.acceptKeyword("else") {
		if p.acceptKeyword("if") {
			var next statement
			next, err = p.ifStatement()
			s.elseBranch = []statement{next}
		} else {
			s.elseBranch, err = p.body()
		}
	}
	return s, err
}

func (p *parser) body() ([]statement, error) {
	if p.current().kind == "{" {
		return p.block()
	}
	s, err := p.statement()
	return []statement{s}, err
}

func (p *parser) block() ([]statement, error) {
	if err := p.expect("{"); err != nil {
		return nil, err
	}
	body, err := p.statements("}")
	if err != nil {
		return nil, err
	}
	return body, p.expect("}")
}

const (
	keywordOrPrecedence = iota + 1
	keywordXorPrecedence
	keywordAndPrecedence
	assignmentPrecedence
	ternaryPrecedence
	coalescePrecedence
	orPrecedence
	andPrecedence
	bitOrPrecedence
	bitXorPrecedence
	bitAndPrecedence
	comparisonPrecedence
	orderingPrecedence
	concatPrecedence
	shiftPrecedence
	addPrecedence
	multiplyPrecedence
	notPrecedence
	instanceofPrecedence
	unaryPrecedence
	exponentPrecedence
)

func precedence(kind string) int {
	switch kind {
	case "or":
		return keywordOrPrecedence
	case "xor":
		return keywordXorPrecedence
	case "and":
		return keywordAndPrecedence
	case "=", "+=", "-=", "*=", "/=", "%=", ".=", "??=", "&=", "|=", "^=", "<<=", ">>=":
		return assignmentPrecedence
	case "?":
		return ternaryPrecedence
	case "??":
		return coalescePrecedence
	case "|":
		return bitOrPrecedence
	case "^":
		return bitXorPrecedence
	case "&":
		return bitAndPrecedence
	case "<<", ">>":
		return shiftPrecedence
	case "||":
		return orPrecedence
	case "&&":
		return andPrecedence
	case "===", "!==", "==", "!=", "<>", "<=>":
		return comparisonPrecedence
	case "<", "<=", ">", ">=":
		return orderingPrecedence
	case ".":
		return concatPrecedence
	case "+", "-":
		return addPrecedence
	case "*", "/", "%":
		return multiplyPrecedence
	case "**":
		return exponentPrecedence
	default:
		return 0
	}
}

func (p *parser) expression(min int, end ...string) (expr, error) {
	t := p.advance()
	var left expr
	switch t.kind {
	case "number":
		left = literal{parseLiteralNumber(t.text), t.line}
	case "string":
		left = literal{value{kind: stringKind, s: t.text}, t.line}
	case "interpolated":
		var err error
		left, err = p.interpolated(t)
		if err != nil {
			return nil, err
		}
	case "variable":
		left = variable{t}
	case "$":
		nameExpr, err := p.expression(unaryPrecedence, end...)
		if err != nil {
			return nil, err
		}
		left = dynamicVariable{nameExpr: nameExpr, line: t.line}
	case "[":
		var err error
		left, err = p.arrayLiteral("]")
		if err != nil {
			return nil, err
		}
	case "name", "\\":
		var nameErr error
		t, nameErr = p.nameAfter(t)
		if nameErr != nil {
			return nil, nameErr
		}
		if p.accept("::") {
			name := t.text
			if !equalName(name, "self") && !equalName(name, "parent") && !equalName(name, "static") {
				name, _ = p.resolveName(name, "class")
			}
			left, nameErr = p.staticMember(name, nil)
			if nameErr != nil {
				return nil, nameErr
			}
			break
		}
		if (t.text == "COUNT_NORMAL" || t.text == "COUNT_RECURSIVE") && p.namespace == "" && p.imports["const"][t.text] == "" {
			n := int64(0)
			if t.text == "COUNT_RECURSIVE" {
				n = 1
			}
			left = literal{integer(n), t.line}
			break
		}
		switch strings.ToLower(t.text) {
		case "clone":
			operand, err := p.expression(unaryPrecedence, end...)
			if err != nil {
				return nil, err
			}
			left = cloneExpr{operand: operand, line: t.line}
		case "yield":
			node := yieldExpr{line: t.line}
			if p.acceptKeyword("from") {
				node.from = true
				var err error
				node.item, err = p.expression(assignmentPrecedence, end...)
				if err != nil {
					return nil, err
				}
			} else if p.current().kind != ";" && p.current().kind != ")" && p.current().kind != "," {
				var err error
				node.item, err = p.expression(assignmentPrecedence, end...)
				if err != nil {
					return nil, err
				}
				if p.accept("=>") {
					node.key = node.item
					node.item, err = p.expression(assignmentPrecedence, end...)
					if err != nil {
						return nil, err
					}
				}
			}
			left = node
		case "throw":
			operand, err := p.expression(assignmentPrecedence, end...)
			if err != nil {
				return nil, err
			}
			left = throwExpr{operand: operand, line: t.line}
		case "__namespace__":
			left = literal{value{kind: stringKind, s: p.namespace}, t.line}
		case "__line__":
			left = literal{integer(int64(t.line)), t.line}
		case "__function__":
			left = literal{value{kind: stringKind, s: p.currentCallable}, t.line}
		case "__method__":
			name := p.currentCallable
			if p.currentClass != nil && name != "" {
				name = p.currentClass.name + "::" + name
			}
			left = literal{value{kind: stringKind, s: name}, t.line}
		case "__class__":
			name := ""
			if p.currentClass != nil {
				name = p.currentClass.name
			}
			left = literal{value{kind: stringKind, s: name}, t.line}
		case "__trait__":
			name := ""
			if p.currentClass != nil && p.currentClass.kind == "trait" {
				name = p.currentClass.name
			}
			left = literal{value{kind: stringKind, s: name}, t.line}
		case "__file__":
			filename := p.filename
			if canonical, err := filepath.EvalSymlinks(filename); err == nil {
				filename = canonical
			}
			left = literal{value{kind: stringKind, s: filename}, t.line}
		case "__dir__":
			filename := p.filename
			if canonical, err := filepath.EvalSymlinks(filename); err == nil {
				filename = canonical
			}
			left = literal{value{kind: stringKind, s: filepath.Dir(filename)}, t.line}
		case "include", "require", "include_once", "require_once":
			path, err := p.expression(assignmentPrecedence, end...)
			if err != nil {
				return nil, err
			}
			left = includeExpr{kind: strings.ToLower(t.text), path: path, line: t.line}
		case "match":
			var err error
			left, err = p.matchExpression(t.line)
			if err != nil {
				return nil, err
			}
		case "isset", "empty":
			args, err := p.languageArguments(strings.ToLower(t.text))
			if err != nil {
				return nil, err
			}
			left = languageExpr{strings.ToLower(t.text), args, t.line}
		case "true":
			left = literal{value{kind: boolKind, i: 1}, t.line}
		case "false":
			left = literal{value{kind: boolKind}, t.line}
		case "null":
			left = literal{value{kind: nullKind}, t.line}
		case "array":
			if err := p.expect("("); err != nil {
				return nil, err
			}
			var err error
			left, err = p.arrayLiteral(")")
			if err != nil {
				return nil, err
			}
		case "new":
			name := p.current()
			if p.acceptKeyword("class") {
				arguments := p.accept("(")
				var args []expr
				if arguments {
					var err error
					args, err = p.arguments()
					if err != nil {
						return nil, err
					}
				}
				p.anonymousClasses++
				generated := fmt.Sprintf("class@anonymous:%s:%d:%d", p.filename, name.line, p.anonymousClasses)
				decl, err := p.classDeclarationNamed("class", token{kind: "name", text: generated, line: name.line})
				if err != nil {
					return nil, err
				}
				left = newObject{name: decl.name, anonymous: decl, line: t.line, arguments: arguments, args: args}
				break
			}
			var classExpr expr
			if name.kind == "(" {
				p.advance()
				var err error
				classExpr, err = p.expression(0, ")")
				if err != nil {
					return nil, err
				}
				if err := p.expect(")"); err != nil {
					return nil, err
				}
			} else if name.kind == "variable" {
				p.advance()
				classExpr = variable{name}
				for {
					if p.accept("->") {
						member := p.current()
						if member.kind != "name" && member.kind != "variable" {
							return nil, unexpected(member, "name", "variable")
						}
						p.advance()
						property := propertyExpr{base: classExpr, line: member.line}
						if member.kind == "name" {
							property.name = member.text
						} else {
							property.nameExpr = variable{member}
						}
						classExpr = property
						continue
					}
					if p.accept("[") {
						key, keyErr := p.expression(0, "]")
						if keyErr != nil {
							return nil, keyErr
						}
						if err := p.expect("]"); err != nil {
							return nil, err
						}
						classExpr = offsetExpr{base: classExpr, key: key, line: name.line}
						continue
					}
					break
				}
			} else {
				var err error
				name, err = p.readName()
				if err != nil {
					return nil, err
				}
				if (equalName(name.text, "self") || equalName(name.text, "parent") || equalName(name.text, "static")) && p.classDepth == 0 {
					return nil, &fatalDiagnostic{name.line, "Cannot use " + strings.ToLower(name.text) + " when no class scope is active"}
				}
				if p.accept("::") {
					className := name.text
					if !equalName(className, "self") && !equalName(className, "parent") && !equalName(className, "static") {
						className, _ = p.resolveName(className, "class")
					}
					if p.current().kind == "variable" {
						member := p.current()
						p.advance()
						classExpr = staticPropertyExpr{className: className, name: member.text, line: member.line}
					} else {
						classExpr, err = p.staticMember(className, nil)
						if err != nil {
							return nil, err
						}
					}
					for p.accept("[") {
						key, keyErr := p.expression(0, "]")
						if keyErr != nil {
							return nil, keyErr
						}
						if err := p.expect("]"); err != nil {
							return nil, err
						}
						classExpr = offsetExpr{base: classExpr, key: key, line: name.line}
					}
				}
			}
			arguments := p.accept("(")
			var args []expr
			if arguments {
				var err error
				args, err = p.arguments()
				if err != nil {
					return nil, err
				}
			}
			resolved := name.text
			if !equalName(resolved, "self") && !equalName(resolved, "parent") && !equalName(resolved, "static") {
				resolved, _ = p.resolveName(name.text, "class")
			}
			left = newObject{name: resolved, classExpr: classExpr, line: name.line, arguments: arguments, args: args}
		case "function":
			var err error
			left, err = p.anonymousFunction(t.line, false)
			if err != nil {
				return nil, err
			}
		case "fn":
			var err error
			left, err = p.arrowFunction(t.line, false, end...)
			if err != nil {
				return nil, err
			}
		case "list":
			if err := p.expect("("); err != nil {
				return nil, err
			}
			var err error
			left, err = p.arrayLiteral(")")
			if err != nil {
				return nil, err
			}
		case "static":
			if p.acceptKeyword("fn") {
				var err error
				left, err = p.arrowFunction(t.line, true, end...)
				if err != nil {
					return nil, err
				}
				break
			}
			if !p.acceptKeyword("function") {
				return nil, syntax(t.line, "unsupported keyword, constant, or function %q", t.text)
			}
			var err error
			left, err = p.anonymousFunction(t.line, true)
			if err != nil {
				return nil, err
			}
		default:
			if strings.EqualFold(t.text, "eval") && p.accept("(") {
				args, err := p.arguments()
				if err != nil {
					return nil, err
				}
				left = functionCall{name: "eval", args: args, line: t.line}
				break
			}
			if strings.EqualFold(t.text, "exit") || strings.EqualFold(t.text, "die") {
				var args []expr
				if p.accept("(") {
					var err error
					args, err = p.arguments()
					if err != nil {
						return nil, err
					}
				}
				left = functionCall{name: strings.ToLower(t.text), args: args, line: t.line}
				break
			}
			if p.accept("(") && !reservedName(t.text) {
				if p.firstClassArgument() {
					name, fallback := p.resolveName(t.text, "function")
					left = firstClassCallable{callee: functionCall{name: name, line: t.line, fallback: fallback}, line: t.line}
					break
				}
				args, err := p.arguments()
				if err != nil {
					return nil, err
				}
				name, fallback := p.resolveName(t.text, "function")
				left = functionCall{name: name, args: args, line: t.line, fallback: fallback}
				break
			}
			if reservedName(t.text) {
				return nil, syntax(t.line, "unsupported keyword, constant, or function %q", t.text)
			}
			name, fallback := p.resolveName(t.text, "const")
			left = constantExpr{name, fallback, t.line}
		}
	case "++", "--":
		target, err := p.expression(unaryPrecedence+1, end...)
		if err != nil {
			return nil, err
		}
		if !writable(target) {
			return nil, syntax(t.line, "invalid increment/decrement target")
		}
		left = increment{target, t, false}
	case "+", "-", "!", "~", "@":
		prec := unaryPrecedence
		if t.kind == "!" {
			prec = notPrecedence
		}
		operand, err := p.expression(prec, end...)
		if err != nil {
			return nil, err
		}
		left = unary{t, operand}
	case "(":
		if p.current().kind == "name" && p.tokens[p.pos+1].kind == ")" {
			name := strings.ToLower(p.current().text)
			switch name {
			case "bool", "int", "float", "string", "array", "object":
				p.advance()
				p.advance()
				operand, err := p.expression(unaryPrecedence, end...)
				if err != nil {
					return nil, err
				}
				left = castExpr{name, operand, t.line}
			case "boolean", "integer", "double", "real", "binary", "unset":
				return nil, syntax(t.line, "unsupported cast (%s)", name)
			}
			if left != nil {
				break
			}
		}
		var err error
		left, err = p.expression(0, ")")
		if err != nil {
			return nil, err
		}
		if err := p.expect(")"); err != nil {
			return nil, err
		}
		if node, ok := left.(newObject); ok {
			node.grouped = true
			left = node
		}
		if node, ok := left.(propertyExpr); ok {
			node.grouped = true
			left = node
		}
		if node, ok := left.(ternaryExpr); ok {
			node.grouped = true
			left = node
		}
	default:
		if t.kind == "eof" {
			return nil, syntax(t.line, "syntax error, unexpected end of file")
		}
		return nil, syntax(t.line, "syntax error, unexpected token %q", t.kind)
	}
	compared := 0
	for {
		op := p.current()
		if op.kind == "name" {
			switch strings.ToLower(op.text) {
			case "and", "or", "xor":
				op.kind = strings.ToLower(op.text)
			}
		}
		if op.kind == "::" {
			p.advance()
			var err error
			left, err = p.staticMember("", left)
			if err != nil {
				return nil, err
			}
			continue
		}
		if op.kind == "name" && strings.EqualFold(op.text, "instanceof") {
			if instanceofPrecedence < min {
				return left, nil
			}
			p.advance()
			if p.current().kind == "variable" {
				name := p.advance()
				left = instanceofExpr{operand: left, classExpr: variable{name}, line: op.line}
				continue
			}
			name, err := p.readName()
			if err != nil {
				return nil, err
			}
			resolved, _ := p.resolveName(name.text, "class")
			if equalName(name.text, "self") || equalName(name.text, "parent") || equalName(name.text, "static") {
				if p.classDepth == 0 {
					return nil, &fatalDiagnostic{name.line, "Cannot use " + strings.ToLower(name.text) + " when no class scope is active"}
				}
				resolved = name.text
			}
			left = instanceofExpr{operand: left, className: resolved, line: op.line}
			continue
		}
		if property, ok := left.(nullsafeExpr); ok && !property.method && p.accept("(") {
			args, err := p.arguments()
			if err != nil {
				return nil, err
			}
			property.method, property.args = true, args
			left = property
			continue
		}
		if property, ok := left.(propertyExpr); ok && !property.grouped && p.accept("(") {
			if p.firstClassArgument() {
				left = firstClassCallable{callee: methodCall{receiver: property.base, name: property.name, nameExpr: property.nameExpr, line: property.line}, line: property.line}
				continue
			}
			args, err := p.arguments()
			if err != nil {
				return nil, err
			}
			left = methodCall{receiver: property.base, name: property.name, nameExpr: property.nameExpr, args: args, line: property.line}
			continue
		}
		if p.accept("(") {
			if p.firstClassArgument() {
				left = firstClassCallable{callee: left, line: op.line}
				continue
			}
			args, err := p.arguments()
			if err != nil {
				return nil, err
			}
			left = dynamicCall{callee: left, args: args, line: op.line}
			continue
		}
		if node, ok := left.(newObject); ok && !node.grouped && (op.kind == "[" || op.kind == "->" || op.kind == "?->") && (!node.arguments || p.version == PHP83) {
			// The new expression is complete in this target's grammar. Report
			// the enclosing expression's terminators (e.g. comma/semicolon for echo).
			return nil, unexpected(op, end...)
		}
		if p.accept("[") {
			var key expr
			var err error
			if p.current().kind != "]" {
				key, err = p.expression(0, "]")
			}
			if err != nil {
				return nil, err
			}
			line := p.current().line
			if err := p.expect("]"); err != nil {
				return nil, err
			}
			left = offsetExpr{left, key, line}
			continue
		}
		if p.accept("->") || p.accept("?->") {
			member := p.current()
			var name string
			var nameExpr expr
			switch member.kind {
			case "name":
				name = p.advance().text
			case "variable":
				nameExpr = variable{p.advance()}
			case "{":
				p.advance()
				var err error
				nameExpr, err = p.expression(0, "}")
				if err != nil {
					return nil, err
				}
				if err := p.expect("}"); err != nil {
					return nil, err
				}
			default:
				return nil, unexpected(member, "name", "variable", "{")
			}
			if op.kind == "?->" {
				left = nullsafeExpr{base: left, name: name, line: member.line, nameExpr: nameExpr}
			} else {
				left = propertyExpr{base: left, name: name, line: member.line, nameExpr: nameExpr}
			}
			continue
		}
		if (op.kind == "++" || op.kind == "--") && min <= unaryPrecedence {
			if !writable(left) {
				return nil, unexpected(op, end...)
			}
			p.advance()
			left = increment{left, op, true}
			continue
		}
		if precedence(op.kind) == assignmentPrecedence && writable(left) {
			p.advance()
			if op.kind == "=" && p.current().kind == "&" {
				p.advance()
				right, err := p.expression(assignmentPrecedence, end...)
				if err != nil {
					return nil, err
				}
				if !referenceTargetSyntax(left) || !referenceSourceSyntax(right) {
					return nil, syntax(op.line, "gophp: unsupported reference assignment target")
				}
				left = referenceAssignment{left, right, op.line}
				continue
			}
			right, err := p.expression(assignmentPrecedence, end...)
			if err != nil {
				return nil, err
			}
			if op.kind == "??=" {
				left = coalesceAssignment{left, right, op.line}
			} else if op.kind == "=" {
				left = assignment{left, right, op.line}
			} else {
				left = compoundAssignment{left, right, op}
			}
			continue
		}
		prec := precedence(op.kind)
		if prec == 0 || prec < min {
			return left, nil
		}
		if prec == assignmentPrecedence {
			return nil, syntax(op.line, "invalid assignment target")
		}
		if op.kind == "?" {
			p.advance()
			var yes expr
			var err error
			if p.current().kind != ":" {
				yes, err = p.expression(0, ":")
				if err != nil {
					return nil, err
				}
			}
			if previous, ok := left.(ternaryExpr); ok && !previous.grouped && (previous.yes != nil || yes != nil) {
				form := "a ? b : c ? d : e"
				alternatives := "(a ? b : c) ? d : e` or `a ? b : (c ? d : e)"
				if previous.yes == nil {
					form = "a ?: b ? c : d"
					alternatives = "(a ?: b) ? c : d` or `a ?: (b ? c : d)"
				}
				if yes == nil {
					form = "a ? b : c ?: d"
					alternatives = "(a ? b : c) ?: d` or `a ? b : (c ?: d)"
				}
				return nil, &fatalDiagnostic{op.line, "Unparenthesized `" + form + "` is not supported. Use either `" + alternatives + "`"}
			}
			if err = p.expect(":"); err != nil {
				return nil, err
			}
			no, err := p.expression(ternaryPrecedence+1, end...)
			if err != nil {
				return nil, err
			}
			left = ternaryExpr{condition: left, yes: yes, no: no, line: op.line}
			continue
		}
		if prec == comparisonPrecedence || prec == orderingPrecedence {
			// Equal-precedence comparisons are non-associative in PHP.
			// A parenthesized expression is parsed by a separate call, so
			// explicit grouping still permits (a === b) === c.
			if compared&(1<<prec) != 0 {
				return nil, syntax(op.line, "syntax error, unexpected token %q", op.text)
			}
			compared |= 1 << prec
		}
		p.advance()
		rightPrec := prec + 1
		if op.kind == "??" || op.kind == "**" {
			rightPrec = prec
		}
		right, err := p.expression(rightPrec, end...)
		if err != nil {
			return nil, err
		}
		left = binary{op, left, right}
	}
}

// The scope token has been consumed. A bare $name denotes a property on a
// read, but denotes a variable method name when followed by a call.
func (p *parser) staticMember(className string, classExpr expr) (expr, error) {
	member := p.current()
	if member.kind == "name" && equalName(member.text, "class") {
		p.advance()
		if classExpr == nil {
			return classNameExpr{name: className, line: member.line}, nil
		}
		return classNameExpr{line: member.line, dynamic: classExpr}, nil
	}
	var name string
	var nameExpr expr
	property := false
	switch member.kind {
	case "name":
		p.advance()
		name = member.text
	case "variable":
		p.advance()
		name = member.text
		property = true
		if p.current().kind == "(" {
			nameExpr = variable{member}
		}
	case "{":
		p.advance()
		var err error
		nameExpr, err = p.expression(0, "}")
		if err != nil {
			return nil, err
		}
		if err = p.expect("}"); err != nil {
			return nil, err
		}
	case "$":
		p.advance()
		if err := p.expect("{"); err != nil {
			return nil, err
		}
		var err error
		nameExpr, err = p.expression(0, "}")
		if err != nil {
			return nil, err
		}
		if err = p.expect("}"); err != nil {
			return nil, err
		}
		property = true
	default:
		return nil, unexpected(member, "name", "variable")
	}
	if p.accept("(") {
		if p.firstClassArgument() {
			return firstClassCallable{callee: staticCallExpr{className: className, name: name, line: member.line, classExpr: classExpr, nameExpr: nameExpr}, line: member.line}, nil
		}
		args, err := p.arguments()
		if err != nil {
			return nil, err
		}
		return staticCallExpr{className: className, name: name, args: args, line: member.line, classExpr: classExpr, nameExpr: nameExpr}, nil
	}
	if property {
		return staticPropertyExpr{className: className, name: name, line: member.line, classExpr: classExpr, nameExpr: nameExpr}, nil
	}
	if nameExpr != nil {
		return nil, unexpected(p.current(), "(")
	}
	return classConstantExpr{className, name, member.line}, nil
}

// A sole ellipsis is the PHP 8.1+ first-class callable syntax, not an
// unpacked argument. The opening parenthesis has already been consumed.
func (p *parser) firstClassArgument() bool {
	if p.current().kind != "..." || p.pos+1 >= len(p.tokens) || p.tokens[p.pos+1].kind != ")" {
		return false
	}
	p.advance()
	p.advance()
	return true
}

// The opening parenthesis has already been consumed.
func (p *parser) arguments() ([]expr, error) {
	var args []expr
	seenUnpack, seenNamed := false, false
	for p.current().kind != ")" {
		start := p.current()
		unpack := p.accept("...")
		named := !unpack && start.kind == "name" && p.tokens[p.pos+1].kind == ":"
		if unpack && seenNamed {
			return nil, &fatalDiagnostic{start.line, "Cannot use argument unpacking after named arguments"}
		}
		if !unpack && !named {
			if seenUnpack {
				return nil, &fatalDiagnostic{start.line, "Cannot use positional argument after argument unpacking"}
			}
			if seenNamed {
				return nil, &fatalDiagnostic{start.line, "Cannot use positional argument after named argument"}
			}
		}
		if named {
			p.advance()
			p.advance()
			seenNamed = true
		}
		x, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		if unpack {
			x = unpackArgument{operand: x, line: start.line}
			seenUnpack = true
		} else if named {
			x = namedArgument{name: start.text, operand: x, line: start.line}
		}
		args = append(args, x)
		if !p.accept(",") {
			break
		}
	}
	return args, p.expect(")")
}

func writable(e expr) bool {
	// Keep nullsafe paths in the write AST so compile validation produces a
	// fatal diagnostic even in an unselected branch.
	if hasNullsafe(e) {
		return true
	}
	switch node := e.(type) {
	case variable:
		return true
	case dynamicVariable:
		return true
	case staticPropertyExpr:
		return true
	case offsetExpr:
		if writable(node.base) {
			return true
		}
		// PHP permits an offset write through a call result. Whether the
		// returned value is writable is determined when the expression runs.
		switch node.base.(type) {
		case functionCall, methodCall, dynamicCall:
			return true
		}
		return false
	case propertyExpr:
		return true
	case arrayLiteral:
		for _, element := range node.elements {
			if element.value != nil && !writable(element.value) {
				return false
			}
		}
		return len(node.elements) != 0
	}
	return false
}

func (p *parser) arrayLiteral(end string) (expr, error) {
	node := arrayLiteral{}
	for p.current().kind != end {
		if p.accept(",") {
			node.elements = append(node.elements, arrayElement{})
			continue
		}
		// At the start of an optional element (including after a comma),
		// the closing delimiter is expected if the statement ends here.
		if t := p.current(); t.kind == ";" || t.kind == "eof" {
			return nil, unexpected(t, end)
		}
		unpack := p.accept("...")
		byRef := p.accept("&")
		v, err := p.expression(0)
		if err != nil {
			return nil, err
		}
		if unpack {
			node.elements = append(node.elements, arrayElement{value: unpackArgument{operand: v, line: expressionLine(v)}})
			if !p.accept(",") {
				break
			}
			continue
		}
		entry := arrayElement{value: v, byRef: byRef}
		if p.accept("=>") {
			entry.key = v
			entry.byRef = p.accept("&")
			entry.value, err = p.expression(0)
			if err != nil {
				return nil, err
			}
		}
		node.elements = append(node.elements, entry)
		if !p.accept(",") {
			break
		}
	}
	node.line = p.current().line
	return node, p.expect(end)
}
