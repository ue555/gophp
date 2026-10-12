package php

import "fmt"

func (p *parser) loopStatement(kind string, line int) (statement, error) {
	s := statement{kind: kind, line: line}
	if kind == "break" || kind == "continue" {
		if p.current().kind != ";" {
			x, err := p.expression(0, ";")
			if err != nil {
				return s, err
			}
			s.expressions = []expr{x}
		}
		return s, p.expect(";")
	}
	var err error
	if kind == "do" {
		s.thenBranch, err = p.body()
		if err != nil {
			return s, err
		}
		if !p.acceptKeyword("while") {
			return s, unexpected(p.current(), "while")
		}
	}
	if err = p.expect("("); err != nil {
		return s, err
	}
	if kind == "for" {
		s.initial, err = p.expressionList(";")
		if err != nil {
			return s, err
		}
		s.expressions, err = p.expressionList(";")
		if err != nil {
			return s, err
		}
		s.update, err = p.expressionList(")")
	} else {
		s.condition, err = p.expression(0, ")")
		if err == nil {
			err = p.expect(")")
		}
	}
	if err != nil {
		return s, err
	}
	if s.condition != nil {
		s.conditionLine = expressionLine(s.condition)
	}
	if kind == "do" {
		return s, p.expect(";")
	}
	s.thenBranch, err = p.body()
	return s, err
}

func (p *parser) expressionList(end string) ([]expr, error) {
	var result []expr
	if p.accept(end) {
		return result, nil
	}
	for {
		x, err := p.expression(0, ",", end)
		if err != nil {
			return nil, err
		}
		result = append(result, x)
		if !p.accept(",") {
			break
		}
	}
	return result, p.expect(end, ",")
}

type loopControl struct {
	kind   string
	levels int64
}

func (s *loopControl) Error() string { return s.kind }

// Only literal depths are accepted by PHP's compiler, even in dead code.
func loopLevel(s statement, depth int) (int64, *runtimeError) {
	level := int64(1)
	if len(s.expressions) != 0 {
		x, ok := s.expressions[0].(literal)

		if !ok {
			return 0, &runtimeError{s.line, "", fmt.Sprintf("'%s' operator with non-integer operand is no longer supported", s.kind)}
		}
		if x.value.kind != intKind || x.value.i < 1 {
			return 0, &runtimeError{s.line, "", fmt.Sprintf("'%s' operator accepts only positive integers", s.kind)}
		}
		level = x.value.i
	}
	if depth == 0 {
		return 0, &runtimeError{s.line, "", fmt.Sprintf("'%s' not in the 'loop' or 'switch' context", s.kind)}
	}
	if level > int64(depth) {
		return 0, &runtimeError{s.line, "", fmt.Sprintf("Cannot '%s' %d levels", s.kind, level)}
	}
	return level, nil
}

func (e *evaluator) evalList(list []expr) (value, error) {
	last := boolean(true)
	for _, x := range list {
		v, err := e.eval(x)
		if err != nil {
			return value{}, err
		}
		last = v
	}
	return last, nil
}

func (e *evaluator) executeLoop(s statement) error {
	if _, err := e.evalList(s.initial); err != nil {
		return err
	}
	for first := true; ; first = false {
		// Tests/fuzzing can opt into a budget. Normal CLI execution is unlimited.
		if e.bounded {
			if e.steps == 0 {
				return &runtimeError{s.line, "Error", "gophp: execution budget exhausted"}
			}
			e.steps--
		}
		if s.kind != "do" || !first {
			var condition value
			var err error
			line := s.conditionLine
			if s.kind == "for" {
				condition, err = e.evalList(s.expressions)
				if len(s.expressions) > 0 {
					line = expressionLine(s.expressions[len(s.expressions)-1])
				}
			} else {
				condition, err = e.eval(s.condition)
			}
			if err != nil {
				return err
			}
			if !e.toBool(condition, line) {
				return nil
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
		if _, err := e.evalList(s.update); err != nil {
			return err
		}
	}
}
