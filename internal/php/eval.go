package php

import (
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"
)

type runtimeError struct {
	line    int
	class   string
	message string
}

func (e *runtimeError) Error() string { return e.message }

type evaluator struct {
	runtimeSymbols
	version               Version
	filename              string
	strictTypes           bool
	stdout, stderr        io.Writer
	variables             variableScope
	rootScope             variableScope
	bindings              map[string]bindingKind
	staticScope           variableScope
	staticScopes          map[*methodDecl]variableScope
	currentMethod         *methodDecl
	predeclared           map[*classDecl]bool
	predeclaredFunctions  map[*methodDecl]bool
	bounded               bool
	steps                 int
	scope                 *classDecl
	lateClass             *classDecl
	frames, failureFrames []callFrame
	loadedFiles           map[string]bool
	includeDepth          int
}

// Run parses the entire source before executing it, so a syntax error cannot
// leave partially executed statements. Each invocation has a fresh scope.
func Run(source, filename string, stdout, stderr io.Writer) int {
	return RunWithVersion(source, filename, stdout, stderr, PHP84)
}

// RunWithVersion keeps the same execution pipeline for all supported versions.
func RunWithVersion(source, filename string, stdout, stderr io.Writer, version Version) int {
	return run(source, filename, stdout, stderr, version, 0)
}

func run(source, filename string, stdout, stderr io.Writer, version Version, budget int) int {
	if _, err := ParseVersion(string(version)); err != nil {
		fmt.Fprintf(stderr, "gophp: %v\n", err)
		return 2
	}
	tokens, err := lex(source)
	var statements []statement
	strictTypes := false
	if err == nil {
		p := parser{tokens: tokens, version: version, filename: filename}
		statements, err = p.program()
		strictTypes = p.strictTypes
	}
	if err != nil {
		if failure, ok := err.(*fatalDiagnostic); ok {
			fmt.Fprintf(stderr, "Fatal error: %s in %s on line %d\n", failure.message, filename, failure.line)
			if version == PHP85 {
				fmt.Fprint(stderr, "Stack trace:\n#0 {main}\n")
			}
			return 255
		}
		e := err.(*syntaxError)
		fmt.Fprintf(stderr, "Parse error: %s in %s on line %d\n", e.message, filename, e.line)
		return 255
	}
	failure := validateStatements(statements)
	if failure == nil {
		failure = duplicateConstants(statements)
	}
	if failure == nil {
		failure = duplicateFunctions(statements, filename, version)
	}
	if failure != nil {
		fmt.Fprintf(stderr, "Fatal error: %s in %s on line %d\n", failure.message, filename, failure.line)
		if version == PHP85 {
			// PHP 8.5 includes a trace for compile errors, even though no
			// statement or method has executed. Parse errors remain trace-free.
			fmt.Fprint(stderr, "Stack trace:\n#0 {main}\n")
		}
		return 255
	}
	e := evaluator{version: version, filename: filename, strictTypes: strictTypes, stdout: stdout, stderr: stderr, variables: make(map[string]value), bounded: budget > 0, steps: budget, loadedFiles: make(map[string]bool)}
	e.rootScope = e.variables
	if absolute, pathErr := filepath.Abs(filename); pathErr == nil {
		if canonical, canonicalErr := filepath.EvalSymlinks(absolute); canonicalErr == nil {
			absolute = canonical
		}
		e.loadedFiles[absolute] = true
	}
	e.callableDiagnostics(statements)
	e.controlDiagnostics(statements, nil)
	if err := e.executeProgram(statements); err != nil {
		errorFilename := filename
		var sourced *sourceError
		if errors.As(err, &sourced) {
			errorFilename = sourced.filename
		}
		if _, ok := err.(*returned); ok {
			return 0
		}
		var parseFailure *syntaxError
		if errors.As(err, &parseFailure) {
			fmt.Fprintf(stderr, "Parse error: %s in %s on line %d\n", parseFailure.message, errorFilename, parseFailure.line)
			return 255
		}
		var fatal *fatalDiagnostic
		if errors.As(err, &fatal) {
			fmt.Fprintf(stderr, "Fatal error: %s in %s on line %d\n", fatal.message, errorFilename, fatal.line)
			if version == PHP85 {
				fmt.Fprintf(stderr, "Stack trace:\n%s", e.stackTrace())
			}
			return 255
		}
		var failure *runtimeError
		var raised *thrown
		if errors.As(err, &raised) {
			data := raised.object.o.exception
			fmt.Fprintf(stderr, "Fatal error: Uncaught %s: %s in %s:%d\nStack trace:\n%s  thrown in %s on line %d\n", raised.object.o.className(), data.message, data.file, data.line, data.trace, data.file, data.line)
			return 255
		}
		if errors.As(err, &failure) {
			if failure.class == "" {
				fmt.Fprintf(stderr, "Fatal error: %s in %s on line %d\n", failure.message, errorFilename, failure.line)
				return 255
			}
			fmt.Fprintf(stderr, "Fatal error: Uncaught %s: %s in %s:%d\nStack trace:\n%s  thrown in %s on line %d\n", failure.class, failure.message, errorFilename, failure.line, e.stackTrace(), errorFilename, failure.line)
			return 255
		}
		fmt.Fprintf(stderr, "gophp: %v\n", err)
		return 1
	}
	return 0
}

