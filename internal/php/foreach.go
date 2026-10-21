package php

import "fmt"

func (p *parser) foreachStatement(line int) (statement, error) {
	s := statement{kind: "foreach", line: line}
	if err := p.expect("("); err != nil {
		return s, err
	}
	var err error
	s.condition, err = p.expression(0)
	if err != nil {
		return s, err
	}
	if !p.acceptKeyword("as") {
		return s, unexpected(p.current(), "as")
	}
	firstByRef := p.accept("&")
	s.valueTarget, err = p.foreachTarget()
	if err != nil {
		return s, err
	}
	if p.accept("=>") {
		if firstByRef {
			return s, syntax(line, "Cannot use a reference as a foreach key")
		}
		s.keyTarget = s.valueTarget
		s.foreachByRef = p.accept("&")
		s.valueTarget, err = p.foreachTarget()
		if err != nil {
			return s, err
		}
	}
	if s.keyTarget == nil {
		s.foreachByRef = firstByRef
	}
	if s.foreachByRef {
		if _, ok := s.valueTarget.(variable); !ok {
			return s, syntax(line, "gophp: foreach reference target must be a variable")
		}
	}
	if err = p.expect(")"); err != nil {
		return s, err
	}
	s.thenBranch, err = p.body()
	return s, err
}

func (p *parser) foreachTarget() (expr, error) {
	t := p.current()
	if t.kind == "&" {
		return nil, syntax(t.line, "unsupported syntax: foreach by reference")
	}
	x, err := p.expression(0)
	if err == nil && !writable(x) {
		return nil, syntax(t.line, "unsupported foreach target: expected variable, array element, or property")
	}
	return x, err
}

func (e *evaluator) executeForeach(s statement) error {
	if s.foreachByRef {
		return e.executeForeachByReference(s)
	}
	array, err := e.eval(s.condition)
	if err != nil {
		return err
	}
	if array.kind == objectKind {
		return e.executeForeachIterator(s, array, 0)
	}
	if array.kind != arrayKind {
		return e.warn(expressionLine(s.condition), fmt.Sprintf("foreach() argument must be of type array|object, %s given", array.diagnosticName()))
	}
	return e.executeForeachEntries(s, array.a.entries)
}

