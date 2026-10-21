package php

import "fmt"

type ternaryExpr struct {
	condition, yes, no expr
	line               int
	grouped            bool
}
type coalesceAssignment struct {
	target, source expr
	line           int
}
type nullsafeExpr struct {
	base     expr
	name     string
	nameExpr expr
	args     []expr
	method   bool
	line     int
}
type matchArm struct {
	conditions []expr
	result     expr
	line       int
}
type matchExpr struct {
	subject expr
	arms    []matchArm
	line    int
}

func (ternaryExpr) expression()        {}
func (coalesceAssignment) expression() {}
func (nullsafeExpr) expression()       {}
func (matchExpr) expression()          {}

func (e *evaluator) evalTernary(n ternaryExpr) (value, error) {
	v, err := e.eval(n.condition)
	if err != nil {
		return value{}, err
	}
	if e.toBool(v, expressionLine(n.condition)) {
		if n.yes == nil {
			return v, nil
		}
		return e.eval(n.yes)
	}
	return e.eval(n.no)
}

// Cache the access expressions for both the quiet fetch and the eventual write.
// The root variable remains a location: RHS mutations must not overwrite an old
// snapshot of its array. Object receivers retain their evaluated identity.
func (e *evaluator) prepareCoalesceTarget(x expr) (expr, error) {
	switch n := x.(type) {
	case variable:
		return n, nil
	case staticPropertyExpr:
		return e.prepareStaticProperty(n)
	case offsetExpr:
		base, err := e.prepareCoalesceTarget(n.base)
		if err != nil {
			return nil, err
		}
		key, err := e.eval(n.key)
		if err != nil {
			return nil, err
		}
		n.base, n.key = base, literal{key, expressionLine(n.key)}
		return n, nil
	case propertyExpr:
		base, err := e.quietAccess(n.base, false)
		if err != nil {
			return nil, err
		}
		n.base = literal{base, expressionLine(n.base)}
		name, err := e.instancePropertyName(n.name, n.nameExpr, n.line)
		if err != nil {
			return nil, err
		}
		n.name, n.nameExpr = name, nil
		return n, nil
	default:
		v, err := e.eval(x)
		return literal{v, expressionLine(x)}, err
	}
}

func (e *evaluator) evalCoalesceAssignment(n coalesceAssignment) (value, error) {
	target, err := e.prepareCoalesceTarget(n.target)
	if err != nil {
		return value{}, err
	}
	v, err := e.quietRead(target)
	if err != nil || v.kind != nullKind {
		return v, err
	}
	v, err = e.eval(n.source)
	if err != nil {
		return value{}, err
	}
	loc, err := e.locate(target)
	if err != nil {
		return value{}, err
	}
	if err := loc.set(v); err != nil {
		return value{}, err
	}
	return loc.get(), nil
}

func hasNullsafe(x expr) bool {
	switch n := x.(type) {
	case nullsafeExpr:
		return true
	case propertyExpr:
		return hasNullsafe(n.base)
	case methodCall:
		return hasNullsafe(n.receiver)
	case offsetExpr:
		return hasNullsafe(n.base)
	}
	return false
}

// Quiet reads used by isset, empty and ?? still need the chain-level stop
// produced by ?->. The ordinary quiet accessor has no way to report that
// state, so a null result there would incorrectly evaluate later dimensions.
func (e *evaluator) quietNullsafeChain(x expr, final, methodReceiver bool) (value, bool, error) {
	var base expr
	switch n := x.(type) {
	case nullsafeExpr:
		base = n.base
	case propertyExpr:
		base = n.base
	case methodCall:
		base = n.receiver
	case offsetExpr:
		base = n.base
	default:
		// The root of a nullsafe call is still fetched quietly. In
		// particular, an undefined variable does not emit a warning.
		if methodReceiver {
			if _, variableRoot := x.(variable); !variableRoot {
				v, err := e.eval(x)
				return v, false, err
			}
		}
		v, err := e.quietAccess(x, final)
		return v, false, err
	}
	// A property or offset which supplies a method receiver must undergo
	// ordinary access checks. A quiet existence check applies only to the
	// terminal value being checked, not to the receiver needed to call it.
	baseIsReceiver := methodReceiver
	if n, ok := x.(nullsafeExpr); ok && n.method {
		baseIsReceiver = true
	}
	if _, ok := x.(methodCall); ok {
		baseIsReceiver = true
	}
	var v value
	var stopped bool
	var err error
	if hasNullsafe(base) {
		v, stopped, err = e.quietNullsafeChain(base, false, baseIsReceiver)
	} else if baseIsReceiver {
		if _, variableRoot := base.(variable); variableRoot {
			v, err = e.quietAccess(base, false)
		} else {
			v, err = e.eval(base)
		}
	} else {
		v, err = e.quietAccess(base, false)
	}
	if stopped || err != nil {
		return v, stopped, err
	}
	switch n := x.(type) {
	case nullsafeExpr:
		if v.kind == nullKind {
			return value{}, true, nil
		}
		if n.method {
			v, err = e.callMethod(methodCall{receiver: literal{v, expressionLine(base)}, name: n.name, nameExpr: n.nameExpr, args: n.args, line: n.line})
		} else if methodReceiver {
			v, err = e.readProperty(propertyExpr{base: literal{v, expressionLine(base)}, name: n.name, line: n.line, nameExpr: n.nameExpr})
		} else {
			v, err = e.quietAccess(propertyExpr{base: literal{v, expressionLine(base)}, name: n.name, line: n.line, nameExpr: n.nameExpr}, final)
		}
	case propertyExpr:
		access := propertyExpr{base: literal{v, expressionLine(base)}, name: n.name, line: n.line, nameExpr: n.nameExpr}
		if methodReceiver {
			v, err = e.readProperty(access)
		} else {
			v, err = e.quietAccess(access, final)
		}
	case methodCall:
		v, err = e.callMethod(methodCall{receiver: literal{v, expressionLine(base)}, name: n.name, nameExpr: n.nameExpr, args: n.args, line: n.line})
	case offsetExpr:
		key, keyErr := e.eval(n.key)
		if keyErr != nil {
			return value{}, false, keyErr
		}
		access := offsetExpr{literal{v, expressionLine(base)}, literal{key, expressionLine(n.key)}, n.line}
		if methodReceiver {
			v, err = e.readOffset(access)
		} else {
			v, err = e.quietAccess(access, final)
		}
	}
	return v, false, err
}

