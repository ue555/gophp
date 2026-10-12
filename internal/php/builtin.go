package php

import (
	"fmt"
	"strings"
)

func builtinName(name string) bool {
	switch strings.ToLower(name) {
	case "strlen", "count", "gettype", "is_int", "is_float", "is_string", "is_bool", "is_null", "is_array", "is_object", "is_numeric":
		return true
	}
	return false
}

func builtinRedeclaration(name string, version Version) string {
	prefix := "Cannot redeclare "
	if version != PHP83 {
		prefix += "function "
	}
	return prefix + name + "()"
}

func (e *evaluator) callBuiltin(node functionCall) (result value, err error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	// Zend compiles single-argument calls to these functions to dedicated
	// opcodes. Arity failures and count($a, $mode) retain an internal call frame.
	framed := len(args) != 1
	if !framed {
		node.line = expressionLine(node.args[0])
	}
	defer func() {
		if err != nil && framed && e.failureFrames == nil {
			e.captureFailure(callFrame{method: name, line: node.line, args: args})
		}
	}()
	max := 1
	if name == "count" {
		max = 2
	}
	if len(args) < 1 || len(args) > max {
		limit, count := "exactly", 1
		if name == "count" {
			limit = "at least"
			if len(args) > 2 {
				limit, count = "at most", 2
			}
		}
		plural := ""
		if count != 1 {
			plural = "s"
		}
		return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("%s() expects %s %d argument%s, %d given", name, limit, count, plural, len(args))}
	}
	v := args[0]
	switch name {
	case "strlen":
		if v.kind == arrayKind || v.kind == objectKind {
			return value{}, builtinTypeError(node.line, name, 1, "string", "string", v)
		}
		if v.kind == nullKind {
			e.diagnostic("Deprecated", node.line, "strlen(): Passing null to parameter #1 ($string) of type string is deprecated")
		}
		s, err := e.toText(v, node.line)
		return integer(int64(len(s))), err
	case "count":
		mode := int64(0)
		if len(args) == 2 {
			mode, err = e.countMode(args[1], node.line)
			if err != nil {
				return value{}, err
			}
			if mode != 0 && mode != 1 {
				return value{}, &runtimeError{node.line, "ValueError", "count(): Argument #2 ($mode) must be either COUNT_NORMAL or COUNT_RECURSIVE"}
			}
		}
		if v.kind != arrayKind {
			return value{}, builtinTypeError(node.line, name, 1, "value", "Countable|array", v)
		}
		return integer(arrayCount(v.a, mode == 1)), nil
	case "gettype":
		return value{kind: stringKind, s: [...]string{"NULL", "integer", "double", "string", "boolean", "array", "object"}[v.kind]}, nil
	case "is_numeric":
		if v.kind == intKind || v.kind == floatKind {
			return boolean(true), nil
		}
		if v.kind != stringKind {
			return boolean(false), nil
		}
		_, trailing, ok := v.numeric()
		return boolean(ok && !trailing), nil
	default:
		kind := map[string]valueKind{"is_int": intKind, "is_float": floatKind, "is_string": stringKind, "is_bool": boolKind, "is_null": nullKind, "is_array": arrayKind, "is_object": objectKind}[name]
		return boolean(v.kind == kind), nil
	}
}

func builtinTypeError(line int, name string, number int, parameter, expected string, v value) *runtimeError {
	return &runtimeError{line, "TypeError", fmt.Sprintf("%s(): Argument #%d ($%s) must be of type %s, %s given", name, number, parameter, expected, v.diagnosticName())}
}

func (e *evaluator) countMode(v value, line int) (int64, error) {
	n, trailing, ok := v.numeric()
	if !ok || trailing || n.kind == floatKind && floatOutsideInt(n.f) {
		return 0, builtinTypeError(line, "count", 2, "mode", "int", v)
	}
	if v.kind == nullKind {
		e.diagnostic("Deprecated", line, "count(): Passing null to parameter #2 ($mode) of type int is deprecated")
	}
	if n.kind == intKind {
		return n.i, nil
	}
	i := floatToInt(n.f)
	if float64(i) != n.f {
		source := "float " + floatText(n.f, -1)
		if v.kind == stringKind {
			source = "float-string \"" + v.s + "\""
		}
		e.diagnostic("Deprecated", line, "Implicit conversion from "+source+" to int loses precision")
	}
	return i, nil
}

func arrayCount(a *phpArray, recursive bool) int64 {
	n := int64(len(a.entries))
	if recursive {
		for _, entry := range a.entries {
			if entry.value.kind == arrayKind {
				n += arrayCount(entry.value.a, true)
			}
		}
	}
	return n
}

func (e *evaluator) captureFailure(extra ...callFrame) {
	if e.failureFrames != nil {
		return
	}
	for _, frame := range e.frames {
		frame.args = append([]value(nil), frame.args...)
		for i, parameter := range frame.parameters {
			if i >= len(frame.args) {
				break
			}
			if v, bound := frame.locals[parameter.name]; bound {
				frame.args[i] = v
			}
		}
		frame.locals = nil
		e.failureFrames = append(e.failureFrames, frame)
	}
	e.failureFrames = append(e.failureFrames, extra...)
}
