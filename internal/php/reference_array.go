package php

func typeIncludesNull(name string) bool {
	t := parseDeclaredType(name, false)
	return t.has("null") || t.has("mixed")
}

// A reference place is a binding, rather than the PHP value stored in it.
// replace changes just this binding; assign follows an existing reference.
type referencePlace struct {
	read    func() value
	assign  func(value) error
	replace func(value) error
	arrays  []*phpArray
	source  *propertyReferenceSource
}

// PHP drops a reference wrapper with no other owner when a copy-on-write
// array is separated. The interpreter's arrays are immutable snapshots, so
// inspect live bindings before copying an entry. Multiple variables holding
// the same array pointer do not count as additional owners of its elements.
func (e *evaluator) referenceSharedOutside(a *phpArray, key arrayKey, cell *referenceCell) bool {
	seenArrays := map[*phpArray]bool{}
	seenObjects := map[*phpObject]bool{}
	seenRefs := map[*referenceCell]bool{}
	var visit func(value) bool
	visit = func(v value) bool {
		switch v.kind {
		case refKind:
			if v.ref == cell {
				return true
			}
			if seenRefs[v.ref] {
				return false
			}
			seenRefs[v.ref] = true
			return visit(v.ref.value)
		case arrayKind:
			if v.a == nil || seenArrays[v.a] {
				return false
			}
			seenArrays[v.a] = true
			for _, entry := range v.a.entries {
				if v.a == a && entry.key == key {
					continue
				}
				if visit(entry.value) {
					return true
				}
			}
		case objectKind:
			if v.o == nil || seenObjects[v.o] {
				return false
			}
			seenObjects[v.o] = true
			if v.o.properties != nil && visit(value{kind: arrayKind, a: v.o.properties}) {
				return true
			}
			if v.o.closure != nil {
				if visit(v.o.closure.boundThis) {
					return true
				}
				for _, captured := range v.o.closure.captures {
					if visit(captured) {
						return true
					}
				}
			}
		}
		return false
	}
	for _, scope := range []variableScope{e.variables, e.rootScope, e.staticScope} {
		for _, v := range scope {
			if visit(v) {
				return true
			}
		}
	}
	for _, scope := range e.suspendedScopes {
		for _, v := range scope {
			if visit(v) {
				return true
			}
		}
	}
	for _, scope := range e.staticScopes {
		for _, v := range scope {
			if visit(v) {
				return true
			}
		}
	}
	for _, class := range e.classes {
		for _, slot := range class.staticProperties {
			if visit(slot.value) {
				return true
			}
		}
	}
	return false
}

func (e *evaluator) arrayCopyForWrite(a *phpArray) *phpArray {
	b := *a
	b.entries = append([]arrayEntry(nil), a.entries...)
	for i, entry := range b.entries {
		if entry.value.kind == refKind && !e.referenceSharedOutside(a, entry.key, entry.value.ref) {
			b.entries[i].value = dereference(entry.value)
		}
	}
	return &b
}

func (e *evaluator) arrayWith(a *phpArray, key arrayKey, v value) *phpArray {
	return e.arrayCopyForWrite(a).with(key, v)
}

func (e *evaluator) arrayWithAssigned(a *phpArray, key arrayKey, v value, line int) (*phpArray, error) {
	if current, ok := a.getRaw(key); ok && current.kind == refKind && e.referenceSharedOutside(a, key, current.ref) {
		if err := e.writeReference(current.ref, v, line); err != nil {
			return nil, err
		}
		return a, nil
	}
	return e.arrayWith(a, key, v), nil
}

func referenceTargetSyntax(x expr) bool {
	switch n := x.(type) {
	case variable:
		return true
	case dynamicVariable:
		return true
	case offsetExpr:
		return referenceTargetSyntax(n.base)
	case propertyExpr, staticPropertyExpr, nullsafeExpr:
		return true
	default:
		return false
	}
}

func referenceSourceSyntax(x expr) bool {
	if referenceTargetSyntax(x) {
		return true
	}
	switch x.(type) {
	case functionCall, methodCall, staticCallExpr, dynamicCall:
		return true
	}
	return false
}

