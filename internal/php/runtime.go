package php

import "sort"

// Variables are bindings in the active global/function/method frame. Values
// remain immutable array snapshots or shared object identities. A slot binds
// to its owning frame, rather than to whichever frame is active when used.
// A reference cell is shared by variable names and array elements. It is never a PHP value:
// reads unwrap it, while ordinary assignments write through it.
type variableScope map[string]value

type referenceCell struct {
	value   value
	sources []propertyReferenceSource
}

func dereference(v value) value {
	if v.kind == refKind {
		return v.ref.value
	}
	return v
}

type bindingKind uint8

const (
	globalBinding bindingKind = iota + 1
	staticBinding
)

func (e *evaluator) variableScopeSlot(s variableScope, name string, line int) location {
	return location{func() value { return dereference(s[name]) }, func(v value) error {
		if current := s[name]; current.kind == refKind {
			return e.writeReference(current.ref, v, line)
		} else {
			s[name] = v
		}
		return nil
	}, false}
}

func (e *evaluator) variableOwner(name string) variableScope {
	if name == "_SERVER" || name == "_ENV" || name == "_GET" || name == "_POST" || name == "_FILES" || name == "_COOKIE" || name == "_REQUEST" || name == "argv" || name == "argc" {
		return e.rootScope
	}
	switch e.bindings[name] {
	case staticBinding:
		return e.staticScope
	case globalBinding:
		return e.rootScope
	}
	return e.variables
}

func (e *evaluator) readVariable(name string) (value, bool) {
	if name == "GLOBALS" {
		names := make([]string, 0, len(e.rootScope))
		for name := range e.rootScope {
			if name != "GLOBALS" {
				names = append(names, name)
			}
		}
		sort.Strings(names)
		result := arrayValue()
		for _, name := range names {
			result.a = result.a.with(arrayKey{text: name, isString: true}, dereference(e.rootScope[name]))
		}
		return result, true
	}
	v, ok := e.variableOwner(name)[name]
	return dereference(v), ok
}

func (e *evaluator) bindReference(target, source string) value {
	if target == source {
		owner := e.variableOwner(source)
		if _, exists := owner[source]; !exists {
			owner[source] = value{kind: nullKind}
		}
		return dereference(owner[source])
	}
	sourceOwner := e.variableOwner(source)
	stored, exists := sourceOwner[source]
	if !exists {
		stored = value{kind: nullKind}
	}
	if stored.kind != refKind {
		stored = value{kind: refKind, ref: &referenceCell{value: stored}}
		sourceOwner[source] = stored
	}
	if e.bindings[target] != 0 {
		// `global $target` and a static local bind this name in the current
		// frame. Rebinding the name detaches that local binding; it does not
		// replace the global or persistent static variable's own slot.
		delete(e.bindings, target)
	}
	e.variables[target] = stored
	return stored.ref.value
}

func (e *evaluator) ensureVariableReference(name string) *referenceCell {
	owner := e.variableOwner(name)
	stored, exists := owner[name]
	if !exists {
		stored = value{kind: nullKind}
	}
	if stored.kind != refKind {
		stored = value{kind: refKind, ref: &referenceCell{value: stored}}
		owner[name] = stored
	}
	return stored.ref
}

func (e *evaluator) variableSlotAt(name string, line int) location {
	return e.variableScopeSlot(e.variableOwner(name), name, line)
}
func (e *evaluator) deleteVariable(name string) {
	// unset detaches the name in this frame; the global or persistent static
	// cell remains available to be rebound by a later declaration.
	delete(e.bindings, name)
	delete(e.variables, name)
}

// One execution owns its registries. Class/function keys are case-insensitive;
// constant keys are case-sensitive. Lexical imports never enter these tables.
type runtimeSymbols struct {
	classes   map[string]*classDecl
	functions map[string]*methodDecl
	constants map[string]value
}
