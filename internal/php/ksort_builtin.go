package php

import (
	"sort"
	"strconv"
	"strings"
)

func (e *evaluator) callKsort(node functionCall) (value, error) {
	if len(node.args) < 1 || len(node.args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "ksort() expects 1 to 2 arguments"}
	}
	if !writable(node.args[0]) {
		return value{}, &runtimeError{node.line, "Error", "ksort(): Argument #1 ($array) cannot be passed by reference"}
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "flags"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "ksort() expects at least 1 argument"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "ksort", 1, "array", "array", bound[0])
	}
	flags := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "ksort", 2, "flags", "int", bound[1])
		}
		flags = v.i
	}
	out := arrayValue()
	out.a.entries = append([]arrayEntry(nil), bound[0].a.entries...)
	out.a.next, out.a.hasInt, out.a.full = bound[0].a.next, bound[0].a.hasInt, bound[0].a.full
	sort.SliceStable(out.a.entries, func(i, j int) bool {
		a, b := out.a.entries[i].key, out.a.entries[j].key
		if !a.isString && !b.isString {
			return a.index < b.index
		}
		as, bs := a.text, b.text
		if !a.isString {
			as = strconv.FormatInt(a.index, 10)
		}
		if !b.isString {
			bs = strconv.FormatInt(b.index, 10)
		}
		if flags&8 != 0 {
			as, bs = strings.ToLower(as), strings.ToLower(bs)
		}
		if flags&6 == 1 {
			af, ae := strconv.ParseFloat(as, 64)
			bf, be := strconv.ParseFloat(bs, 64)
			if ae == nil && be == nil {
				return af < bf
			}
		}
		if flags&6 == 6 {
			return naturalKeyLess(as, bs)
		}
		return as < bs
	})
	if _, err := e.evalAssignment(assignment{target: node.args[0], source: literal{value: out, line: node.line}, line: node.line}); err != nil {
		return value{}, err
	}
	return boolean(true), nil
}

func naturalKeyLess(a, b string) bool {
	for len(a) > 0 && len(b) > 0 {
		ac, bc := a[0], b[0]
		if ac >= '0' && ac <= '9' && bc >= '0' && bc <= '9' {
			i, j := 0, 0
			for i < len(a) && a[i] >= '0' && a[i] <= '9' {
				i++
			}
			for j < len(b) && b[j] >= '0' && b[j] <= '9' {
				j++
			}
			an, bn := strings.TrimLeft(a[:i], "0"), strings.TrimLeft(b[:j], "0")
			if len(an) != len(bn) {
				return len(an) < len(bn)
			}
			if an != bn {
				return an < bn
			}
			a, b = a[i:], b[j:]
			continue
		}
		if ac != bc {
			return ac < bc
		}
		a, b = a[1:], b[1:]
	}
	return len(a) < len(b)
}