func (e *evaluator) referencePlace(x expr, asSource bool) (referencePlace, error) {
	switch n := x.(type) {
	case functionCall, methodCall, staticCallExpr, dynamicCall:
		if !asSource {
			return referencePlace{}, &runtimeError{expressionLine(x), "Error", "gophp: reference-returning calls cannot be reference assignment targets"}
		}
		var result value
		var err error
		switch call := n.(type) {
		case functionCall:
			result, err = e.callFunction(call)
		case methodCall:
			result, err = e.callMethod(call)
		case staticCallExpr:
			result, err = e.callStaticMethod(call)
		case dynamicCall:
			result, err = e.callDynamic(call)
		}
		if err != nil {
			return referencePlace{}, err
		}
		if result.kind != refKind {
			return referencePlace{}, &runtimeError{expressionLine(x), "Error", "gophp: call did not return a reference"}
		}
		return referencePlace{
			read:   func() value { return result },
			assign: func(v value) error { return e.writeReference(result.ref, v, expressionLine(x)) },
			replace: func(value) error {
				return &runtimeError{expressionLine(x), "Error", "gophp: cannot rebind a returned reference"}
			},
		}, nil
	case dynamicVariable:
		name, err := e.dynamicVariableName(n)
		if err != nil {
			return referencePlace{}, err
		}
		return e.referencePlace(variable{token{kind: "variable", text: name, line: n.line}}, asSource)
	case variable:
		name := n.token.text
		owner := e.variableOwner(name)
		return referencePlace{
			read:    func() value { return owner[name] },
			assign:  e.variableScopeSlot(owner, name, n.token.line).set,
			replace: func(v value) error { owner[name] = v; return nil },
		}, nil
	case propertyExpr:
		base, err := e.eval(n.base)
		if err != nil {
			return referencePlace{}, err
		}
		name, err := e.instancePropertyName(n.name, n.nameExpr, n.line)
		if err != nil {
			return referencePlace{}, err
		}
		if base.kind != objectKind {
			return referencePlace{}, &runtimeError{n.line, "Error", "Attempt to assign property \"" + name + "\" on " + base.typeName()}
		}
		object := base.o
		key, err := e.propertyKey(object, name, n.line)
		if err != nil {
			return referencePlace{}, err
		}
		property := propertyForKey(object, key)
		if property == nil && object.class != nil && object.class.readonly {
			return referencePlace{}, &runtimeError{n.line, "Error", "Cannot create dynamic property " + object.class.name + "::$" + name}
		}
		if object.enumCase != "" || property != nil && (property.readonly || property.owner.readonly) {
			return referencePlace{}, &runtimeError{n.line, "Error", "Cannot indirectly modify readonly property " + object.className() + "::$" + name}
		}
		var constraint *propertyReferenceSource
		if property != nil && property.typeName != "" {
			constraint = &propertyReferenceSource{property: *property, object: object, key: key}
		}
		if _, exists := object.properties.getRaw(key); !exists && constraint != nil && asSource {
			if !property.nullable && !typeIncludesNull(property.typeName) {
				return referencePlace{}, &runtimeError{n.line, "Error", "Cannot access uninitialized non-nullable property " + property.owner.name + "::$" + property.name + " by reference"}
			}
			object.properties = object.properties.with(key, value{kind: nullKind})
		}
		if _, exists := object.properties.getRaw(key); !exists && constraint == nil {
			if err := e.storeObjectProperty(object, key, property, value{kind: nullKind}, n.line); err != nil {
				return referencePlace{}, err
			}
			if object.dynamicProperty(name, e.scope) {
				if err := e.diagnostic("Deprecated", n.line, "Creation of dynamic property "+object.className()+"::$"+name+" is deprecated"); err != nil {
					return referencePlace{}, err
				}
			}
		}
		return referencePlace{
			read:   func() value { v, _ := object.properties.getRaw(key); return v },
			assign: func(v value) error { return e.storeObjectProperty(object, key, property, v, n.line) },
			replace: func(v value) error {
				if old, ok := object.properties.getRaw(key); ok && old.kind == refKind && constraint != nil {
					old.ref.removeSource(*constraint)
				}
				object.properties = object.properties.with(key, v)
				if v.kind == refKind && constraint != nil {
					v.ref.addSource(*constraint)
				}
				return nil
			},
			source: constraint,
		}, nil
	case staticPropertyExpr:
		prepared, err := e.prepareStaticProperty(n)
		if err != nil {
			return referencePlace{}, err
		}
		class, err := e.resolveClassName(prepared.className, n.line)
		if err != nil {
			return referencePlace{}, err
		}
		if _, err := e.staticPropertyLocation(prepared); err != nil {
			return referencePlace{}, err
		}
		slot := class.staticProperties[prepared.name]
		var constraint *propertyReferenceSource
		if slot.decl.typeName != "" {
			constraint = &propertyReferenceSource{property: slot.decl, slot: slot}
		}
		if !slot.initialized && asSource && constraint != nil {
			if !slot.decl.nullable && !typeIncludesNull(slot.decl.typeName) {
				return referencePlace{}, &runtimeError{n.line, "Error", "Cannot access uninitialized non-nullable property " + slot.decl.owner.name + "::$" + slot.decl.name + " by reference"}
			}
			slot.value, slot.initialized = value{kind: nullKind}, true
		}
		return referencePlace{
			read: func() value { return slot.value },
			assign: func(v value) error {
				if slot.value.kind == refKind {
					return e.writeReference(slot.value.ref, v, n.line)
				} else {
					slot.value = v
				}
				return nil
			},
			replace: func(v value) error {
				if slot.value.kind == refKind && constraint != nil {
					slot.value.ref.removeSource(*constraint)
				}
				slot.value, slot.initialized = v, true
				if v.kind == refKind && constraint != nil {
					v.ref.addSource(*constraint)
				}
				return nil
			},
			source: constraint,
		}, nil
	case offsetExpr:
		parent, err := e.referencePlace(n.base, true)
		if err != nil {
			return referencePlace{}, err
		}
		return e.offsetReferencePlace(parent, n)
	}
	return referencePlace{}, &runtimeError{expressionLine(x), "Error", "Only variables and array elements can be assigned by reference"}
}

