package php

import "fmt"

func (p *parser) languageArguments(name string) ([]expr, error) {
	if err := p.expect("("); err != nil {
		return nil, err
	}
	if p.current().kind == ")" {
		return nil, unexpected(p.current())
	}
	var args []expr
	for {
		x, err := p.expression(0, ",", ")")
		if err != nil {
			return nil, err
		}
		if name == "unset" && !writable(x) {
			return nil, syntax(expressionLine(x), "invalid unset target")
		}
		args = append(args, x)
		if name == "empty" || !p.accept(",") {
			break
		}
		if p.current().kind == ")" {
			break
		}
	}
	return args, p.expect(")")
}

func issetTarget(x expr) bool {
	switch x.(type) {
	case variable, offsetExpr, propertyExpr:
		return true
	}
	return false
}

func validateUnset(x expr) *runtimeError {
	switch n := x.(type) {
	case variable:
		if n.token.text == "this" {
			return &runtimeError{n.token.line, "", "Cannot unset $this"}
		}
	case offsetExpr:
		if n.key == nil {
			return &runtimeError{n.line, "", "Cannot use [] for unsetting"}
		}
		if err := validateUnset(n.base); err != nil {
			return err
		}
		return validateExpression(n.key, false)
	case propertyExpr:
		// $this may be the receiver, even though it cannot itself be unset.
		return validateExpression(n.base, false)
	}
	return nil
}

func (e *evaluator) evalLanguage(node languageExpr) (value, error) {
	for _, arg := range node.args {
		v, err := e.quietRead(arg)
		if err != nil {
			return value{}, err
		}
		if node.name == "empty" {
			return boolean(!e.toBool(v, expressionLine(arg))), nil
		}
		if v.kind == nullKind {
			return boolean(false), nil
		}
	}
	return boolean(true), nil
}

// Quiet fetches apply only to the access path. Key expressions and calls still
// run normally, with their own diagnostics and side effects.
func (e *evaluator) quietRead(x expr) (value, error) {
	return e.quietAccess(x, true)
}

func (e *evaluator) quietAccess(x expr, final bool) (value, error) {
	switch n := x.(type) {
	case variable:
		if n.token.text == "this" && !final {
			return e.eval(n)
		}
		return e.variables[n.token.text], nil
	case propertyExpr:
		base, err := e.quietAccess(n.base, false)
		if err != nil {
			return value{}, err
		}
		if base.kind != objectKind {
			return value{}, nil
		}
		if e.propertyAccess(base.o, n.name, n.line) != nil {
			return value{}, nil
		}
		v, _ := base.o.properties.get(arrayKey{text: n.name, isString: true})
		return v, nil
	case offsetExpr:
		var base, key value
		var err error
		if _, bare := n.base.(variable); bare {
			key, err = e.eval(n.key)
			if err == nil {
				base, err = e.quietAccess(n.base, false)
			}
		} else {
			base, err = e.quietAccess(n.base, false)
			if err == nil {
				key, err = e.eval(n.key)
			}
		}
		if err != nil {
			return value{}, err
		}
		switch base.kind {
		case arrayKind:
			k, err := e.accessKey(key, n.line, "isset")
			if err != nil {
				return value{}, err
			}
			v, _ := base.a.get(k)
			return v, nil
		case objectKind:
			return value{}, &runtimeError{n.line, "Error", "Cannot use object of type " + base.o.className() + " as array"}
		case stringKind:
			if key.kind == arrayKind || key.kind == objectKind {
				if !final {
					return value{}, &runtimeError{n.line, "TypeError", "Cannot access offset of type " + key.typeName() + " on string"}
				}
				return value{}, nil
			}
			if key.kind == stringKind {
				num, trailing, ok := key.numeric()
				if !ok || (trailing && final) || num.kind != intKind {
					return value{}, nil
				}
				if trailing {
					e.warn(n.line, "Illegal string offset \""+key.s+"\"")
				}
				key = num
			}
			if key.kind == floatKind {
				// String existence checks use strict integer conversion, unlike
				// an explicit (int) cast (fractional keys are deprecated).
				if final {
					k, err := e.arrayKey(key, n.line)
					if err != nil {
						return value{}, err
					}
					key = integer(k.index)
				} else {
					key = integer(floatToInt(key.f))
				}
			}
			i := key.i
			if i < 0 {
				i += int64(len(base.s))
			}
			if i < 0 || i >= int64(len(base.s)) {
				return value{}, nil
			}
			return value{kind: stringKind, s: base.s[i : i+1]}, nil
		}
		return value{}, nil
	default:
		return e.eval(x)
	}
}

