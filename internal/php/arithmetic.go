package php

import (
	"fmt"
	"math"
	"strings"
)

func (e *evaluator) modulo(left, right value, line int) (value, error) {
	a, ok, err := e.integerOperand(left, line)
	if err != nil {
		if left.kind == floatKind {
			// The failed float-to-int conversion supplies the zero value
			// observed by typed-property compound assignment.
			return integer(0), err
		}
		return value{}, err
	}
	if !ok {
		return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Unsupported operand types: %s %% %s", left.typeName(), right.typeName())}
	}
	b, ok, err := e.integerOperand(right, line)
	if err != nil {
		return value{}, err
	}
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

func (e *evaluator) integerOperand(v value, line int) (int64, bool, error) {
	// stdClass and user-defined objects cannot be converted for % or %=.
	// numeric rejects them without the warning used by numeric comparisons.
	n, trailing, ok := v.numeric()
	if !ok {
		return 0, false, nil
	}
	if trailing {
		if err := e.warn(line, "A non-numeric value encountered"); err != nil {
			return 0, false, err
		}
	}
	if n.kind == intKind {
		return n.i, true, nil
	}
	if v.kind == floatKind {
		// Integer operand conversion uses the same loss/range diagnostics as keys.
		key, err := e.arrayKey(v, line)
		if err != nil {
			return 0, false, err
		}
		return key.index, true, nil
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
		if err := e.diagnostic("Deprecated", line, fmt.Sprintf("Implicit conversion from float-string \"%s\" to int loses precision", v.s)); err != nil {
			return 0, false, err
		}
	}
	return i, true, nil
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
	op := strings.TrimSuffix(node.op.kind, "=")
	left := target.get()
	v, err := e.binaryOperation(op, left, right, node.op.line)
	if err != nil {
		// Zend type-checks the temporary result of an arithmetic assignment
		// to a typed property even when these operations raise an exception.
		// Ordinary variables and array elements keep their original values.
		// Keep this to the cases checked against PHP; bitwise assignments have
		// different operators and are not covered by that evidence.
		if target.typed && compoundErrorStoresZero(op, left, right, v, err) {
			_ = target.set(integer(0))
		}
		return value{}, err
	}
	if err := target.set(v); err != nil {
		return value{}, err
	}
	// A typed property may coerce the operation's result while storing it.
	// The compound assignment expression yields the value actually stored.
	return target.get(), nil
}

func compoundErrorStoresZero(op string, left, right, partial value, err error) bool {
	if op == "%" && partial.kind == intKind {
		// The float-to-int warning already produced a zero temporary.
		return true
	}
	if op == "." {
		// A string left operand is concatenated in place and stays unchanged.
		return left.kind == intKind || left.kind == floatKind
	}
	if left.kind != intKind && left.kind != floatKind {
		return false
	}
	if op == "/" || op == "%" {
		runtimeErr, ok := err.(*runtimeError)
		return ok && runtimeErr.class == "DivisionByZeroError"
	}
	if op == "+" {
		if runtimeErr, ok := err.(*runtimeError); ok {
			return runtimeErr.class == "TypeError" && right.kind == arrayKind
		}
		if right.kind == stringKind {
			_, trailing, numeric := right.numeric()
			return numeric && trailing
		}
	}
	return false
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
	// Typed properties can coerce the incremented value while storing it.
	return target.get(), nil
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
		if err := e.warn(line, action+" on type bool has no effect, this will change in the next major version of PHP"); err != nil {
			return value{}, err
		}
		return v, nil
	case nullKind:
		if inc {
			return integer(1), nil
		}
		if err := e.warn(line, "Decrement on type null has no effect, this will change in the next major version of PHP"); err != nil {
			return value{}, err
		}
		return v, nil
	case stringKind:
		if n, trailing, ok := v.numeric(); ok && !trailing {
			return arithmetic(op, n, integer(1)), nil
		}
		if !inc {
			if v.s == "" {
				if err := e.diagnostic("Deprecated", line, "Decrement on empty string is deprecated as non-numeric"); err != nil {
					return value{}, err
				}
				return integer(-1), nil
			}
			if err := e.diagnostic("Deprecated", line, "Decrement on non-numeric string has no effect and is deprecated"); err != nil {
				return value{}, err
			}
			return v, nil
		}
		if e.version == PHP85 {
			if err := e.diagnostic("Deprecated", line, "Increment on non-numeric string is deprecated, use str_increment() instead"); err != nil {
				return value{}, err
			}
		} else if !asciiAlphanumeric(v.s) {
			if err := e.diagnostic("Deprecated", line, "Increment on non-alphanumeric string is deprecated"); err != nil {
				return value{}, err
			}
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
