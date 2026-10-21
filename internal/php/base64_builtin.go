package php

import (
	"encoding/base64"
	"strings"
)

func (e *evaluator) callBase64(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	name := strings.ToLower(node.name)
	parameters := []parameter{{name: "string"}}
	if name == "base64_decode" {
		parameters = append(parameters, parameter{name: "strict"})
	}
	bound, present, _, err := bindCallArguments(args, parameters, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > len(parameters) {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects 1 to 2 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, name, 1, "string", "string", bound[0])
	}
	if name == "base64_encode" {
		return value{kind: stringKind, s: base64.StdEncoding.EncodeToString([]byte(s.s))}, nil
	}
	input := s.s
	strict := present[1] && bound[1].truthy()
	if !strict {
		input = strings.Map(func(r rune) rune {
			if strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=", r) {
				return r
			}
			return -1
		}, input)
	}
	decoded, decodeErr := base64.StdEncoding.DecodeString(input)
	if decodeErr != nil {
		decoded, decodeErr = base64.RawStdEncoding.DecodeString(strings.TrimRight(input, "="))
	}
	if decodeErr != nil {
		return boolean(false), nil
	}
	return value{kind: stringKind, s: string(decoded)}, nil
}
