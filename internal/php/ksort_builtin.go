package php

import (
	"sort"
	"strings"
)

func (e *evaluator) callKsort(node functionCall) (value, error) { return e.callBuiltinSort(node, true) }

func (e *evaluator) callBuiltinSort(node functionCall, byKey bool) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArgumentsCapturingArray(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "flags"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects 1 to 2 arguments"}
	}
	arg := argumentForParameter(args, "array", 0)
	if arg == nil || arg.place == nil && arg.ref == nil {
		return value{}, &runtimeError{node.line, "Error", name + "(): Argument #1 ($array) could not be passed by reference"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, name, 1, "array", "array", bound[0])
	}
	flags := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, name, 2, "flags", "int", bound[1])
		}
		flags = v.i
	}
	out := arrayValue()
	out.a.entries = append([]arrayEntry(nil), bound[0].a.entries...)
	out.a.next, out.a.hasInt, out.a.full = bound[0].a.next, bound[0].a.hasInt, bound[0].a.full
	keyValue := func(k arrayKey) value {
		if k.isString {
			return value{kind: stringKind, s: k.text}
		}
		return integer(k.index)
	}
	var comparisonError error
	sort.SliceStable(out.a.entries, func(i, j int) bool {
		if comparisonError != nil {
			return false
		}
		a, b := keyValue(out.a.entries[i].key), keyValue(out.a.entries[j].key)
		if !byKey {
			a, b = dereference(out.a.entries[i].value), dereference(out.a.entries[j].value)
		}
		order := 0
		switch flags &^ 8 {
		case 1:
			na, _, ok := a.numeric()
			if !ok {
				na = integer(0)
			}
			nb, _, ok := b.numeric()
			if !ok {
				nb = integer(0)
			}
			order = numberOrder(na, nb)
		case 2, 5, 6:
			as, bs := a.text(), b.text()
			if flags&8 != 0 {
				as = asciiKeyLower(as)
				bs = asciiKeyLower(bs)
			}
			if flags&^8 == 6 {
				if naturalKeyLess(as, bs) {
					order = -1
				} else if naturalKeyLess(bs, as) {
					order = 1
				}
			} else {
				order = strings.Compare(as, bs)
			}
		default:
			order, comparisonError = e.valueOrder(a, b, node.line, make(map[*phpArray]bool))
		}
		if name == "krsort" || name == "rsort" || name == "arsort" {
			return order > 0
		}
		return order < 0
	})
	if comparisonError != nil {
		return value{}, comparisonError
	}
	if name == "sort" || name == "rsort" {
		for i := range out.a.entries {
			out.a.entries[i].key = arrayKey{index: int64(i)}
		}
		out.a.next = int64(len(out.a.entries))
		out.a.hasInt = len(out.a.entries) > 0
		out.a.full = false
	}

	if arg.place != nil {
		err = arg.place.assign(out)
	} else {
		err = e.writeReference(arg.ref, out, node.line)
	}
	if err != nil {
		return value{}, err
	}
	return boolean(true), nil
}

func asciiKeyLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

func naturalKeyLess(a, b string) bool {
	for len(a) > 0 && len(b) > 0 {
		a = strings.TrimLeft(a, " \t\n\r\v\f")
		b = strings.TrimLeft(b, " \t\n\r\v\f")
		if len(a) == 0 || len(b) == 0 {
			break
		}
		ac, bc := a[0], b[0]
		if ac >= '0' && ac <= '9' && bc >= '0' && bc <= '9' {
			i, j := 0, 0
			for i < len(a) && a[i] >= '0' && a[i] <= '9' {
				i++
			}
			for j < len(b) && b[j] >= '0' && b[j] <= '9' {
				j++
			}
			an, bn := a[:i], b[:j]
			if a[0] == '0' || b[0] == '0' {
				if an != bn {
					return an < bn
				}
				a, b = a[i:], b[j:]
				continue
			}
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
