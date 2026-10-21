package php

import (
	encodinghex "encoding/hex"
	"math/big"
)

func (e *evaluator) callHex2bin(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "hex2bin() expects exactly 1 argument"}
	}
	v, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "hex2bin", 1, "string", "string", bound[0])
	}
	if len(v.s)%2 != 0 {
		if err := e.warn(node.line, "hex2bin(): Hexadecimal input string must have an even length"); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	decoded, err := encodinghex.DecodeString(v.s)
	if err != nil {
		if err := e.warn(node.line, "hex2bin(): Input string must be hexadecimal string"); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	return value{kind: stringKind, s: string(decoded)}, nil
}

func (e *evaluator) callHexdec(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "hex_string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "hexdec() expects exactly 1 argument"}
	}
	v, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "hexdec", 1, "hex_string", "string", bound[0])
	}
	n := new(big.Int)
	invalid := false
	for _, c := range v.s {
		var digit int64
		switch {
		case c >= '0' && c <= '9':
			digit = int64(c - '0')
		case c >= 'a' && c <= 'f':
			digit = int64(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			digit = int64(c - 'A' + 10)
		default:
			invalid = true
			continue
		}
		n.Mul(n, big.NewInt(16))
		n.Add(n, big.NewInt(digit))
	}
	if invalid {
		if err := e.diagnostic("Deprecated", node.line, "Invalid characters passed for attempted conversion, these have been ignored"); err != nil {
			return value{}, err
		}
	}
	if n.IsInt64() {
		return integer(n.Int64()), nil
	}
	f, _ := new(big.Float).SetInt(n).Float64()
	return floating(f), nil
}