func (e *evaluator) executeProgram(statements []statement) error {
	oldClasses, oldFunctions := e.predeclared, e.predeclaredFunctions
	defer func() { e.predeclared, e.predeclaredFunctions = oldClasses, oldFunctions }()
	e.predeclared = make(map[*classDecl]bool)
	e.predeclaredFunctions = make(map[*methodDecl]bool)
	// Unconditional top-level classes with resolvable parents are available
	// before their declaration. A parent supplied by a preceding include is
	// resolved when execution reaches the child declaration.
	pending := make(map[string]*classDecl)
	for _, s := range statements {
		if s.kind == "class" {
			pending[strings.ToLower(s.class.name)] = s.class
		}
	}
	visiting := make(map[*classDecl]bool)
	// Only declarations whose parents are already known (or declared in this
	// file) can be bound eagerly. Executing a prior statement may register a
	// class through a function, branch, switch, or include; leave an unresolved
	// child pending until execution reaches its declaration in that case.
	canDefineBefore := make(map[*classDecl]bool)
	canDefine := false
	hasExecutedBefore := make(map[*classDecl]bool)
	hasExecutable := false
	for _, s := range statements {
		if s.kind == "class" {
			canDefineBefore[s.class] = canDefine
			hasExecutedBefore[s.class] = hasExecutable
		}
		if statementMayDefineClass(s) {
			canDefine = true
		}
		if s.kind != "class" && s.kind != "function" {
			hasExecutable = true
		}
	}
	var predeclareClass func(*classDecl) (bool, error)
	predeclareClass = func(decl *classDecl) (bool, error) {
		if e.predeclared[decl] {
			return true, nil
		}
		// Interface implementation and trait composition are bound when the
		// declaration is reached after executable code in the same file.
		if hasExecutedBefore[decl] && (len(decl.interfaceNames) != 0 || len(decl.traitUses) != 0) {
			e.initClasses()
			resolved := true
			for _, name := range decl.interfaceNames {
				key := strings.ToLower(name)
				resolved = resolved && (pending[key] != nil || e.classes[key] != nil)
			}
			for _, use := range decl.traitUses {
				for _, name := range use.names {
					key := strings.ToLower(name)
					resolved = resolved && (pending[key] != nil || e.classes[key] != nil)
				}
			}
			if resolved {
				return false, nil
			}
		}
		if visiting[decl] {
			// Let declaration execution report the cycle at its source line.
			return false, nil
		}
		for _, use := range decl.traitUses {
			for _, name := range use.names {
				traitKey := strings.ToLower(name)
				if dependency := pending[traitKey]; dependency != nil {
					visiting[decl] = true
					ready, err := predeclareClass(dependency)
					delete(visiting, decl)
					if err != nil || !ready {
						return false, err
					}
				} else {
					e.initClasses()
					if e.classes[traitKey] == nil {
						if canDefineBefore[decl] {
							return false, nil
						}
						return false, &runtimeError{decl.line, "Error", "Trait \"" + name + "\" not found"}
					}
				}
			}
		}
		for _, name := range decl.interfaceNames {
			interfaceKey := strings.ToLower(name)
			if dependency := pending[interfaceKey]; dependency != nil {
				visiting[decl] = true
				ready, err := predeclareClass(dependency)
				delete(visiting, decl)
				if err != nil || !ready {
					return false, err
				}
			} else {
				e.initClasses()
				if e.classes[interfaceKey] == nil {
					if canDefineBefore[decl] {
						return false, nil
					}
					return false, &runtimeError{decl.line, "Error", "Interface \"" + name + "\" not found"}
				}
			}
		}
		if decl.parentName != "" {
			parentKey := strings.ToLower(decl.parentName)
			if parent := pending[parentKey]; parent != nil {
				visiting[decl] = true
				ready, err := predeclareClass(parent)
				delete(visiting, decl)
				if err != nil || !ready {
					return false, err
				}
			} else {
				e.initClasses()
				if e.classes[parentKey] == nil {
					if canDefineBefore[decl] {
						return false, nil
					}
					return false, &runtimeError{decl.line, "Error", "Class \"" + decl.parentName + "\" not found"}
				}
			}
		}
		if err := e.declareClass(decl); err != nil {
			return false, err
		}
		e.predeclared[decl] = true
		return true, nil
	}
	for _, s := range statements {
		if s.kind == "function" {
			if err := e.declareFunction(s.function); err != nil {
				return err
			}
			e.predeclaredFunctions[s.function] = true
		}
		if s.kind == "class" {
			if _, err := predeclareClass(s.class); err != nil {
				return err
			}
		}
	}
	return e.execute(statements)
}

