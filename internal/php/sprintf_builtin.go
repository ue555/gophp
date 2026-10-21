package php

import (
	"fmt"
	"strconv"
	"strings"
)

func (e *evaluator) callSprintf(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) == 0 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects at least 1 argument"}
	}
	format, _, ok := e.coerceType(args[0].value, "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, node.name, 1, "format", "string", args[0].value)
	}
	values := args[1:]
	if strings.EqualFold(node.name, "vsprintf") {
		if len(args) != 2 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "vsprintf() expects exactly 2 arguments"}
		}
		if args[1].value.kind != arrayKind {
			return value{}, builtinTypeError(node.line, "vsprintf", 2, "values", "array", args[1].value)
		}
		values = nil
		for _, entry := range args[1].value.a.entries {
			values = append(values, callArgument{value: dereference(entry.value)})
		}
	}
	var out strings.Builder
	index := 0
	for i := 0; i < len(format.s); i++ {
		if format.s[i] != '%' {
			out.WriteByte(format.s[i])
			continue
		}
		if i+1 < len(format.s) && format.s[i+1] == '%' {
			out.WriteByte('%')
			i++
			continue
		}
		i++
		argumentIndex := index
		formatStart := i
		for i < len(format.s) && format.s[i] >= '0' && format.s[i] <= '9' {
			i++
		}
		if i > formatStart && i < len(format.s) && format.s[i] == '$' {
			n, _ := strconv.Atoi(format.s[formatStart:i])
			if n < 1 {
				return value{}, &runtimeError{node.line, "ValueError", "Argument number specifier must be greater than zero"}
			}
			argumentIndex = n - 1
			i++
		} else {
			i = formatStart
			index++
		}
		verbStart := i
		for i < len(format.s) && strings.ContainsRune("-+ 0#123456789.$", rune(format.s[i])) {
			i++
		}
		if i >= len(format.s) {
			return value{}, &runtimeError{node.line, "ValueError", "Missing format specifier at end of string"}
		}
		if argumentIndex >= len(values) {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "Too few arguments"}
		}
		spec := format.s[i]
		item := values[argumentIndex].value
		verb := "%" + format.s[verbStart:i+1]
		var converted any
		switch spec {
		case 's':
			txt, err := e.toText(item, node.line)
			if err != nil {
				return value{}, err
			}
			converted = txt
		case 'd', 'u', 'b', 'o', 'x', 'X', 'c':
			num, _, ok := e.coerceType(item, "int", false, false)
			if !ok {
				return value{}, builtinTypeError(node.line, node.name, index+1, "value", "int", item)
			}
			converted = num.i
		case 'f', 'F', 'e', 'E', 'g', 'G':
			num, _, ok := e.coerceType(item, "float", false, false)
			if !ok {
				return value{}, builtinTypeError(node.line, node.name, index+1, "value", "float", item)
			}
			converted = num.f
		default:
			return value{}, &runtimeError{node.line, "ValueError", "Unknown format specifier"}
		}
		out.WriteString(fmt.Sprintf(verb, converted))
	}
	return value{kind: stringKind, s: out.String()}, nil
}
