package php

import (
	"crypto/rand"
	"math/big"
)

func (e *evaluator) callRandomInt(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	b, p, _, err := bindCallArguments(args, []parameter{{name: "min"}, {name: "max"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !p[0] || !p[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "random_int() expects exactly 2 arguments"}
	}
	for i, n := range []string{"min", "max"} {
		v, _, ok := e.coerceType(b[i], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "random_int", i+1, n, "int", b[i])
		}
		b[i] = v
	}
	if b[0].i > b[1].i {
		return value{}, &runtimeError{node.line, "ValueError", "random_int(): Argument #1 ($min) must be less than or equal to argument #2 ($max)"}
	}
	low := big.NewInt(b[0].i)
	size := new(big.Int).Sub(big.NewInt(b[1].i), low)
	size.Add(size, big.NewInt(1))
	r, err := rand.Int(rand.Reader, size)
	if err != nil {
		return value{}, &runtimeError{node.line, "Error", "random_int(): Failed to obtain cryptographically secure random bytes"}
	}
	return integer(r.Add(r, low).Int64()), nil
}
