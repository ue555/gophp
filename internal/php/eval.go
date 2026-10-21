package php

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type runtimeError struct {
	line    int
	class   string
	message string
}

type exitSignal struct{ code int }

func (s *exitSignal) Error() string { return "script exited" }

func (e *runtimeError) Error() string { return e.message }

type evaluator struct {
	runtimeSymbols
	version               Version
	filename              string
	strictTypes           bool
	stdout, stderr        io.Writer
	variables             variableScope
	rootScope             variableScope
	suspendedScopes       []variableScope
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
	includePath           string
	includeDepth          int
	errorHandler          value
	errorHandlerMask      int64
	errorReporting        int64
	constantAttributes    map[string][]attributeDecl
	exceptionHandler      value
	exceptionHandlerStack []value
	shutdownFunctions     []shutdownCall
	lastError             value
	jsonLastError         int64
	pregLastError         int64
	responseHeaders       []string
	responseStatus        int64
	requestBody           string
	strtokRest            string
	errorHandlerStack     []errorHandlerState
	pendingDiagnostic     error
	suppressedDiagnostics int
	autoloaders           []resolvedCallable
	autoloading           map[string]bool
	activeGenerator       *phpGenerator
	generatorRunning      bool
	generatorEntering     *methodDecl
	generators            *generatorRegistry
	objectIDs             map[*phpObject]int64
	defaultTimezone       string
	defaultMimetype       string
	foreachTempCounter    int
	sqliteConnections     []*sqliteConnection
	processes             []*phpProcess
	openStreams           []*phpStream
	outputBuffers         []outputBuffer
}

// Run parses the entire source before executing it, so a syntax error cannot
// leave partially executed statements. Each invocation has a fresh scope.
func Run(source, filename string, stdout, stderr io.Writer) int {
	return RunWithVersion(source, filename, stdout, stderr, PHP84)
}

// RunWithVersion keeps the same execution pipeline for all supported versions.
func RunWithVersion(source, filename string, stdout, stderr io.Writer, version Version) int {
	return runWithArguments(source, filename, []string{filename}, stdout, stderr, version, 0, nil, nil)
}

func RunWithArguments(source, filename string, arguments []string, stdout, stderr io.Writer, version Version) int {
	return runWithArguments(source, filename, arguments, stdout, stderr, version, 0, nil, nil)
}

func run(source, filename string, stdout, stderr io.Writer, version Version, budget int) int {
	return runWithArguments(source, filename, []string{filename}, stdout, stderr, version, budget, nil, nil)
}

type HTTPRequest struct {
	Method, URI, Host, Body, RemoteAddr string
	Headers                             map[string]string
}
type HTTPResponse struct {
	Status  int
	Headers []string
}

func RunHTTPRequest(source, filename string, request HTTPRequest, stdout, stderr io.Writer, version Version) (int, HTTPResponse) {
	response := HTTPResponse{Status: 200}
	code := runWithArguments(source, filename, []string{filename}, stdout, stderr, version, 0, &request, &response)
	return code, response
}

