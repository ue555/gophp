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
	s.valueTarget, err = p.foreachTarget()
	if err != nil {
		return s, err
	}
	if p.accept("=>") {
		s.keyTarget = s.valueTarget
		s.valueTarget, err = p.foreachTarget()
		if err != nil {
			return s, err
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
	array, err := e.eval(s.condition)
	if err != nil {
		return err
	}
	if array.kind == objectKind {
		return &runtimeError{expressionLine(s.condition), "Error", "gophp: foreach over objects is not supported"}
	}
	if array.kind != arrayKind {
		e.warn(expressionLine(s.condition), fmt.Sprintf("foreach() argument must be of type array|object, %s given", array.diagnosticName()))
		return nil
	}
	// Published arrays are immutable. Retaining this array gives PHP's by-value
	// snapshot, including nested value copies and shared object identities.
	for _, entry := range array.a.entries {
		if e.bounded {
			if e.steps == 0 {
				return &runtimeError{s.line, "Error", "gophp: execution budget exhausted"}
			}
			e.steps--
		}
		if _, err := e.evalAssignment(assignment{s.valueTarget, literal{entry.value, s.line}, s.line}); err != nil {
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
