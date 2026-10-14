package php

import (
	"fmt"
	"strings"
)

func constantBuiltin(name string) bool {
	switch strings.ToLower(name) {
	case "define", "defined", "constant":
		return true
	}
	return false
}

func (e *evaluator) callConstantBuiltin(node functionCall) (result value, err error) {
	name := strings.ToLower(node.name)
	evaluated, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "constant_name"}}
	if name == "define" {
		params = append(params, parameter{name: "value"}, parameter{name: "case_insensitive"})
	}
	bound, present, _, err := bindCallArguments(evaluated, params, node.line)
	if err != nil {
		return value{}, err
	}
	args := argumentValues(evaluated)
	defer func() {
		if err != nil && e.failureFrames == nil {
			e.captureFailure(callFrame{method: name, line: node.line, filename: e.filename, args: args})
		}
	}()
	min, max := 1, 1
	if name == "define" {
		min, max = 2, 3
	}
	if len(args) < min || len(args) > max || !present[0] || name == "define" && !present[1] {
		which, count := "exactly", min
		if min != max {
			if len(args) < min {
				which = "at least"
			} else {
				which, count = "at most", max
			}
		}
		plural := ""
		if count != 1 {
			plural = "s"
		}
		return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("%s() expects %s %d argument%s, %d given", name, which, count, plural, len(args))}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, name, 1, "constant_name", "string", bound[0])
	}
	if bound[0].kind == nullKind {
		e.diagnostic("Deprecated", node.line, name+"(): Passing null to parameter #1 ($constant_name) of type string is deprecated")
	}
	constant, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	if name != "define" {
		constant = strings.TrimPrefix(constant, "\\")
	}
	if strings.Contains(constant, "::") {
		return value{}, &runtimeError{node.line, "Error", "gophp: class constants are not supported"}
	}
	if name == "define" {
		if present[2] {
			if bound[2].kind == arrayKind || bound[2].kind == objectKind {
				return value{}, builtinTypeError(node.line, name, 3, "case_insensitive", "bool", bound[2])
			}
			if bound[2].kind == nullKind {
				e.diagnostic("Deprecated", node.line, "define(): Passing null to parameter #3 ($case_insensitive) of type bool is deprecated")
			}
			if e.toBool(bound[2], node.line) {
				e.warn(node.line, "define(): Argument #3 ($case_insensitive) is ignored since declaration of case-insensitive constants is no longer supported")
			}
		}
		return boolean(e.defineConstant(constant, bound[1], node.line)), nil
	}
	e.initConstants()
	v, ok := e.constants[constant]
	if name == "defined" {
		return boolean(ok), nil
	}
	if !ok {
		return value{}, &runtimeError{node.line, "Error", fmt.Sprintf("Undefined constant %q", constant)}
	}
	return v, nil
}
