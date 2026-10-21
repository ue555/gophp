package php

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	hexenc "encoding/hex"
	"hash"
	"strings"
)

func (e *evaluator) callHMAC(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	equal := strings.EqualFold(node.name, "hash_equals")
	params := []parameter{{name: "algo"}, {name: "data"}, {name: "key"}, {name: "binary"}}
	required := 3
	if equal {
		params = []parameter{{name: "known_string"}, {name: "user_string"}}
		required = 2
	}
	b, p, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < required || len(args) > len(params) {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "(): Invalid argument count"}
	}
	for i := 0; i < required; i++ {
		if !p[i] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "(): Missing argument"}
		}
		v, _, ok := e.coerceType(b[i], "string", false, e.strictTypes || equal)
		if !ok {
			return value{}, builtinTypeError(node.line, node.name, i+1, params[i].name, "string", b[i])
		}
		b[i] = v
	}
	if equal {
		return boolean(subtle.ConstantTimeCompare([]byte(b[0].s), []byte(b[1].s)) == 1), nil
	}
	var factory func() hash.Hash
	switch strings.ToLower(b[0].s) {
	case "md5":
		factory = md5.New
	case "sha1":
		factory = sha1.New
	case "sha224":
		factory = sha256.New224
	case "sha256":
		factory = sha256.New
	case "sha384":
		factory = sha512.New384
	case "sha512":
		factory = sha512.New
	default:
		return value{}, &runtimeError{node.line, "ValueError", "hash_hmac(): Argument #1 ($algo) must be a valid cryptographic hashing algorithm"}
	}
	binary := false
	if p[3] {
		v, _, ok := e.coerceType(b[3], "bool", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, node.name, 4, "binary", "bool", b[3])
		}
		binary = v.truthy()
	}
	h := hmac.New(factory, []byte(b[2].s))
	h.Write([]byte(b[1].s))
	sum := h.Sum(nil)
	s := hexenc.EncodeToString(sum)
	if binary {
		s = string(sum)
	}
	return value{kind: stringKind, s: s}, nil
}
