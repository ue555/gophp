package php

import (
	"strconv"
	"strings"
)

func compareNumericVersions(left, right string) int {
	a := strings.FieldsFunc(left, func(r rune) bool { return r == '.' || r == '-' || r == '_' })
	b := strings.FieldsFunc(right, func(r rune) bool { return r == '.' || r == '-' || r == '_' })
	max := len(a)
	if len(b) > max {
		max = len(b)
	}
	for i := 0; i < max; i++ {
		x, y := "0", "0"
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		xi, xe := strconv.ParseInt(x, 10, 64)
		yi, ye := strconv.ParseInt(y, 10, 64)
		if xe == nil && ye == nil {
			if xi < yi {
				return -1
			}
			if xi > yi {
				return 1
			}
			continue
		}
		if x < y {
			return -1
		}
		if x > y {
			return 1
		}
	}
	return 0
}

func (e *evaluator) callVersionCompare(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "version1"}, {name: "version2"}, {name: "operator"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "version_compare() expects 2 to 3 arguments"}
	}
	a, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "version_compare", 1, "version1", "string", bound[0])
	}
	b, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "version_compare", 2, "version2", "string", bound[1])
	}
	comparison := compareNumericVersions(a.s, b.s)
	if !present[2] {
		return integer(int64(comparison)), nil
	}
	op, _, ok := e.coerceType(bound[2], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "version_compare", 3, "operator", "?string", bound[2])
	}
	switch strings.ToLower(op.s) {
	case "<", "lt":
		return boolean(comparison < 0), nil
	case "<=", "le":
		return boolean(comparison <= 0), nil
	case ">", "gt":
		return boolean(comparison > 0), nil
	case ">=", "ge":
		return boolean(comparison >= 0), nil
	case "==", "=", "eq":
		return boolean(comparison == 0), nil
	case "!=", "<>", "ne":
		return boolean(comparison != 0), nil
	}
	return value{}, &runtimeError{node.line, "ValueError", "version_compare(): Argument #3 ($operator) must be a valid comparison operator"}
}
