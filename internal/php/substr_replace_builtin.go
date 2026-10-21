package php

func (e *evaluator) callSubstrReplace(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	names := []string{"string", "replace", "offset", "length"}
	b, p, _, err := bindCallArguments(args, []parameter{{name: names[0]}, {name: names[1]}, {name: names[2]}, {name: names[3]}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 3 || len(args) > 4 || !p[0] || !p[1] || !p[2] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "substr_replace() expects 3 to 4 arguments"}
	}
	for i := range b {
		if b[i].kind == arrayKind || i == 3 && b[i].kind == nullKind {
			continue
		}
		typ := "string"
		union := "array|string"
		if i >= 2 {
			typ = "int"
			union = "array|int"
			if i == 3 {
				union += "|null"
			}
		}
		v, _, ok := e.coerceType(b[i], typ, false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "substr_replace", i+1, names[i], union, b[i])
		}
		b[i] = v
	}
	if b[0].kind != arrayKind {
		for i := 2; i < 4; i++ {
			if b[i].kind == arrayKind {
				return value{}, &runtimeError{node.line, "TypeError", "substr_replace(): Argument #" + string(rune('1'+i)) + " ($" + names[i] + ") cannot be an array when working on a single string"}
			}
		}
	}
	pick := func(v value, i int, def value) value {
		if v.kind != arrayKind {
			return v
		}
		if i >= len(v.a.entries) {
			return def
		}
		return dereference(v.a.entries[i].value)
	}
	replace := func(s value, i int) (value, error) {
		str, err := e.toText(s, node.line)
		if err != nil {
			return value{}, err
		}
		rep, err := e.toText(pick(b[1], i, value{kind: stringKind}), node.line)
		if err != nil {
			return value{}, err
		}
		off := pick(b[2], i, integer(0))
		length := pick(b[3], i, value{})
		offInt, err := e.cast(off, "int", node.line)
		if err != nil {
			return value{}, err
		}
		start := offInt.i
		n := int64(len(str))
		if start < 0 {
			start += n
		}
		if start < 0 {
			start = 0
		}
		if start > n {
			start = n
		}
		end := n
		if length.kind != nullKind {
			lenInt, err := e.cast(length, "int", node.line)
			if err != nil {
				return value{}, err
			}
			l := lenInt.i
			if l < 0 {
				end = n + l
			} else if l <= n-start {
				end = start + l
			}
			if end < start {
				end = start
			}
		}
		return value{kind: stringKind, s: str[:start] + rep + str[end:]}, nil
	}
	if b[0].kind != arrayKind {
		return replace(b[0], 0)
	}
	out := arrayValue()
	for i, entry := range b[0].a.entries {
		v, err := replace(dereference(entry.value), i)
		if err != nil {
			return value{}, err
		}
		out.a = out.a.with(entry.key, v)
	}
	return out, nil
}
