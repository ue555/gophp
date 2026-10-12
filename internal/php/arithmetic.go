package php

import (
	"fmt"
	"math"
	"strings"
)

func (e *evaluator) modulo(left, right value, line int) (value, error) {
	a, ok := e.integerOperand(left, line)
	if !ok {
		return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Unsupported operand types: %s %% %s", left.typeName(), right.typeName())}
	}
	b, ok := e.integerOperand(right, line)
	if !ok {
		return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Unsupported operand types: %s %% %s", left.typeName(), right.typeName())}
	}
	if b == 0 {
		return value{}, &runtimeError{line, "DivisionByZeroError", "Modulo by zero"}
	}
	if b == -1 {
		return integer(0), nil
	}
	return integer(a % b), nil
}

func (e *evaluator) integerOperand(v value, line int) (int64, bool) {
	// stdClass and user-defined objects cannot be converted for % or %=.
	// numeric rejects them without the warning used by numeric comparisons.
	n, trailing, ok := v.numeric()
	if !ok {
		return 0, false
	}
	if trailing {
		e.warn(line, "A non-numeric value encountered")
	}
	if n.kind == intKind {
		return n.i, true
	}
	if v.kind == floatKind {
		// Integer operand conversion uses the same loss/range diagnostics as keys.
		key, _ := e.arrayKey(v, line)
		return key.index, true
	}
	// PHP saturates finite float-strings, but wraps actual float operands.
	i := int64(0)
	if !math.IsInf(n.f, 0) && !math.IsNaN(n.f) {
		switch {
		case n.f >= 9223372036854775808.0:
			i = math.MaxInt64
		case n.f < -9223372036854775808.0:
			i = math.MinInt64
		default:
			i = int64(n.f)
		}
	}
	if float64(i) != n.f {
		e.diagnostic("Deprecated", line, fmt.Sprintf("Implicit conversion from float-string \"%s\" to int loses precision", v.s))
	}
	return i, true
}

func (e *evaluator) evalCompound(node compoundAssignment) (value, error) {
	prepared, err := e.prepareTarget(node.target)
	if err != nil {
		return value{}, err
	}
	// DIM_OP fetches its writable element before a bare RHS. ASSIGN_OP and
	// OBJ_OP read the RHS first. A DIM RHS referring to its own container is
	// copied before the writable fetch, matching Zend's assign-to-self handling.
	rhs, bare := node.source.(variable)
	_, dimension := node.target.(offsetExpr)
	deferRHS := bare && dimension && rhs.token.text != rootVariable(node.target)
	var right value
	if !deferRHS {
		right, err = e.eval(node.source)
		if err != nil {
			return value{}, err
		}
	}
	target, err := e.locate(prepared, locationReadWrite)
	if err != nil {
		return value{}, err
	}
	if deferRHS {
		right, err = e.eval(node.source)
		if err != nil {
			return value{}, err
		}
	}
	v, err := e.binaryOperation(strings.TrimSuffix(node.op.kind, "="), target.get(), right, node.op.line)
	if err != nil {
		return value{}, err
	}
	return v, target.set(v)
}

func (e *evaluator) evalIncrement(node increment) (value, error) {
	prepared, err := e.prepareTarget(node.target)
	if err != nil {
		return value{}, err
	}
	target, err := e.locate(prepared, locationIncrement)
	if err != nil {
		return value{}, err
	}
	before := target.get()
	after, err := e.incrementValue(before, node.op.kind == "++", node.op.line)
	if err != nil {
		return value{}, err
	}
	if err := target.set(after); err != nil {
		return value{}, err
	}
	if node.post {
		return before, nil
	}
	return after, nil
}

func (e *evaluator) incrementValue(v value, inc bool, line int) (value, error) {
	op, action := "-", "Decrement"
	if inc {
		op, action = "+", "Increment"
	}
	switch v.kind {
	case intKind, floatKind:
		return arithmetic(op, v, integer(1)), nil
	case boolKind:
		e.warn(line, action+" on type bool has no effect, this will change in the next major version of PHP")
		return v, nil
	case nullKind:
		if inc {
			return integer(1), nil
		}
		e.warn(line, "Decrement on type null has no effect, this will change in the next major version of PHP")
		return v, nil
	case stringKind:
		if n, trailing, ok := v.numeric(); ok && !trailing {
			return arithmetic(op, n, integer(1)), nil
		}
		if !inc {
			if v.s == "" {
				e.diagnostic("Deprecated", line, "Decrement on empty string is deprecated as non-numeric")
				return integer(-1), nil
			}
			e.diagnostic("Deprecated", line, "Decrement on non-numeric string has no effect and is deprecated")
			return v, nil
		}
		if e.version == PHP85 {
			e.diagnostic("Deprecated", line, "Increment on non-numeric string is deprecated, use str_increment() instead")
		} else if !asciiAlphanumeric(v.s) {
			e.diagnostic("Deprecated", line, "Increment on non-alphanumeric string is deprecated")
		}
		return value{kind: stringKind, s: incrementString(v.s)}, nil
	default:
		name := v.typeName()
		if v.kind == objectKind {
			name = v.o.className()
		}
		return value{}, &runtimeError{line, "TypeError", "Cannot " + strings.ToLower(action) + " " + name}
	}
}

func asciiAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	for i := range s {
		c := s[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func incrementString(s string) string {
	if s == "" {
		return "1"
	}
	b := []byte(s)
	prefix := byte('1')
	for i := len(b) - 1; i >= 0; i-- {
		switch c := b[i]; {
		case c == 'z':
			b[i], prefix = 'a', 'a'
		case c == 'Z':
			b[i], prefix = 'A', 'A'
		case c == '9':
			b[i], prefix = '0', '1'
		case c >= 'a' && c < 'z', c >= 'A' && c < 'Z', c >= '0' && c < '9':
			b[i]++
			return string(b)
		default:
			return string(b)
		}
	}
	return string(prefix) + string(b)
}