func statementMayDefineClass(s statement) bool {
	switch s.kind {
	case "class", "function", "const":
		return false // The body of a function declaration does not execute here.
	case "echo":
		for _, x := range s.expressions {
			if expressionMayDefineClass(x) {
				return true
			}
		}
		return false
	default:
		// Runtime statements are barriers even when their bodies are not
		// obviously loads: user code called from them may include a file.
		return true
	}
}

func expressionMayDefineClass(x expr) bool {
	switch n := x.(type) {
	case includeExpr:
		return true
	case assignment:
		return expressionMayDefineClass(n.target) || expressionMayDefineClass(n.source)
	case compoundAssignment:
		return expressionMayDefineClass(n.target) || expressionMayDefineClass(n.source)
	case coalesceAssignment:
		return expressionMayDefineClass(n.target) || expressionMayDefineClass(n.source)
	case increment:
		return expressionMayDefineClass(n.target)
	case binary:
		return expressionMayDefineClass(n.left) || expressionMayDefineClass(n.right)
	case unary:
		return expressionMayDefineClass(n.operand)
	case castExpr:
		return expressionMayDefineClass(n.operand)
	case throwExpr:
		return expressionMayDefineClass(n.operand)
	case ternaryExpr:
		return expressionMayDefineClass(n.condition) || expressionMayDefineClass(n.yes) || expressionMayDefineClass(n.no)
	case offsetExpr:
		return expressionMayDefineClass(n.base) || expressionMayDefineClass(n.key)
	case propertyExpr:
		return expressionMayDefineClass(n.base)
	case methodCall:
		return true
	case nullsafeExpr:
		if expressionMayDefineClass(n.base) {
			return true
		}
		for _, arg := range n.args {
			if expressionMayDefineClass(arg) {
				return true
			}
		}
	case arrayLiteral:
		for _, element := range n.elements {
			if expressionMayDefineClass(element.key) || expressionMayDefineClass(element.value) {
				return true
			}
		}
	case dynamicCall:
		return true
	case namedArgument:
		return expressionMayDefineClass(n.operand)
	case newObject:
		return true
	case languageExpr:
		for _, arg := range n.args {
			if expressionMayDefineClass(arg) {
				return true
			}
		}
	case instanceofExpr:
		return expressionMayDefineClass(n.operand)
	case matchExpr:
		if expressionMayDefineClass(n.subject) {
			return true
		}
		for _, arm := range n.arms {
			for _, condition := range arm.conditions {
				if expressionMayDefineClass(condition) {
					return true
				}
			}
			if expressionMayDefineClass(arm.result) {
				return true
			}
		}
	case functionCall:
		return true
	}
	return false
}