// The parent place is resolved separately so callers that already evaluated
// an argument can use that same binding for all of its array elements.
func (e *evaluator) offsetReferencePlace(parent referencePlace, n offsetExpr) (referencePlace, error) {
	base := dereference(parent.read())
	if base.kind == nullKind || base.kind == boolKind && base.i == 0 {
		if base.kind == boolKind {
			if err := e.diagnostic("Deprecated", n.line, "Automatic conversion of false to array is deprecated"); err != nil {
				return referencePlace{}, err
			}
		}
		base = arrayValue()
		if err := parent.assign(base); err != nil {
			return referencePlace{}, err
		}
	}
	if base.kind != arrayKind {
		return referencePlace{}, &runtimeError{n.line, "Error", "Cannot use a scalar value as an array"}
	}
	var err error
	var key arrayKey
	if n.key == nil {
		key, err = base.a.appendKey(n.line)
	} else {
		var raw value
		raw, err = e.eval(n.key)
		if err == nil {
			key, err = e.arrayKey(raw, n.line)
		}
	}
	if err != nil {
		return referencePlace{}, err
	}
	if _, exists := base.a.getRaw(key); !exists {
		base.a = e.arrayWith(base.a, key, value{kind: nullKind})
		if err := parent.assign(base); err != nil {
			return referencePlace{}, err
		}
	}
	return referencePlace{
		read: func() value {
			latest := dereference(parent.read())
			if latest.kind != arrayKind {
				return value{}
			}
			v, _ := latest.a.getRaw(key)
			return v
		},
		assign: func(v value) error {
			latest := dereference(parent.read())
			if latest.kind != arrayKind {
				return &runtimeError{n.line, "Error", "Cannot use a scalar value as an array"}
			}
			latest.a, err = e.arrayWithAssigned(latest.a, key, v, n.line)
			if err != nil {
				return err
			}
			return parent.assign(latest)
		},
		replace: func(v value) error {
			latest := dereference(parent.read())
			if latest.kind != arrayKind {
				return &runtimeError{n.line, "Error", "Cannot use a scalar value as an array"}
			}
			latest.a = e.arrayWith(latest.a, key, v)
			return parent.assign(latest)
		},
		arrays: append(parent.arrays, base.a),
	}, nil
}

func (e *evaluator) evalReferenceAssignment(n referenceAssignment) (value, error) {
	if target, ok := n.target.(variable); ok {
		if source, ok := n.source.(variable); ok {
			return e.bindReference(target.token.text, source.token.text), nil
		}
	}
	target, err := e.referencePlace(n.target, false)
	if err != nil {
		return value{}, err
	}
	source, err := e.referencePlace(n.source, true)
	if err != nil {
		return value{}, err
	}
	return e.bindReferencePlaces(n, target, source)
}

func (e *evaluator) bindReferencePlaces(n referenceAssignment, target, source referencePlace) (value, error) {
	stored := source.read()
	if candidate := dereference(stored); candidate.kind == arrayKind {
		for _, ancestor := range target.arrays {
			if candidate.a == ancestor {
				return value{}, &runtimeError{n.line, "Error", "gophp: cyclic array references are not supported"}
			}
		}
	}
	if stored.kind != refKind {
		stored = value{kind: refKind, ref: &referenceCell{value: dereference(stored)}}
	}
	e.pruneReferenceSources(stored.ref)
	// Verify the prospective sources before changing either binding. In
	// particular, a failed typed-property bind must not leave a new alias.
	sources := append([]propertyReferenceSource(nil), stored.ref.sources...)
	addSource := func(candidate *propertyReferenceSource) {
		if candidate == nil {
			return
		}
		for _, existing := range sources {
			if existing.sameBinding(*candidate) {
				return
			}
		}
		sources = append(sources, *candidate)
	}
	addSource(source.source)
	addSource(target.source)
	converted, err := e.coerceReferenceValue(stored.ref.value, sources, n.line, true)
	if err != nil {
		return value{}, err
	}
	if source.read().kind != refKind {
		if err := source.replace(stored); err != nil {
			return value{}, err
		}
	}
	if v, ok := n.target.(variable); ok && e.bindings[v.token.text] != 0 {
		delete(e.bindings, v.token.text)
		e.variables[v.token.text] = stored
	} else if err := target.replace(stored); err != nil {
		return value{}, err
	}
	stored.ref.value = converted
	return stored.ref.value, nil
}
