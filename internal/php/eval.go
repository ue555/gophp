package php

import (
	"fmt"
	"io"
	"math"
)

type runtimeError struct {
	line    int
	class   string
	message string
}

func (e *runtimeError) Error() string { return e.message }

type evaluator struct {
	version               Version
	filename              string
	stdout, stderr        io.Writer
	variables             map[string]value
	classes               map[string]*classDecl
	predeclared           map[*classDecl]bool
	functions             map[string]*methodDecl
	predeclaredFunctions  map[*methodDecl]bool
	bounded               bool
	steps                 int
	scope                 *classDecl
	frames, failureFrames []callFrame
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
	if err == nil {
		statements, err = parseWithVersion(tokens, version)
	}
	if err != nil {
		e := err.(*syntaxError)
		fmt.Fprintf(stderr, "Parse error: %s in %s on line %d\n", e.message, filename, e.line)
		return 255
	}
	failure := validateStatements(statements)
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
	e := evaluator{version: version, filename: filename, stdout: stdout, stderr: stderr, variables: make(map[string]value), bounded: budget > 0, steps: budget}
	e.callableDiagnostics(statements)
	if err := e.executeProgram(statements); err != nil {
		if _, ok := err.(*returned); ok {
			return 0
		}
		if failure, ok := err.(*fatalDiagnostic); ok {
			fmt.Fprintf(stderr, "Fatal error: %s in %s on line %d\n", failure.message, filename, failure.line)
			if version == PHP85 {
				fmt.Fprintf(stderr, "Stack trace:\n%s", e.stackTrace())
			}
			return 255
		}
		if failure, ok := err.(*runtimeError); ok {
			if failure.class == "" {
				fmt.Fprintf(stderr, "Fatal error: %s in %s on line %d\n", failure.message, filename, failure.line)
				return 255
			}
			fmt.Fprintf(stderr, "Fatal error: Uncaught %s: %s in %s:%d\nStack trace:\n%s  thrown in %s on line %d\n", failure.class, failure.message, filename, failure.line, e.stackTrace(), filename, failure.line)
			return 255
		}
		fmt.Fprintf(stderr, "gophp: %v\n", err)
		return 1
	}
	return 0
}

func (e *evaluator) executeProgram(statements []statement) error {
	e.predeclared = make(map[*classDecl]bool)
	e.predeclaredFunctions = make(map[*methodDecl]bool)
	// Unconditional top-level classes are available before their declaration.
	// Classes inside branches are registered only if execution reaches them.
	for _, s := range statements {
		if s.kind == "function" {
			if err := e.declareFunction(s.function); err != nil {
				return err
			}
			e.predeclaredFunctions[s.function] = true
		}
		if s.kind == "class" {
			if err := e.declareClass(s.class); err != nil {
				return err
			}
			e.predeclared[s.class] = true
		}
	}
	return e.execute(statements)
}

func (e *evaluator) execute(statements []statement) error {
	for _, statement := range statements {
		switch statement.kind {
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
			return &returned{v}
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
	switch node := expression.(type) {
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
		if v, ok := e.variables[node.token.text]; ok {
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
			values[i], ok = e.variables[v.token.text]
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
