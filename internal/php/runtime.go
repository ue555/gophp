package php

// Variables are bindings in the active global/function/method frame. Values
// remain immutable array snapshots or shared object identities. A slot binds
// to its owning frame, rather than to whichever frame is active when used.
// References are not implemented; their future aliasing belongs at this layer.
type variableScope map[string]value

type bindingKind uint8

const (
	globalBinding bindingKind = iota + 1
	staticBinding
)

func (s variableScope) slot(name string) location {
	return location{func() value { return s[name] }, func(v value) error { s[name] = v; return nil }}
}

func (e *evaluator) variableOwner(name string) variableScope {
	switch e.bindings[name] {
	case staticBinding:
		return e.staticScope
	case globalBinding:
		return e.rootScope
	}
	return e.variables
}

func (e *evaluator) readVariable(name string) (value, bool) {
	v, ok := e.variableOwner(name)[name]
	return v, ok
}

func (e *evaluator) variableSlot(name string) location { return e.variableOwner(name).slot(name) }
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
