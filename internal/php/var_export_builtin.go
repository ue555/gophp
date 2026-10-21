package php

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (e *evaluator) callVarExport(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}, {name: "return"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "var_export() expects 1 or 2 arguments"}
	}
	ret := false
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "bool", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "var_export", 2, "return", "bool", bound[1])
		}
		ret = v.truthy()
	}
	text, err := e.exportValue(bound[0], 0, map[*phpArray]bool{}, node.line)
	if err != nil {
		return value{}, err
	}
	if ret {
		return value{kind: stringKind, s: text}, nil
	}
	_, err = fmt.Fprint(e.stdout, text)
	return value{}, err
}

func exportString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\x00", "' . \"\\0\" . '")
	return "'" + s + "'"
}

// Package/service manifests contain only scalars and arrays. Reject other
// object exports explicitly until __set_state and object recursion are supported.
func (e *evaluator) exportValue(v value, depth int, seen map[*phpArray]bool, line int) (string, error) {
	v = dereference(v)
	switch v.kind {
	case nullKind, resourceKind:
		return "NULL", nil
	case boolKind:
		if v.truthy() {
			return "true", nil
		}
		return "false", nil
	case intKind:
		return strconv.FormatInt(v.i, 10), nil
	case floatKind:
		if math.IsNaN(v.f) {
			return "NAN", nil
		}
		if math.IsInf(v.f, 1) {
			return "INF", nil
		}
		if math.IsInf(v.f, -1) {
			return "-INF", nil
		}
		s := floatText(v.f, -1)
		if math.Abs(v.f) >= 1e-4 && math.Abs(v.f) < 1e17 {
			s = strconv.FormatFloat(v.f, 'f', -1, 64)
		}
		if !strings.ContainsAny(s, ".E") {
			s += ".0"
		}
		return s, nil
	case stringKind:
		return exportString(v.s), nil
	case arrayKind:
		if seen[v.a] {
			if err := e.warn(line, "var_export does not handle circular references"); err != nil {
				return "", err
			}
			return "NULL", nil
		}
		seen[v.a] = true
		defer delete(seen, v.a)
		var b strings.Builder
		b.WriteString("array (\n")
		for _, entry := range v.a.entries {
			b.WriteString(strings.Repeat(" ", depth+2))
			if entry.key.isString {
				b.WriteString(exportString(entry.key.text))
			} else {
				b.WriteString(strconv.FormatInt(entry.key.index, 10))
			}
			b.WriteString(" => ")
			val := dereference(entry.value)
			if val.kind == arrayKind {
				b.WriteString("\n" + strings.Repeat(" ", depth+2))
			}
			s, err := e.exportValue(val, depth+2, seen, line)
			if err != nil {
				return "", err
			}
			b.WriteString(s)
			b.WriteString(",\n")
		}
		b.WriteString(strings.Repeat(" ", depth) + ")")
		return b.String(), nil
	default:
		return "", &runtimeError{line, "Error", "gophp: var_export of objects is not implemented"}
	}
}