func (e *evaluator) execute(statements []statement) error {
	for _, statement := range statements {
		switch statement.kind {
		case "static":
			if e.currentMethod == nil {
				return &runtimeError{statement.line, "", "static variables require a function or method scope"}
			}
			if e.bindings == nil {
				e.bindings = make(map[string]bindingKind)
			}
			for i, name := range statement.names {
				if _, exists := e.staticScope[name]; !exists {
					initial := value{kind: nullKind}
					if statement.expressions[i] != nil {
						var err error
						initial, err = e.eval(statement.expressions[i])
						if err != nil {
							return err
						}
					}
					e.staticScope[name] = initial
				}
				e.bindings[name] = staticBinding
			}
			continue
		case "global":
			if e.bindings == nil {
				e.bindings = make(map[string]bindingKind)
			}
			for _, name := range statement.names {
				if _, exists := e.rootScope[name]; !exists {
					e.rootScope[name] = value{kind: nullKind}
				}
				e.bindings[name] = globalBinding
			}
			continue
		case "const":
			decl := statement.constant
			v, err := e.eval(decl.initial)
			if err != nil {
				return err
			}
			e.defineConstant(decl.name, v, decl.line)
			continue
		case "switch":
			if err := e.executeSwitch(statement); err != nil {
				return err
			}
			continue
		case "try":
			if err := e.executeTry(statement); err != nil {
				return err
			}
			continue
		case "unset":
			for _, target := range statement.expressions {
				if err := e.unset(target); err != nil {
					return err
				}
			}
			continue
		case "foreach":
			if err := e.executeForeach(statement); err != nil {
				return err
			}
			continue
		case "function":
			if !e.predeclaredFunctions[statement.function] {
				if err := e.declareFunction(statement.function); err != nil {
					return err
				}
			}
			continue
		case "while", "for", "do":
			if err := e.executeLoop(statement); err != nil {
				return err
			}
			continue
		case "break", "continue":
			level, _ := loopLevel(statement, int(^uint(0)>>1))
			return &loopControl{statement.kind, level}
		}
		if statement.kind == "class" {
			if !e.predeclared[statement.class] {
				if err := e.declareClass(statement.class); err != nil {
					return err
				}
			}
			continue
		}
		if statement.kind == "return" {
			var v value
			var err error
			if len(statement.expressions) != 0 {
				v, err = e.eval(statement.expressions[0])
			}
			if err != nil {
				return err
			}
			return &returned{v, statement.line}
		}
		if statement.kind == "block" {
			if err := e.execute(statement.thenBranch); err != nil {
				return err
			}
			continue
		}
		if statement.kind == "if" {
			condition, err := e.eval(statement.condition)
			if err != nil {
				return err
			}
			branch := statement.elseBranch
			if e.toBool(condition, statement.conditionLine) {
				branch = statement.thenBranch
			}
			// Blocks share the surrounding PHP variable scope. Unchosen
			// statements (including nested conditions) are never evaluated.
			if err := e.execute(branch); err != nil {
				return err
			}
			continue
		}
		for _, expression := range statement.expressions {
			v, err := e.eval(expression)
			if err != nil {
				return err
			}
			switch statement.kind {
			case "echo":
				text, err := e.toText(v, expressionLine(expression))
				if err != nil {
					return err
				}
				if _, err := io.WriteString(e.stdout, text); err != nil {
					return fmt.Errorf("write stdout: %w", err)
				}
			}
		}
	}
	return nil
}

