package php

import "fmt"

// A source stops constraining the cell after its object is no longer reachable
// from the PHP execution state. The source's Go pointer is bookkeeping, not a
// PHP ownership edge.
func (e *evaluator) referenceSourceAlive(source propertyReferenceSource) bool {
	if source.slot != nil {
		return true
	}
	seenArrays := map[*phpArray]bool{}
	seenObjects := map[*phpObject]bool{}
	seenRefs := map[*referenceCell]bool{}
	var visit func(value) bool
	visit = func(v value) bool {
		switch v.kind {
		case refKind:
			if v.ref == nil || seenRefs[v.ref] {
				return false
			}
			seenRefs[v.ref] = true
			return visit(v.ref.value)
		case arrayKind:
			if v.a == nil || seenArrays[v.a] {
				return false
			}
			seenArrays[v.a] = true
			for _, item := range v.a.entries {
				if visit(item.value) {
					return true
				}
			}
		case objectKind:
			if v.o == source.object {
				return true
			}
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

func (e *evaluator) pruneReferenceSources(cell *referenceCell) {
	kept := cell.sources[:0]
	for _, source := range cell.sources {
		if e.referenceSourceAlive(source) {
			kept = append(kept, source)
		}
	}
	cell.sources = kept
}

// A type source belongs to a property binding, not to a variable alias. A
// reference may have several sources with different declarations.
type propertyReferenceSource struct {
	property propertyDecl
	object   *phpObject
	key      arrayKey
	slot     *staticPropertySlot
}

func (s propertyReferenceSource) sameBinding(other propertyReferenceSource) bool {
	if s.slot != nil || other.slot != nil {
		return s.slot != nil && s.slot == other.slot
	}
	return s.object == other.object && s.key == other.key
}

func (cell *referenceCell) addSource(source propertyReferenceSource) {
	for _, old := range cell.sources {
		if old.sameBinding(source) {
			return
		}
	}
	cell.sources = append(cell.sources, source)
}

func (cell *referenceCell) removeSource(source propertyReferenceSource) {
	for i, old := range cell.sources {
		if old.sameBinding(source) {
			cell.sources = append(cell.sources[:i], cell.sources[i+1:]...)
			return
		}
	}
}

func (e *evaluator) coerceReferenceValue(v value, sources []propertyReferenceSource, line int, binding bool) (value, error) {
	if len(sources) == 0 {
		return v, nil
	}
	var result value
	for i, source := range sources {
		decl := source.property
		converted, selected, ok := e.coerceType(v, relativeTypeName(decl.typeName, decl.owner, nil), decl.nullable, e.strictTypes)
		if !ok {
			if binding && i == len(sources)-1 {
				if i == 0 {
					return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Cannot assign %s to property %s::$%s of type %s", v.typeName(), decl.owner.name, decl.name, declaredTypeText(decl.typeName, decl.nullable))}
				}
				first := sources[0].property
				return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Reference with value of type %s held by property %s::$%s of type %s is not compatible with property %s::$%s of type %s", v.typeName(), first.owner.name, first.name, declaredTypeText(first.typeName, first.nullable), decl.owner.name, decl.name, declaredTypeText(decl.typeName, decl.nullable))}
			}
			return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Cannot assign %s to reference held by property %s::$%s of type %s", v.typeName(), decl.owner.name, decl.name, declaredTypeText(decl.typeName, decl.nullable))}
		}
		if i != 0 && !converted.identical(result) {
			if binding {
				first := sources[0].property
				return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Reference with value of type %s held by property %s::$%s of type %s is not compatible with property %s::$%s of type %s", v.typeName(), first.owner.name, first.name, declaredTypeText(first.typeName, first.nullable), decl.owner.name, decl.name, declaredTypeText(decl.typeName, decl.nullable))}
			}
			first := sources[0].property
			return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Cannot assign %s to reference held by property %s::$%s of type %s and property %s::$%s of type %s, as this would result in an inconsistent type conversion", v.typeName(), first.owner.name, first.name, declaredTypeText(first.typeName, first.nullable), decl.owner.name, decl.name, declaredTypeText(decl.typeName, decl.nullable))}
		}
		if err := e.warnLossyIntTypeCoercion(v, selected, e.strictTypes, line); err != nil {
			// A handler exception interrupts the lossy conversion. PHP reports a
			// TypeError for the reference assignment with that exception as its
			// previous cause, rather than exposing the handler exception alone.
			if prior, ok := e.catchable(err); ok {
				// The diagnostic is no longer the pending failure: the reference
				// assignment replaces it with a TypeError. Otherwise eval's
				// deferred diagnostic propagation would restore the old exception.
				e.pendingDiagnostic = nil
				e.failureFrames = nil
				failure := &runtimeError{line, "TypeError", fmt.Sprintf("Cannot assign %s to reference held by property %s::$%s of type %s", v.typeName(), decl.owner.name, decl.name, declaredTypeText(decl.typeName, decl.nullable))}
				if replacement, ok := e.catchable(failure); ok {
					chainException(replacement.object, prior.object)
					return value{}, replacement
				}
			}
			return value{}, err
		}
		result = converted
	}
	return result, nil
}

func (e *evaluator) writeReference(cell *referenceCell, v value, line int) error {
	e.pruneReferenceSources(cell)
	converted, err := e.coerceReferenceValue(dereference(v), cell.sources, line, false)
	if err != nil {
		return err
	}
	cell.value = converted
	return nil
}