func (e *evaluator) accessKey(v value, line int, mode string) (arrayKey, error) {
	if v.kind == arrayKind || v.kind == objectKind {
		message := "Cannot access offset of type " + v.typeName() + " in isset or empty"
		if mode == "unset" {
			message = "Cannot unset offset of type " + v.typeName() + " on array"
		}
		return arrayKey{}, &runtimeError{line, "TypeError", message}
	}
	// UNSET_DIM deliberately has no null-offset deprecation, even in 8.5.
	if mode == "unset" && v.kind == nullKind {
		return arrayKey{isString: true}, nil
	}
	return e.arrayKey(v, line)
}

func (a *phpArray) without(k arrayKey) *phpArray {
	i := a.find(k)
	if i < 0 {
		return a
	}
	b := *a
	b.entries = make([]arrayEntry, 0, len(a.entries)-1)
	b.entries = append(b.entries, a.entries[:i]...)
	b.entries = append(b.entries, a.entries[i+1:]...)
	// Deletion never rewinds the next index, including deletion of the last
	// element. At LONG_MAX the next index equals the removed key itself.
	return &b
}

func (e *evaluator) unset(x expr) error {
	if v, ok := x.(variable); ok {
		delete(e.variables, v.token.text)
		return nil
	}
	prepared, err := e.prepareTarget(x)
	if err != nil {
		return err
	}
	_, err = e.unsetLocation(prepared, true)
	return err
}

// Resolve an unset path without creating missing containers. Intermediate
// locations allow nested array deletion to copy only the affected array path.
func (e *evaluator) unsetLocation(x expr, remove bool) (location, error) {
	noop := location{func() value { return value{} }, func(value) error { return nil }}
	switch n := x.(type) {
	case variable:
		if _, exists := e.variables[n.token.text]; !exists {
			e.warn(n.token.line, "Undefined variable $"+n.token.text)
		}
		return location{func() value { return e.variables[n.token.text] }, func(v value) error { e.variables[n.token.text] = v; return nil }}, nil
	case propertyExpr:
		var base value
		var err error
		if v, bare := n.base.(variable); bare && v.token.text == "this" {
			base, err = e.eval(n.base)
		} else if bare && remove {
			base, err = e.quietRead(n.base)
		} else if issetTarget(n.base) {
			var loc location
			loc, err = e.unsetLocation(n.base, false)
			if err == nil {
				base = loc.get()
			}
		} else {
			base, err = e.eval(n.base)
		}
		if err != nil {
			return noop, err
		}
		if base.kind != objectKind {
			return noop, nil
		}
		if err := e.propertyAccess(base.o, n.name, n.line); err != nil {
			return noop, err
		}
		k := arrayKey{text: n.name, isString: true}
		if remove {
			base.o.properties = base.o.properties.without(k)
		}
		return location{func() value { v, _ := base.o.properties.get(k); return v }, func(v value) error { base.o.properties = base.o.properties.with(k, v); return nil }}, nil
	case offsetExpr:
		parent, err := e.unsetLocation(n.base, false)
		if err != nil {
			return noop, err
		}
		base := parent.get()
		key, err := e.eval(n.key)
		if err != nil {
			return noop, err
		}
		if base.kind != arrayKind {
			switch base.kind {
			case nullKind:
				return noop, nil
			case boolKind:
				if base.i == 0 {
					e.diagnostic("Deprecated", n.line, "Automatic conversion of false to array is deprecated")
					return noop, nil
				}
			case objectKind:
				return noop, &runtimeError{n.line, "Error", "Cannot use object of type " + base.o.className() + " as array"}
			case stringKind:
				return noop, &runtimeError{n.line, "Error", "Cannot unset string offsets"}
			}
			return noop, &runtimeError{n.line, "Error", "Cannot unset offset in a non-array variable"}
		}
		k, err := e.accessKey(key, n.line, "unset")
		if err != nil {
			return noop, err
		}
		if remove {
			base.a = base.a.without(k)
			if err = parent.set(base); err != nil {
				return noop, err
			}
		}
		return location{func() value {
			v := parent.get()
			if v.kind != arrayKind {
				return value{}
			}
			child, _ := v.a.get(k)
			return child
		}, func(v value) error {
			latest := parent.get()
			if latest.kind != arrayKind {
				return nil
			}
			latest.a = latest.a.with(k, v)
			return parent.set(latest)
		}}, nil
	}
	return noop, fmt.Errorf("unsupported unset target")
}