func (e *evaluator) warn(line int, message string) {
	e.diagnostic("Warning", line, message)
}

func (e *evaluator) diagnostic(level string, line int, message string) {
	fmt.Fprintf(e.stderr, "%s: %s in %s on line %d\n", level, message, e.filename, line)
}

func (e *evaluator) toText(v value, line int) (string, error) {
	if e.version == PHP85 && v.kind == floatKind && math.IsNaN(v.f) {
		e.warn(line, "unexpected NAN value was coerced to string")
	}
	if v.kind == arrayKind {
		e.warn(line, "Array to string conversion")
	}
	if v.kind == objectKind {
		return "", &runtimeError{line, "Error", "Object of class " + v.o.className() + " could not be converted to string"}
	}
	return v.text(), nil
}

func (e *evaluator) toBool(v value, line int) bool {
	if e.version == PHP85 && v.kind == floatKind && math.IsNaN(v.f) {
		e.warn(line, "unexpected NAN value was coerced to bool")
	}
	return v.truthy()
}

func (e *evaluator) eval(expression expr) (value, error) {
	if hasNullsafe(expression) {
		v, _, err := e.evalChain(expression)
		return v, err
	}
	switch node := expression.(type) {
	case throwExpr:
		v, err := e.eval(node.operand)
		if err != nil {
			return value{}, err
		}
		return value{}, e.throwValue(v, node.line)
	case dynamicCall:
		callee, err := e.eval(node.callee)
		if err != nil {
			return value{}, err
		}
		if callee.kind == stringKind {
			return e.callFunction(functionCall{name: strings.TrimPrefix(callee.s, "\\"), args: node.args, line: node.line})
		}
		if callee.kind == objectKind && callee.o.closure != nil {
			return e.invokeClosure(callee.o.closure, node.args, node.line)
		}
		if callee.kind == arrayKind && len(callee.a.entries) == 2 && !callee.a.entries[0].key.isString && callee.a.entries[0].key.index == 0 && !callee.a.entries[1].key.isString && callee.a.entries[1].key.index == 1 {
			receiver, member := callee.a.entries[0].value, callee.a.entries[1].value
			if member.kind == stringKind && receiver.kind == objectKind {
				return e.callMethod(methodCall{receiver: literal{value: receiver, line: node.line}, name: member.s, args: node.args, line: node.line})
			}
			if member.kind == stringKind && receiver.kind == stringKind {
				e.initClasses()
				class := e.classes[strings.ToLower(strings.TrimPrefix(receiver.s, "\\"))]
				if class == nil {
					return value{}, &runtimeError{node.line, "Error", fmt.Sprintf("Class %q not found", receiver.s)}
				}
				method := e.lookupMethod(class, member.s)
				if method == nil {
					return value{}, &runtimeError{node.line, "Error", "Call to undefined method " + class.name + "::" + member.s + "()"}
				}
				if err := e.methodAccess(class, method, node.line); err != nil {
					return value{}, err
				}
				if !method.static {
					return value{}, &runtimeError{node.line, "Error", "Non-static method " + method.owner.name + "::" + method.name + "() cannot be called statically"}
				}
				args, err := e.evalArguments(node.args)
				if err != nil {
					return value{}, err
				}
				return e.invokeWithClass(value{}, class, method, args, node.line)
			}
		}
		return value{}, &runtimeError{node.line, "Error", "Value of type " + callee.typeName() + " is not callable"}
	case instanceofExpr:
		v, err := e.eval(node.operand)
		if err != nil {
			return value{}, err
		}
		if v.kind != objectKind {
			return boolean(false), nil
		}
		className := node.className
		if node.classExpr != nil {
			right, err := e.eval(node.classExpr)
			if err != nil {
				return value{}, err
			}
			if right.kind == objectKind {
				className = right.o.className()
			} else if right.kind == stringKind {
				className = strings.TrimPrefix(right.s, "\\")
			} else {
				return boolean(false), nil
			}
		}
		if equalName(className, "self") || equalName(className, "parent") || equalName(className, "static") {
			class, err := e.resolveClassName(className, node.line)
			if err != nil {
				return value{}, err
			}
			className = class.name
		}
		for c := v.o.class; c != nil; c = c.parent {
			if equalName(c.name, className) {
				return boolean(true), nil
			}
			if classImplements(c, className) {
				return boolean(true), nil
			}
		}
		return boolean(false), nil
	case closureExpr:
		return e.evalClosure(node), nil
	case includeExpr:
		return e.evalInclude(node)
	case constantExpr:
		return e.constant(node)
	case classNameExpr:
		if equalName(node.name, "self") || equalName(node.name, "parent") || equalName(node.name, "static") {
			class, err := e.resolveClassName(node.name, node.line)
			if err != nil {
				return value{}, err
			}
			return value{kind: stringKind, s: class.name}, nil
		}
		return value{kind: stringKind, s: node.name}, nil
	case classConstantExpr:
		return e.readClassConstant(node)
	case staticPropertyExpr:
		loc, err := e.staticPropertyLocation(node)
		if err != nil {
			return value{}, err
		}
		return loc.get(), nil
	case staticCallExpr:
		return e.callStaticMethod(node)
	case interpolatedString:
		return e.evalInterpolated(node)
	case ternaryExpr:
		return e.evalTernary(node)
	case coalesceAssignment:
		return e.evalCoalesceAssignment(node)
	case matchExpr:
		return e.evalMatch(node)
	case languageExpr:
		return e.evalLanguage(node)
	case castExpr:
		v, err := e.eval(node.operand)
		if err != nil {
			return value{}, err
		}
		return e.cast(v, node.name, expressionLine(node.operand))
	case literal:
		return node.value, nil
	case arrayLiteral:
		return e.evalArray(node)
	case offsetExpr:
		return e.readOffset(node)
	case propertyExpr:
		return e.readProperty(node)
	case newObject:
		return e.newInstance(node)
	case methodCall:
		return e.callMethod(node)
	case assignment:
		return e.evalAssignment(node)
	case compoundAssignment:
		return e.evalCompound(node)
	case increment:
		return e.evalIncrement(node)
	case functionCall:
		return e.callFunction(node)
	case variable:
		if v, ok := e.readVariable(node.token.text); ok {
			return v, nil
		}
		if node.token.text == "this" {
			return value{}, &runtimeError{node.token.line, "Error", "Using $this when not in object context"}
		}
		e.warn(node.token.line, "Undefined variable $"+node.token.text)
		return value{kind: nullKind}, nil
	case unary:
		v, err := e.eval(node.operand)
		if err != nil {
			return value{}, err
		}
		if node.op.kind == "!" {
			return boolean(!e.toBool(v, expressionLine(node.operand))), nil
		}
		if node.op.kind == "~" {
			return e.bitwiseNot(v, node.op.line)
		}
		n, warn, ok := v.numeric()
		if !ok {
			return value{}, &runtimeError{node.op.line, "TypeError", "Unsupported operand types: " + v.typeName()}
		}
		if warn {
			e.warn(node.op.line, "A non-numeric value encountered")
		}
		if node.op.kind == "+" {
			return n, nil
		}
		if n.kind == floatKind {
			return floating(-n.f), nil
		}
		return arithmetic("-", integer(0), n), nil
	case binary:
		if node.op.kind == "??" {
			v, err := e.quietRead(node.left)
			if err != nil || v.kind != nullKind {
				return v, err
			}
			return e.eval(node.right)
		}
		switch node.op.kind {
		case "==", "!=", "<>", "<=>", "<", "<=", ">", ">=", "===", "!==":
			return e.evalComparison(node)
		}
		if node.op.kind == "&&" || node.op.kind == "||" {
			left, err := e.eval(node.left)
			if err != nil {
				return value{}, err
			}
			truth := e.toBool(left, expressionLine(node.left))
			if node.op.kind == "&&" && !truth || node.op.kind == "||" && truth {
				return boolean(truth), nil
			}
			right, err := e.eval(node.right)
			if err != nil {
				return value{}, err
			}
			return boolean(e.toBool(right, expressionLine(node.right))), nil
		}
		values, err := e.evalOperands(node.left, node.right, expressionLine(node), false)
		if err != nil {
			return value{}, err
		}
		return e.binaryOperation(node.op.kind, values[0], values[1], node.op.line)
	}
	panic("invalid AST expression")
}