// A nullsafe short circuit propagates through the access chain, but never out
// of an argument, key, or other independent subexpression.
func (e *evaluator) evalChain(x expr) (value, bool, error) {
	var base expr
	switch n := x.(type) {
	case nullsafeExpr:
		base = n.base
	case propertyExpr:
		base = n.base
	case methodCall:
		base = n.receiver
	case offsetExpr:
		base = n.base
	default:
		v, err := e.eval(x)
		return v, false, err
	}
	v, stopped, err := e.evalChain(base)
	if stopped || err != nil {
		return v, stopped, err
	}
	if n, ok := x.(nullsafeExpr); ok {
		if v.kind == nullKind {
			return value{}, true, nil
		}
		if n.method {
			x = methodCall{receiver: literal{v, expressionLine(base)}, name: n.name, nameExpr: n.nameExpr, args: n.args, line: n.line}
		} else {
			x = propertyExpr{base: literal{v, expressionLine(base)}, name: n.name, line: n.line, nameExpr: n.nameExpr}
		}
	}
	switch n := x.(type) {
	case propertyExpr:
		n.base = literal{v, expressionLine(base)}
		v, err = e.readProperty(n)
	case methodCall:
		n.receiver = literal{v, expressionLine(base)}
		v, err = e.callMethod(n)
	case offsetExpr:
		n.base = literal{v, expressionLine(base)}
		v, err = e.readOffset(n)
	}
	return v, false, err
}

func (p *parser) matchExpression(line int) (expr, error) {
	n := matchExpr{line: line}
	if err := p.expect("("); err != nil {
		return nil, err
	}
	var err error
	n.subject, err = p.expression(0, ")")
	if err != nil {
		return nil, err
	}
	if err = p.expect(")"); err != nil {
		return nil, err
	}
	if err = p.expect("{"); err != nil {
		return nil, err
	}
	for p.current().kind != "}" {
		arm := matchArm{line: p.current().line}
		if !p.acceptKeyword("default") {
			for {
				x, err := p.expression(0, ",", "=>")
				if err != nil {
					return nil, err
				}
				arm.conditions = append(arm.conditions, x)
				if !p.accept(",") || p.current().kind == "=>" {
					break
				}
			}
		}
		if err = p.expect("=>"); err != nil {
			return nil, err
		}
		arm.result, err = p.expression(0, ",", "}")
		if err != nil {
			return nil, err
		}
		n.arms = append(n.arms, arm)
		if !p.accept(",") {
			break
		}
	}
	return n, p.expect("}")
}

func (e *evaluator) evalMatch(n matchExpr) (value, error) {
	subject, err := e.eval(n.subject)
	if err != nil {
		return value{}, err
	}
	var fallback expr
	for _, arm := range n.arms {
		if len(arm.conditions) == 0 {
			fallback = arm.result
			continue
		}
		for _, condition := range arm.conditions {
			v, err := e.eval(condition)
			if err != nil {
				return value{}, err
			}
			if subject.identical(v) {
				return e.eval(arm.result)
			}
		}
	}
	if fallback != nil {
		return e.eval(fallback)
	}
	label := traceValue(subject)
	if subject.kind == objectKind {
		label = "of type " + subject.o.className()
	}
	if subject.kind == arrayKind {
		label = "of type array"
	}
	return value{}, &runtimeError{n.line, "UnhandledMatchError", fmt.Sprintf("Unhandled match case %s", label)}
}
