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
	case variable, offsetExpr, propertyExpr, staticPropertyExpr:
		return true
	case nullsafeExpr:
		return !x.(nullsafeExpr).method
	}
	return false
}

func validateUnset(x expr) *runtimeError {
	if hasNullsafe(x) {
		return &runtimeError{expressionLine(x), "", "Can't use nullsafe operator in write context"}
	}
	switch n := x.(type) {
	case staticPropertyExpr:
		// PHP rejects an executed unset, not an unreachable branch containing it.
		return nil
	case variable:
		if n.token.text == "this" {
			return &runtimeError{n.token.line, "", "Cannot unset $this"}
		}
	case offsetExpr:
		if n.key == nil {
			return &runtimeError{n.line, "", "Cannot use [] for unsetting"}
		}
		if _, staticBase := n.base.(staticPropertyExpr); !staticBase {
			if err := validateUnset(n.base); err != nil {
				return err
			}
		}
		return validateExpression(n.key, false)
	case propertyExpr:
		// $this may be the receiver, even though it cannot itself be unset.
		if err := validateExpression(n.base, false); err != nil {
			return err
		}
		return validateExpression(n.nameExpr, false)
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
	if hasNullsafe(x) {
		v, _, err := e.quietNullsafeChain(x, final, false)
		return v, err
	}
	switch n := x.(type) {
	case staticPropertyExpr:
		class, err := e.resolveStaticClass(n.className, n.classExpr, n.line)
		if err != nil {
			return value{}, err
		}
		name, err := e.staticMemberName(n.name, n.nameExpr, n.line)
		if err != nil {
			return value{}, err
		}
		if class.staticProperties[name] == nil {
			return value{}, nil
		}
		n.className, n.classExpr, n.name, n.nameExpr = class.name, nil, name, nil
		loc, err := e.staticPropertyLocation(n)
		if err != nil {
			return value{}, err
		}
		if !class.staticProperties[name].initialized {
			return value{}, nil
		}
		return loc.get(), nil
	case variable:
		if n.token.text == "this" && !final {
			return e.eval(n)
		}
		v, _ := e.readVariable(n.token.text)
		return v, nil
	case propertyExpr:
		base, err := e.quietAccess(n.base, false)
		if err != nil {
			return value{}, err
		}
		name, err := e.instancePropertyName(n.name, n.nameExpr, n.line)
		if err != nil {
			return value{}, err
		}
		n.name = name
		if base.kind != objectKind {
			return value{}, nil
		}
		k, err := e.propertyKey(base.o, n.name, n.line)
		if err != nil {
			return value{}, nil
		}
		v, _ := base.o.properties.get(k)
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
			if classImplements(base.o.class, "ArrayAccess") {
				existsMethod := e.lookupMethod(base.o.class, "offsetExists")
				if existsMethod == nil {
					return value{}, &runtimeError{n.line, "Error", "ArrayAccess::offsetExists() is missing"}
				}
				exists, callErr := e.invoke(base, existsMethod, []callArgument{{value: key, line: n.line}}, n.line)
				if callErr != nil {
					return value{}, callErr
				}
				if !exists.truthy() {
					return value{kind: nullKind}, nil
				}
				getMethod := e.lookupMethod(base.o.class, "offsetGet")
				if getMethod == nil {
					return value{}, &runtimeError{n.line, "Error", "ArrayAccess::offsetGet() is missing"}
				}
				return e.invoke(base, getMethod, []callArgument{{value: key, line: n.line}}, n.line)
			}
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
	if n, ok := x.(staticPropertyExpr); ok {
		return &fatalDiagnostic{n.line, "Attempt to unset static property"}
	}
	if v, ok := x.(variable); ok {
		e.deleteVariable(v.token.text)
		return nil
	}
	if hasDynamicPropertyTarget(x) {
		// unsetLocation already resolves each component once and performs its
		// quiet receiver fetch. Preparing it as an assignment would warn for
		// an undefined receiver before unset gets to suppress that warning.
		_, err := e.unsetLocation(x, true)
		return err
	}
	prepared, err := e.prepareTarget(x)
	if err != nil {
		return err
	}
	_, err = e.unsetLocation(prepared, true)
	return err
}

func hasDynamicPropertyTarget(x expr) bool {
	switch n := x.(type) {
	case propertyExpr:
		return n.nameExpr != nil || hasDynamicPropertyTarget(n.base)
	case offsetExpr:
		return hasDynamicPropertyTarget(n.base)
	}
	return false
}

// Resolve an unset path without creating missing containers. Intermediate
// locations allow nested array deletion to copy only the affected array path.
func (e *evaluator) unsetLocation(x expr, remove bool) (location, error) {
	noop := location{func() value { return value{} }, func(value) error { return nil }, false}
	switch n := x.(type) {
	case staticPropertyExpr:
		return e.staticPropertyLocation(n)
	case variable:
		if _, exists := e.readVariable(n.token.text); !exists {
			e.warn(n.token.line, "Undefined variable $"+n.token.text)
		}
		return e.variableSlotAt(n.token.text, n.token.line), nil
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
		name, err := e.instancePropertyName(n.name, n.nameExpr, n.line)
		if err != nil {
			return noop, err
		}
		n.name = name
		if base.kind != objectKind {
			return noop, nil
		}
		if base.o.enumCase != "" {
			return noop, &runtimeError{n.line, "Error", "Cannot unset readonly property " + base.o.className() + "::$" + n.name}
		}
		k, err := e.propertyKey(base.o, n.name, n.line)
		if err != nil {
			return noop, err
		}
		property := propertyForKey(base.o, k)
		if remove && property != nil && (property.readonly || property.owner.readonly) {
			if _, initialized := base.o.properties.get(k); initialized {
				return noop, &runtimeError{n.line, "Error", "Cannot unset readonly property " + property.owner.name + "::$" + property.name}
			}
			if !e.canInitializeReadonlyProperty(property) {
				return noop, e.readonlyScopeError(property, n.line, "unset")
			}
		}
		if remove {
			if old, ok := base.o.properties.getRaw(k); ok && old.kind == refKind && property != nil && property.typeName != "" {
				old.ref.removeSource(propertyReferenceSource{property: *property, object: base.o, key: k})
			}
			base.o.properties = base.o.properties.without(k)
		}
		return location{func() value { v, _ := base.o.properties.get(k); return v }, func(v value) error { return e.storeObjectProperty(base.o, k, property, v, n.line) }, false}, nil
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
		if base.kind == objectKind && base.o != nil && classImplements(base.o.class, "ArrayAccess") {
			if remove {
				method := e.lookupMethod(base.o.class, "offsetUnset")
				_, err := e.invoke(base, method, []callArgument{{value: key, line: n.line}}, n.line)
				return noop, err
			}
			return location{get: func() value {
				method := e.lookupMethod(base.o.class, "offsetGet")
				v, callErr := e.invoke(base, method, []callArgument{{value: key, line: n.line}}, n.line)
				if callErr != nil {
					e.pendingDiagnostic = callErr
				}
				return v
			}, set: func(v value) error {
				method := e.lookupMethod(base.o.class, "offsetSet")
				_, callErr := e.invoke(base, method, []callArgument{{value: key, line: n.line}, {value: v, line: n.line}}, n.line)
				return callErr
			}}, nil
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
			base.a = e.arrayCopyForWrite(base.a).without(k)
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
			latest.a, err = e.arrayWithAssigned(latest.a, k, v, n.line)
			if err != nil {
				return err
			}
			return parent.set(latest)
		}, false}, nil
	}
	return noop, fmt.Errorf("unsupported unset target")
}