func runWithArguments(source, filename string, arguments []string, stdout, stderr io.Writer, version Version, budget int, request *HTTPRequest, response *HTTPResponse) (exitCode int) {
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
	e := evaluator{version: version, filename: filename, strictTypes: strictTypes, stdout: stdout, stderr: stderr, variables: make(map[string]value), bounded: budget > 0, steps: budget, loadedFiles: make(map[string]bool), includePath: defaultIncludePath, generators: &generatorRegistry{}, objectIDs: make(map[*phpObject]int64), defaultTimezone: "UTC", defaultMimetype: "text/html", errorReporting: allErrorLevels(version)}
	defer e.closeSQLiteConnections()
	defer e.closeProcesses()
	defer e.closeOpenStreams()
	defer e.flushOutputBuffers()
	if response != nil {
		defer func() {
			if e.responseStatus > 0 {
				response.Status = int(e.responseStatus)
			}
			response.Headers = append([]string(nil), e.responseHeaders...)
		}()
	}
	defer e.generators.closeAll()
	defer func() {
		if err := e.runShutdownFunctions(); err != nil {
			fmt.Fprintf(stderr, "Fatal error: %s in %s\n", err.Error(), e.filename)
			exitCode = 255
		}
	}()
	e.rootScope = e.variables
	argv := arrayValue()
	for i, argument := range arguments {
		argv.a = argv.a.with(arrayKey{index: int64(i)}, value{kind: stringKind, s: argument})
	}
	argc := integer(int64(len(arguments)))
	server := arrayValue()
	for _, entry := range os.Environ() {
		if cut := strings.IndexByte(entry, '='); cut >= 0 {
			server.a = server.a.with(arrayKey{isString: true, text: entry[:cut]}, value{kind: stringKind, s: entry[cut+1:]})
		}
	}
	server.a = server.a.with(arrayKey{isString: true, text: "argv"}, argv)
	server.a = server.a.with(arrayKey{isString: true, text: "argc"}, argc)
	server.a = server.a.with(arrayKey{isString: true, text: "PHP_SELF"}, value{kind: stringKind, s: filename})
	if request != nil {
		e.requestBody = request.Body
		fields := map[string]string{"REQUEST_METHOD": request.Method, "REQUEST_URI": request.URI, "SCRIPT_NAME": "/" + filepath.Base(filename), "SCRIPT_FILENAME": filename, "SERVER_PROTOCOL": "HTTP/1.1", "HTTP_HOST": request.Host, "REMOTE_ADDR": request.RemoteAddr}
		if uri, err := url.ParseRequestURI(request.URI); err == nil {
			fields["QUERY_STRING"] = uri.RawQuery
		}
		for name, content := range request.Headers {
			key := strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
			if key == "CONTENT_TYPE" || key == "CONTENT_LENGTH" {
				fields[key] = content
			}
			fields["HTTP_"+key] = content
		}
		for key, content := range fields {
			server.a = server.a.with(arrayKey{isString: true, text: key}, value{kind: stringKind, s: content})
		}
	}
	e.rootScope["_SERVER"] = server
	e.rootScope["_ENV"] = arrayValue()
	for _, name := range []string{"_GET", "_POST", "_FILES", "_COOKIE", "_REQUEST"} {
		e.rootScope[name] = arrayValue()
	}
	if request != nil {
		e.initializeRequestArrays(*request)
	}
	e.rootScope["argv"] = argv
	e.rootScope["argc"] = argc
	if absolute, pathErr := filepath.Abs(filename); pathErr == nil {
		if canonical, canonicalErr := filepath.EvalSymlinks(absolute); canonicalErr == nil {
			absolute = canonical
		}
		e.loadedFiles[absolute] = true
	}
	e.callableDiagnostics(statements)
	e.controlDiagnostics(statements, nil)
	if err := e.executeProgram(statements); err != nil {
		var stopped *exitSignal
		if errors.As(err, &stopped) {
			return stopped.code
		}
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
		if e.exceptionHandler.kind != nullKind {
			if caught, ok := e.catchable(err); ok {
				if target, valid := e.resolveCallable(e.exceptionHandler, caught.object.o.exception.line); valid {
					_, handlerErr := e.invokeResolved(target, []callArgument{{value: caught.object, line: caught.object.o.exception.line}}, caught.object.o.exception.line)
					if handlerErr == nil {
						return 0
					}
					err = handlerErr
				}
			}
		}
		if errors.As(err, &raised) {
			data := raised.object.o.exception
			fmt.Fprintf(stderr, "Fatal error: Uncaught %s\n  thrown in %s on line %d\n", exceptionDiagnostic(raised.object), data.file, data.line)
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
		// Enum built-in method redeclarations are checked when the enum is
		// bound. A preceding statement must run before this declaration error.
		if hasExecutedBefore[decl] && decl.kind == "enum" {
			for name := range decl.methods {
				if name == "__construct" || name == "cases" || decl.enumBacking != "" && (name == "from" || name == "tryfrom") {
					return false, nil
				}
			}
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
					if e.classes[traitKey] == nil && len(e.autoloaders) != 0 {
						if err := e.tryAutoload(name, decl.line); err != nil {
							return false, err
						}
					}
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
				if e.classes[interfaceKey] == nil && len(e.autoloaders) != 0 {
					if err := e.tryAutoload(name, decl.line); err != nil {
						return false, err
					}
				}
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
				if e.classes[parentKey] == nil && len(e.autoloaders) != 0 {
					if err := e.tryAutoload(decl.parentName, decl.line); err != nil {
						return false, err
					}
				}
				if e.classes[parentKey] == nil {
					if canDefineBefore[decl] {
						return false, nil
					}
					return false, &runtimeError{decl.line, "Error", "Class \"" + decl.parentName + "\" not found"}
				}
			}
		}
		// A method contract may refer to a class declared later in this file.
		// Unlike an extends dependency, that type must not be bound early just
		// to establish a variance relationship. PHP checks the relationship at
		// this class declaration and can report that the type is unavailable.
		if decl.parentName != "" || len(decl.interfaceNames) != 0 {
			e.initClasses()
			references := classTypeReferences(decl)
			if parent := e.classes[strings.ToLower(decl.parentName)]; parent != nil {
				references = append(references, classTypeReferences(parent)...)
			}
			for _, interfaceName := range decl.interfaceNames {
				if iface := e.classes[strings.ToLower(interfaceName)]; iface != nil {
					references = append(references, classTypeReferences(iface)...)
				}
			}
			// A prior runtime statement may have required the missing type.
			// Defer this declaration so that execution can define it first.
			if hasExecutedBefore[decl] {
				for _, typeName := range references {
					key := strings.ToLower(typeName)
					if pending[key] == nil && e.classes[key] == nil {
						return false, nil
					}
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
		return expressionMayDefineClass(n.base) || expressionMayDefineClass(n.nameExpr)
	case methodCall:
		return true
	case nullsafeExpr:
		if expressionMayDefineClass(n.base) || expressionMayDefineClass(n.nameExpr) {
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
	case firstClassCallable:
		return expressionMayDefineClass(n.callee)
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

type gotoSignal struct {
	label string
	line  int
}

func (g *gotoSignal) Error() string { return "goto " + g.label }

func (e *evaluator) execute(statements []statement) error {
	labels := make(map[string]int)
	for i, s := range statements {
		if s.kind == "label" {
			labels[s.names[0]] = i
		}
	}
	for start := 0; start < len(statements); {
		err := e.executeFrom(statements, start)
		jump, ok := err.(*gotoSignal)
		if !ok {
			return err
		}
		index, found := labels[jump.label]
		if !found {
			return err
		}
		start = index + 1
	}
	return nil
}

func (e *evaluator) executeFrom(statements []statement, start int) (err error) {
	defer func() {
		if e.pendingDiagnostic != nil {
			err = e.pendingDiagnostic
			e.pendingDiagnostic = nil
		}
	}()
	for _, statement := range statements[start:] {
		if e.pendingDiagnostic != nil {
			return e.pendingDiagnostic
		}
		switch statement.kind {
		case "label":
			continue
		case "goto":
			return &gotoSignal{statement.names[0], statement.line}
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
			if e.defineConstant(decl.name, v, decl.line) && len(decl.attributes) != 0 {
				if e.constantAttributes == nil {
					e.constantAttributes = make(map[string][]attributeDecl)
				}
				e.constantAttributes[decl.name] = decl.attributes
			}
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
				x := statement.expressions[0]
				if e.currentMethod != nil && e.currentMethod.returnsByReference && !methodContainsYield(e.currentMethod.body) {
					if !referenceSourceSyntax(x) {
						return &runtimeError{statement.line, "Error", "gophp: returning an expression by reference is not implemented"}
					}
					e.foreachTempCounter++
					name := fmt.Sprintf("\x00return_reference_%d", e.foreachTempCounter)
					target := variable{token{kind: "variable", text: name, line: statement.line}}
					_, err = e.evalReferenceAssignment(referenceAssignment{target: target, source: x, line: statement.line})
					if err == nil {
						v = e.variables[name]
						delete(e.variables, name)
					}
				} else {
					v, err = e.eval(x)
				}
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
		if statement.kind == "expression" && len(statement.expressions) == 1 {
			if _, unusedVariable := statement.expressions[0].(variable); unusedVariable {
				// PHP does not fetch a simple variable whose value is discarded.
				continue
			}
		}
		if statement.kind == "if" {
			condition, err := e.eval(statement.condition)
			if err != nil {
				return err
			}
			branch := statement.elseBranch
			truth := e.toBool(condition, statement.conditionLine)
			if e.pendingDiagnostic != nil {
				return e.pendingDiagnostic
			}
			if truth {
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

func (e *evaluator) warn(line int, message string) error {
	return e.diagnostic("Warning", line, message)
}

func (e *evaluator) diagnostic(level string, line int, message string) error {
	if e.pendingDiagnostic != nil {
		return e.pendingDiagnostic
	}
	if err := e.handleDiagnostic(level, line, message); err != nil {
		e.pendingDiagnostic = err
		return err
	}
	return nil
}

func (e *evaluator) toText(v value, line int) (string, error) {
	if e.version == PHP85 && v.kind == floatKind && math.IsNaN(v.f) {
		if err := e.warn(line, "unexpected NAN value was coerced to string"); err != nil {
			return "", err
		}
	}
	if v.kind == arrayKind {
		if err := e.warn(line, "Array to string conversion"); err != nil {
			return "", err
		}
	}
	if v.kind == objectKind {
		if v.o != nil && v.o.class != nil && v.o.class.name == "ReflectionNamedType" {
			name := v.o.reflectedTypeName
			if v.o.reflectedNullable && name != "mixed" && name != "null" && !strings.Contains(name, "|") {
				name = "?" + name
			}
			return name, nil
		}
		if method := e.lookupMethod(v.o.class, "__toString"); method != nil {
			result, err := e.invoke(v, method, nil, line)
			if err != nil {
				return "", err
			}
			if result.kind != stringKind {
				return "", &runtimeError{line, "TypeError", v.o.className() + "::__toString(): Return value must be of type string, " + result.diagnosticName() + " returned"}
			}
			return result.s, nil
		}
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

func (e *evaluator) eval(expression expr) (result value, err error) {
	if e.pendingDiagnostic != nil {
		return value{}, e.pendingDiagnostic
	}
	defer func() {
		if e.pendingDiagnostic != nil {
			result, err = value{}, e.pendingDiagnostic
			e.pendingDiagnostic = nil
		}
	}()
	if hasNullsafe(expression) {
		v, _, err := e.evalChain(expression)
		return v, err
	}
	switch node := expression.(type) {
	case firstClassCallable:
		return e.evalFirstClass(node)
	case throwExpr:
		v, err := e.eval(node.operand)
		if err != nil {
			return value{}, err
		}
		return value{}, e.throwValue(v, node.line)
	case dynamicCall:
		result, err := e.callDynamic(node)
		return dereference(result), err
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
		if node.dynamic != nil {
			selector, err := e.eval(node.dynamic)
			if err != nil {
				return value{}, err
			}
			if selector.kind == objectKind && selector.o != nil {
				return value{kind: stringKind, s: selector.o.className()}, nil
			}
			if selector.kind != stringKind {
				return value{}, &runtimeError{node.line, "TypeError", "Class name must be a valid object or a string"}
			}
			class, err := e.resolveClassName(selector.s, node.line)
			if err != nil {
				return value{}, err
			}
			return value{kind: stringKind, s: class.name}, nil
		}
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
		prepared, err := e.prepareStaticProperty(node)
		if err != nil {
			return value{}, err
		}
		node = prepared
		loc, err := e.staticPropertyLocation(node)
		if err != nil {
			return value{}, err
		}
		class, err := e.resolveClassName(node.className, node.line)
		if err != nil {
			return value{}, err
		}
		if slot := class.staticProperties[node.name]; slot != nil && !slot.initialized {
			return value{}, uninitializedPropertyError(slot.decl, node.line)
		}
		return loc.get(), nil
	case staticCallExpr:
		v, err := e.callStaticMethod(node)
		return dereference(v), err
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
	case yieldExpr:
		return e.evalYield(node)
	case cloneExpr:
		return e.evalClone(node)
	case dynamicVariable:
		return e.evalDynamicVariable(node)
	case methodCall:
		v, err := e.callMethod(node)
		return dereference(v), err
	case assignment:
		return e.evalAssignment(node)
	case referenceAssignment:
		return e.evalReferenceAssignment(node)
	case compoundAssignment:
		return e.evalCompound(node)
	case increment:
		return e.evalIncrement(node)
	case functionCall:
		v, err := e.callFunction(node)
		return dereference(v), err
	case variable:
		if v, ok := e.readVariable(node.token.text); ok {
			return v, nil
		}
		if node.token.text == "this" {
			return value{}, &runtimeError{node.token.line, "Error", "Using $this when not in object context"}
		}
		if err := e.warn(node.token.line, "Undefined variable $"+node.token.text); err != nil {
			return value{}, err
		}
		return value{kind: nullKind}, nil
	case unary:
		if node.op.kind == "@" {
			e.suppressedDiagnostics++
			defer func() { e.suppressedDiagnostics-- }()
			return e.eval(node.operand)
		}
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
			if err := e.warn(node.op.line, "A non-numeric value encountered"); err != nil {
				return value{}, err
			}
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
		if node.op.kind == "&&" || node.op.kind == "||" || node.op.kind == "and" || node.op.kind == "or" || node.op.kind == "xor" {
			left, err := e.eval(node.left)
			if err != nil {
				return value{}, err
			}
			truth := e.toBool(left, expressionLine(node.left))
			if (node.op.kind == "&&" || node.op.kind == "and") && !truth || (node.op.kind == "||" || node.op.kind == "or") && truth {
				return boolean(truth), nil
			}
			right, err := e.eval(node.right)
			if err != nil {
				return value{}, err
			}
			rightTruth := e.toBool(right, expressionLine(node.right))
			if node.op.kind == "xor" {
				return boolean(truth != rightTruth), nil
			}
			return boolean(rightTruth), nil
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
				if err := e.warn(line, "Undefined variable $"+v.token.text); err != nil {
					return values, err
				}
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
		if err := e.warn(line, "A non-numeric value encountered"); err != nil {
			return value{}, err
		}
	}
	b, warnB, okB := right.numeric()
	if !okB {
		return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Unsupported operand types: %s %s %s", left.typeName(), op, right.typeName())}
	}
	if warnB {
		if err := e.warn(line, "A non-numeric value encountered"); err != nil {
			return value{}, err
		}
	}
	if op == "/" && b.asFloat() == 0 {
		return value{}, &runtimeError{line, "DivisionByZeroError", "Division by zero"}
	}
	return arithmetic(op, a, b), nil
}

func (e *evaluator) callDynamic(node dynamicCall) (value, error) {
	callee, err := e.eval(node.callee)
	if err != nil {
		return value{}, err
	}
	if callee.kind == stringKind {
		if strings.Contains(callee.s, "::") {
			target, ok := e.resolveCallable(callee, node.line)
			if !ok {
				return value{}, &runtimeError{node.line, "Error", "Call to invalid callback " + callee.s}
			}
			args, err := e.evalArguments(node.args)
			if err != nil {
				return value{}, err
			}
			return e.invokeResolved(target, args, node.line)
		}
		return e.callFunction(functionCall{name: strings.TrimPrefix(callee.s, "\\"), args: node.args, line: node.line})
	}
	if callee.kind == objectKind && callee.o.closure != nil {
		return e.invokeClosure(callee.o.closure, node.args, node.line)
	}
	if callee.kind == objectKind {
		return e.callMethod(methodCall{receiver: literal{value: callee, line: node.line}, name: "__invoke", args: node.args, line: node.line})
	}
	if callee.kind == arrayKind && len(callee.a.entries) == 2 && !callee.a.entries[0].key.isString && callee.a.entries[0].key.index == 0 && !callee.a.entries[1].key.isString && callee.a.entries[1].key.index == 1 {
		receiver, member := dereference(callee.a.entries[0].value), dereference(callee.a.entries[1].value)
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
}
