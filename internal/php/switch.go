package php

import "fmt"

type switchCase struct {
	condition expr
	body      []statement
	line      int
}

func (p *parser) switchStatement(line int) (statement, error) {
	s := statement{kind: "switch", line: line}
	if err := p.expect("("); err != nil {
		return s, err
	}
	var err error
	s.condition, err = p.expression(0, ")")
	if err != nil {
		return s, err
	}
	if err = p.expect(")"); err != nil {
		return s, err
	}
	colonForm := p.accept(":")
	if !colonForm {
		if err = p.expect("{"); err != nil {
			return s, err
		}
	}
	p.accept(";")
	for p.current().kind != "}" && !(colonForm && p.keyword("endswitch")) {
		c := switchCase{line: p.current().line}
		if p.acceptKeyword("case") {
			c.condition, err = p.expression(0, ":", ";")
			if err != nil {
				return s, err
			}
		} else if !p.acceptKeyword("default") {
			return s, unexpected(p.current(), "case", "default")
		}
		if !p.accept(";") {
			if err = p.expect(":"); err != nil {
				return s, err
			}
		}
		for p.current().kind != "}" && !(colonForm && p.keyword("endswitch")) && !p.keyword("case") && !p.keyword("default") {
			stmt, err := p.statement()
			if err != nil {
				return s, err
			}
			c.body = append(c.body, stmt)
		}
		s.cases = append(s.cases, c)
	}
	if colonForm {
		if !p.acceptKeyword("endswitch") {
			return s, unexpected(p.current(), "endswitch")
		}
		return s, p.expect(";")
	}
	return s, p.expect("}")
}

func (p *parser) keyword(name string) bool {
	t := p.current()
	return t.kind == "name" && equalName(t.text, name)
}

func (e *evaluator) executeSwitch(s statement) error {
	subject, err := e.eval(s.condition)
	if err != nil {
		return err
	}
	selected, fallback := -1, -1
	for i, c := range s.cases {
		if c.condition == nil {
			fallback = i
			continue
		}
		v, err := e.eval(c.condition)
		if err != nil {
			return err
		}
		match, err := e.comparisonValue("==", subject, v, c.line)
		if err != nil {
			return err
		}
		if match.truthy() {
			selected = i
			break
		}
	}
	if selected < 0 {
		selected = fallback
	}
	if selected < 0 {
		return nil
	}
	for _, c := range s.cases[selected:] {
		err := e.execute(c.body)
		if control, ok := err.(*loopControl); ok {
			if control.levels == 1 {
				return nil
			}
			return &loopControl{control.kind, control.levels - 1}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// PHP emits continue-to-switch warnings at compile time, including dead code.
func (e *evaluator) controlDiagnostics(statements []statement, contexts []string) {
	for _, s := range statements {
		if s.kind == "continue" {
			level, err := loopLevel(s, len(contexts))
			if err == nil && contexts[len(contexts)-int(level)] == "switch" {
				message := "\"continue\" targeting switch is equivalent to \"break\""
				if level > 1 {
					message = fmt.Sprintf("\"continue %d\" targeting switch is equivalent to \"break %d\"", level, level)
				}
				if int(level) < len(contexts) {
					message += fmt.Sprintf(". Did you mean to use \"continue %d\"?", level+1)
				}
				e.warn(s.line, message)
			}
		}
		if s.function != nil {
			e.controlDiagnostics(s.function.body, nil)
		}
		if s.class != nil {
			for _, m := range s.class.methodOrder {
				e.controlDiagnostics(m.body, nil)
			}
		}
		child := contexts
		switch s.kind {
		case "switch", "for", "while", "do", "foreach":
			child = append(append([]string{}, contexts...), s.kind)
		}
		e.controlDiagnostics(s.thenBranch, child)
		e.controlDiagnostics(s.elseBranch, contexts)
		for _, c := range s.cases {
			e.controlDiagnostics(c.body, child)
		}
	}
}