func (e *evaluator) evalOperands(left, right expr, line int, reverse bool) ([2]value, error) {
	operands := [2]expr{left, right}
	var values [2]value
	// Compile/evaluate compound operands first; Zend reads bare variables
	// at the opcode, after both subexpressions. Logical ops do not use this.
	for i, operand := range operands {
		if operand == nil {
			continue
		}
		if v, deferred := operand.(variable); deferred && v.token.text != "this" {
			continue
		}
		v, err := e.eval(operand)
		if err != nil {
			return values, err
		}
		values[i] = v
	}
	order := [2]int{0, 1}
	if reverse {
		order = [2]int{1, 0}
	}
	for _, i := range order {
		if v, deferred := operands[i].(variable); deferred && v.token.text != "this" {
			var ok bool
			values[i], ok = e.readVariable(v.token.text)
			if !ok {
				e.warn(line, "Undefined variable $"+v.token.text)
			}
		}
	}
	return values, nil
}

func (e *evaluator) evalComparison(node binary) (value, error) {
	values, err := e.evalOperands(node.left, node.right, expressionLine(node), node.op.kind == ">" || node.op.kind == ">=")
	if err != nil {
		return value{}, err
	}

	if node.op.kind == "===" || node.op.kind == "!==" {
		equal := values[0].identical(values[1])
		if node.op.kind == "!==" {
			equal = !equal
		}
		return boolean(equal), nil
	}
	return e.comparisonValue(node.op.kind, values[0], values[1], expressionLine(node))
}