func (e *evaluator) executeForeachEntries(s statement, entries []arrayEntry) error {
	// Published arrays are immutable. Retaining this array gives PHP's by-value
	// snapshot, including nested value copies and shared object identities.
	for _, entry := range entries {
		if e.bounded {
			if e.steps == 0 {
				return &runtimeError{s.line, "Error", "gophp: execution budget exhausted"}
			}
			e.steps--
		}
		if _, err := e.evalAssignment(assignment{s.valueTarget, literal{dereference(entry.value), s.line}, s.line}); err != nil {
			return err
		}
		if s.keyTarget != nil {
			key := integer(entry.key.index)
			if entry.key.isString {
				key = value{kind: stringKind, s: entry.key.text}
			}
			if _, err := e.evalAssignment(assignment{s.keyTarget, literal{key, s.line}, s.line}); err != nil {
				return err
			}
		}
		err := e.execute(s.thenBranch)
		if control, ok := err.(*loopControl); ok {
			if control.levels > 1 {
				return &loopControl{control.kind, control.levels - 1}
			}
			if control.kind == "break" {
				return nil
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (e *evaluator) iteratorMethod(object value, name string, line int) (value, error) {
	if object.o != nil && object.o.appendIterators != nil {
		return e.appendIteratorMethod(object, name, line)
	}
	if object.o != nil && object.o.cachingArray != nil {
		return e.cachingIteratorMethod(object, name, line)
	}
	if object.o != nil && object.o.iteratorInner.kind == objectKind {
		result, err := e.iteratorWrapperMethod(object, name, nil, line)
		if err != nil && e.failureFrames == nil {
			for _, frame := range e.frames {
				e.failureFrames = append(e.failureFrames, snapshotCallFrame(frame))
			}
			e.failureFrames = append(e.failureFrames, callFrame{class: object.o.className(), method: name, object: object, filename: e.filename, line: line})
		}
		return result, err
	}
	if object.o.iteratorArray != nil {
		it := object.o
		switch name {
		case "rewind":
			it.iteratorIndex = 0
			return value{kind: nullKind}, nil
		case "valid":
			return boolean(it.iteratorIndex < len(it.iteratorArray.entries)), nil
		case "next":
			it.iteratorIndex++
			return value{kind: nullKind}, nil
		case "key", "current":
			if it.iteratorIndex >= len(it.iteratorArray.entries) {
				return value{kind: nullKind}, nil
			}
			entry := it.iteratorArray.entries[it.iteratorIndex]
			if name == "current" {
				return dereference(entry.value), nil
			}
			if entry.key.isString {
				return value{kind: stringKind, s: entry.key.text}, nil
			}
			return integer(entry.key.index), nil
		}
	}
	if object.o.class != nil && object.o.class.name == "EmptyIterator" {
		switch name {
		case "rewind", "next":
			return value{kind: nullKind}, nil
		case "valid":
			return boolean(false), nil
		case "key", "current":
			return value{}, &runtimeError{line, "Error", "Iterator is not valid"}
		}
	}
	if object.o.generator != nil {
		g := object.o.generator
		prime := func() error {
			if g.started {
				return nil
			}
			err := g.step()
			if err != nil && e.failureFrames == nil {
				e.failureFrames = append([]callFrame(nil), g.failureFrames...)
			}
			return err
		}
		switch name {
		case "rewind":
			if g.pastFirstYield {
				return value{}, &runtimeError{line, "Exception", "Cannot rewind a generator that was already run"}
			}
			return value{kind: nullKind}, prime()
		case "valid":
			if err := prime(); err != nil {
				return value{}, err
			}
			return boolean(!g.finished), nil
		case "key":
			if err := prime(); err != nil {
				return value{}, err
			}
			if g.finished {
				return value{kind: nullKind}, nil
			}
			return g.current.key, nil
		case "current":
			if err := prime(); err != nil {
				return value{}, err
			}
			if g.finished {
				return value{kind: nullKind}, nil
			}
			return dereference(g.current.item), nil
		case "next":
			if err := prime(); err != nil {
				return value{}, err
			}
			err := g.step()
			if err != nil && e.failureFrames == nil {
				e.failureFrames = append([]callFrame(nil), g.failureFrames...)
			}
			return value{kind: nullKind}, err
		}
	}
	method := e.lookupMethod(object.o.class, name)
	if method == nil {
		return value{}, &runtimeError{line, "Error", "Call to undefined method " + object.o.className() + "::" + name + "()"}
	}
	return e.invoke(object, method, nil, line)
}

func (e *evaluator) executeForeachIterator(s statement, object value, depth int) error {
	if depth >= 16 {
		return &runtimeError{s.line, "Error", "gophp: iterator nesting limit exceeded"}
	}
	if classImplements(object.o.class, "IteratorAggregate") {
		inner, err := e.iteratorMethod(object, "getIterator", s.line)
		if err != nil {
			return err
		}
		if inner.kind != objectKind || !classImplements(inner.o.class, "Traversable") {
			return &runtimeError{s.line, "TypeError", object.o.className() + "::getIterator() must return Traversable"}
		}
		return e.executeForeachIterator(s, inner, depth+1)
	}
	if !classImplements(object.o.class, "Iterator") {
		return e.executeForeachEntries(s, e.visibleObjectProperties(object.o))
	}
	if _, err := e.iteratorMethod(object, "rewind", s.line); err != nil {
		return err
	}
	for {
		valid, err := e.iteratorMethod(object, "valid", s.line)
		if err != nil {
			return err
		}
		if !valid.truthy() {
			return nil
		}
		if e.bounded {
			if e.steps == 0 {
				return &runtimeError{s.line, "Error", "gophp: execution budget exhausted"}
			}
			e.steps--
		}
		key, err := e.iteratorMethod(object, "key", s.line)
		if err != nil {
			return err
		}
		item, err := e.iteratorMethod(object, "current", s.line)
		if err != nil {
			return err
		}
		if s.keyTarget != nil {
			if _, err := e.evalAssignment(assignment{s.keyTarget, literal{key, s.line}, s.line}); err != nil {
				return err
			}
		}
		if _, err := e.evalAssignment(assignment{s.valueTarget, literal{item, s.line}, s.line}); err != nil {
			return err
		}
		err = e.execute(s.thenBranch)
		if control, ok := err.(*loopControl); ok {
			if control.levels > 1 {
				return &loopControl{control.kind, control.levels - 1}
			}
			if control.kind == "break" {
				return nil
			}
		} else if err != nil {
			return err
		}
		if _, err := e.iteratorMethod(object, "next", s.line); err != nil {
			return err
		}
	}
}

// Ordinary objects iterate properties visible from the scope where foreach
// executes. Private slots use an internal storage name, which must never be
// exposed as the user-facing key.
func (e *evaluator) visibleObjectProperties(object *phpObject) []arrayEntry {
	if object == nil || object.properties == nil {
		return nil
	}
	entries := make([]arrayEntry, 0, len(object.properties.entries))
	for _, entry := range object.properties.entries {
		property := propertyForKey(object, entry.key)
		name := entry.key.text
		if property != nil {
			switch property.visibility {
			case "private":
				if e.scope != property.owner {
					continue
				}
			case "protected":
				if !sameHierarchy(e.scope, property.owner) {
					continue
				}
			}
			name = property.name
		}
		entries = append(entries, arrayEntry{key: arrayKey{text: name, isString: true}, value: entry.value})
	}
	return entries
}

func (e *evaluator) executeForeachByReference(s statement) error {
	arrayVariable, ok := s.condition.(variable)
	if !ok {
		if dynamic, yes := s.condition.(dynamicVariable); yes {
			name, err := e.dynamicVariableName(dynamic)
			if err != nil {
				return err
			}
			arrayVariable = variable{token{kind: "variable", text: name, line: dynamic.line}}
		} else {
			// A value-producing expression (for example, a coalesce on a
			// static array property) has its own temporary array container.
			// The reference target may outlive the loop, so keep its element
			// cells while removing only the temporary variable binding.
			v, err := e.eval(s.condition)
			if err != nil {
				return err
			}
			if v.kind == objectKind {
				return e.executeForeachIteratorByReference(s, v, 0)
			}
			if v.kind != arrayKind {
				return &runtimeError{s.line, "Error", "An iterator cannot be used with foreach by reference"}
			}
			e.foreachTempCounter++
			name := fmt.Sprintf("\x00foreach_temp_%d", e.foreachTempCounter)
			e.variables[name] = v
			defer delete(e.variables, name)
			arrayVariable = variable{token{kind: "variable", text: name, line: s.line}}
		}
	}
	if current, exists := e.readVariable(arrayVariable.token.text); exists && current.kind == objectKind {
		return e.executeForeachIteratorByReference(s, current, 0)
	}
	if _, ok := s.valueTarget.(variable); !ok {
		return &runtimeError{s.line, "Error", "gophp: foreach reference target must be a variable"}
	}
	for index := 0; ; index++ {
		current, exists := e.readVariable(arrayVariable.token.text)
		if !exists || current.kind != arrayKind || index >= len(current.a.entries) {
			return nil
		}
		if e.bounded {
			if e.steps == 0 {
				return &runtimeError{s.line, "Error", "gophp: execution budget exhausted"}
			}
			e.steps--
		}
		entry := current.a.entries[index]
		keyValue := integer(entry.key.index)
		if entry.key.isString {
			keyValue = value{kind: stringKind, s: entry.key.text}
		}
		source := offsetExpr{base: arrayVariable, key: literal{keyValue, s.line}, line: s.line}
		if _, err := e.evalReferenceAssignment(referenceAssignment{s.valueTarget, source, s.line}); err != nil {
			return err
		}
		if s.keyTarget != nil {
			if _, err := e.evalAssignment(assignment{s.keyTarget, literal{keyValue, s.line}, s.line}); err != nil {
				return err
			}
		}
		err := e.execute(s.thenBranch)
		if control, ok := err.(*loopControl); ok {
			if control.levels > 1 {
				return &loopControl{control.kind, control.levels - 1}
			}
			if control.kind == "break" {
				return nil
			}
		} else if err != nil {
			return err
		}
	}
}

func (e *evaluator) executeForeachIteratorByReference(s statement, object value, depth int) error {
	if depth >= 16 {
		return &runtimeError{s.line, "Error", "gophp: iterator nesting limit exceeded"}
	}
	if object.o == nil {
		return &runtimeError{s.line, "Error", "An iterator cannot be used with foreach by reference"}
	}
	if classImplements(object.o.class, "IteratorAggregate") {
		inner, err := e.iteratorMethod(object, "getIterator", s.line)
		if err != nil {
			return err
		}
		if inner.kind != objectKind || inner.o == nil || !classImplements(inner.o.class, "Traversable") {
			return &runtimeError{s.line, "TypeError", object.o.className() + "::getIterator() must return Traversable"}
		}
		return e.executeForeachIteratorByReference(s, inner, depth+1)
	}
	g := object.o.generator
	if g == nil {
		if !classImplements(object.o.class, "Iterator") {
			return e.executeForeachObjectByReference(s, object)
		}
		return &runtimeError{s.line, "Error", "An iterator cannot be used with foreach by reference"}
	}
	if !g.yieldsByReference {
		return &runtimeError{s.line, "Exception", "You can only iterate a generator by-reference if it declared that it yields by-reference"}
	}
	if _, err := e.iteratorMethod(object, "rewind", s.line); err != nil {
		return err
	}
	target := s.valueTarget.(variable).token.text
	for {
		valid, err := e.iteratorMethod(object, "valid", s.line)
		if err != nil {
			return err
		}
		if !valid.truthy() {
			return nil
		}
		if e.bounded {
			if e.steps == 0 {
				return &runtimeError{s.line, "Error", "gophp: execution budget exhausted"}
			}
			e.steps--
		}
		if g.current.item.kind != refKind {
			return &runtimeError{s.line, "Error", "gophp: yielding a non-variable by reference is not implemented"}
		}
		if s.keyTarget != nil {
			if _, err := e.evalAssignment(assignment{s.keyTarget, literal{g.current.key, s.line}, s.line}); err != nil {
				return err
			}
		}
		e.variableOwner(target)[target] = g.current.item
		err = e.execute(s.thenBranch)
		if control, ok := err.(*loopControl); ok {
			if control.levels > 1 {
				return &loopControl{control.kind, control.levels - 1}
			}
			if control.kind == "break" {
				return nil
			}
		} else if err != nil {
			return err
		}
		if _, err := e.iteratorMethod(object, "next", s.line); err != nil {
			return err
		}
	}
}

func (e *evaluator) executeForeachObjectByReference(s statement, object value) error {
	entries := e.visibleObjectProperties(object.o)
	for _, entry := range entries {
		if e.bounded {
			if e.steps == 0 {
				return &runtimeError{s.line, "Error", "gophp: execution budget exhausted"}
			}
			e.steps--
		}
		source := propertyExpr{base: literal{object, s.line}, name: entry.key.text, line: s.line}
		if _, err := e.evalReferenceAssignment(referenceAssignment{s.valueTarget, source, s.line}); err != nil {
			return err
		}
		if s.keyTarget != nil {
			if _, err := e.evalAssignment(assignment{s.keyTarget, literal{value{kind: stringKind, s: entry.key.text}, s.line}, s.line}); err != nil {
				return err
			}
		}
		err := e.execute(s.thenBranch)
		if control, ok := err.(*loopControl); ok {
			if control.levels > 1 {
				return &loopControl{control.kind, control.levels - 1}
			}
			if control.kind == "break" {
				return nil
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}
