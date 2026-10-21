package php

import "fmt"

func (e *evaluator) bitwise(op string, left, right value, line int) (value, error) {
	if (op == "&" || op == "|" || op == "^") && left.kind == stringKind && right.kind == stringKind {
		a, b := []byte(left.s), []byte(right.s)
		if len(a) > len(b) {
			a, b = b, a
		}
		out := make([]byte, len(a))
		for i := range a {
			switch op {
			case "&":
				out[i] = a[i] & b[i]
			case "|":
				out[i] = a[i] | b[i]
			case "^":
				out[i] = a[i] ^ b[i]
			}
		}
		if op == "|" {
			out = append(out, b[len(a):]...)
		}
		return value{kind: stringKind, s: string(out)}, nil
	}
	a, ok, err := e.integerOperand(left, line)
	if err != nil {
		return value{}, err
	}
	if !ok {
		return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Unsupported operand types: %s %s %s", left.typeName(), op, right.typeName())}
	}
	b, ok, err := e.integerOperand(right, line)
	if err != nil {
		return value{}, err
	}
	if !ok {
		return value{}, &runtimeError{line, "TypeError", fmt.Sprintf("Unsupported operand types: %s %s %s", left.typeName(), op, right.typeName())}
	}
	switch op {
	case "&":
		return integer(a & b), nil
	case "|":
		return integer(a | b), nil
	case "^":
		return integer(a ^ b), nil
	}
	if b < 0 {
		return value{}, &runtimeError{line, "ArithmeticError", "Bit shift by negative number"}
	}
	if op == "<<" {
		return integer(a << uint64(b)), nil
	}
	return integer(a >> uint64(b)), nil
}

func (e *evaluator) bitwiseNot(v value, line int) (value, error) {
	if v.kind == stringKind {
		b := []byte(v.s)
		for i := range b {
			b[i] = ^b[i]
		}
		return value{kind: stringKind, s: string(b)}, nil
	}
	if v.kind != intKind && v.kind != floatKind {
		return value{}, &runtimeError{line, "TypeError", "Cannot perform bitwise not on " + v.typeName()}
	}
	i, _, err := e.integerOperand(v, line)
	if err != nil {
		return value{}, err
	}
	return integer(^i), nil
}