func (e *evaluator) binaryOperation(op string, left, right value, line int) (value, error) {
	switch op {
	case "&", "|", "^", "<<", ">>":
		return e.bitwise(op, left, right, line)
	}
	if op == "%" {
		return e.modulo(left, right, line)
	}
	if op == "." {
		a, err := e.toText(left, line)
		if err != nil {
			return value{}, err
		}
		b, err := e.toText(right, line)
		return value{kind: stringKind, s: a + b}, err
	}
	if op == "+" && left.kind == arrayKind && right.kind == arrayKind {
		result := left
		for _, entry := range right.a.entries {
			if result.a.find(entry.key) < 0 {
				result.a = result.a.with(entry.key, entry.value)
			}
		}
		return result, nil
	}
	a, warnA, okA := left.numeric()
	if !okA {
		return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Unsupported operand types: %s %s %s", left.typeName(), op, right.typeName())}
	}
	if warnA {
		e.warn(line, "A non-numeric value encountered")
	}
	b, warnB, okB := right.numeric()
	if !okB {
		return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Unsupported operand types: %s %s %s", left.typeName(), op, right.typeName())}
	}
	if warnB {
		e.warn(line, "A non-numeric value encountered")
	}
	if op == "/" && b.asFloat() == 0 {
		return value{}, &runtimeError{line, "DivisionByZeroError", "Division by zero"}
	}
	return arithmetic(op, a, b), nil
}
