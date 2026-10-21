package php

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

func (e *evaluator) callJSONDecode(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	if name == "json_last_error" || name == "json_last_error_msg" {
		args, err := e.evalArguments(node.args)
		if err != nil {
			return value{}, err
		}
		if len(args) != 0 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects exactly 0 arguments"}
		}
		if name == "json_last_error" {
			return integer(e.jsonLastError), nil
		}
		if e.jsonLastError == 0 {
			return value{kind: stringKind, s: "No error"}, nil
		}
		return value{kind: stringKind, s: "Syntax error"}, nil
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "json"}, {name: "associative"}, {name: "depth"}, {name: "flags"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 4 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "json_decode() expects 1 to 4 arguments"}
	}
	input, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "json_decode", 1, "json", "string", bound[0])
	}
	assoc := present[1] && bound[1].kind != nullKind && bound[1].truthy()
	flags := int64(0)
	if present[3] {
		v, _, ok := e.coerceType(bound[3], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "json_decode", 4, "flags", "int", bound[3])
		}
		flags = v.i
	}
	if !present[1] && flags&1 != 0 {
		assoc = true
	}
	depth := int64(512)
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "json_decode", 3, "depth", "int", bound[2])
		}
		depth = v.i
	}
	if depth < 1 {
		return value{}, &runtimeError{node.line, "ValueError", "json_decode(): Argument #3 ($depth) must be greater than 0"}
	}
	decoder := json.NewDecoder(strings.NewReader(input.s))
	decoder.UseNumber()
	v, err := e.decodeJSONValue(decoder, assoc, 0, int(depth))
	if err == nil {
		var trailing any
		if endErr := decoder.Decode(&trailing); endErr != io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}
	if err != nil {
		e.jsonLastError = 4
		if flags&4194304 != 0 {
			return value{}, &runtimeError{node.line, "JsonException", "Syntax error"}
		}
		return value{kind: nullKind}, nil
	}
	e.jsonLastError = 0
	return v, nil
}

func (e *evaluator) decodeJSONValue(decoder *json.Decoder, assoc bool, depth, limit int) (value, error) {
	if depth > limit {
		return value{}, io.ErrUnexpectedEOF
	}
	token, err := decoder.Token()
	if err != nil {
		return value{}, err
	}
	switch item := token.(type) {
	case nil:
		return value{kind: nullKind}, nil
	case bool:
		return boolean(item), nil
	case string:
		return value{kind: stringKind, s: item}, nil
	case json.Number:
		if n, err := strconv.ParseInt(string(item), 10, 64); err == nil {
			return integer(n), nil
		}
		f, err := strconv.ParseFloat(string(item), 64)
		if err != nil {
			return value{}, err
		}
		return floating(f), nil
	case json.Delim:
		if item == '[' {
			out := arrayValue()
			for decoder.More() {
				part, err := e.decodeJSONValue(decoder, assoc, depth+1, limit)
				if err != nil {
					return value{}, err
				}
				out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, part)
			}
			_, err := decoder.Token()
			return out, err
		}
		if item == '{' {
			out := arrayValue()
			for decoder.More() {
				name, err := decoder.Token()
				if err != nil {
					return value{}, err
				}
				part, err := e.decodeJSONValue(decoder, assoc, depth+1, limit)
				if err != nil {
					return value{}, err
				}
				out.a = out.a.with(arrayKey{isString: true, text: name.(string)}, part)
			}
			if _, err := decoder.Token(); err != nil {
				return value{}, err
			}
			if assoc {
				return out, nil
			}
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["stdclass"], properties: out.a}}, nil
		}
	}
	return value{}, io.ErrUnexpectedEOF
}
