package php

func (e *evaluator) callIconvSlice(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	lengthMode := node.name == "iconv_strlen"
	params := []parameter{{name: "string"}, {name: "encoding"}}
	if !lengthMode {
		params = []parameter{{name: "string"}, {name: "offset"}, {name: "length"}, {name: "encoding"}}
	}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	min := 1
	if !lengthMode {
		min = 2
	}
	if len(args) < min || len(args) > len(params) || !present[0] || !lengthMode && !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() has an invalid argument count"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, node.name, 1, "string", "string", bound[0])
	}
	encodingIndex := 1
	if !lengthMode {
		encodingIndex = 3
	}
	encoding := "UTF-8"
	if present[encodingIndex] && bound[encodingIndex].kind != nullKind {
		v, _, ok := e.coerceType(bound[encodingIndex], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, node.name, encodingIndex+1, "encoding", "?string", bound[encodingIndex])
		}
		encoding = v.s
	}
	if normalizedEncoding(encoding) != "UTF8" {
		return value{}, &runtimeError{node.line, "Error", "gophp: iconv encoding is not implemented: " + encoding}
	}
	runes := []rune(s.s)
	if lengthMode {
		return integer(int64(len(runes))), nil
	}
	offset, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, node.name, 2, "offset", "int", bound[1])
	}
	start := int(offset.i)
	if start < 0 {
		start += len(runes)
	}
	if start < 0 {
		start = 0
	}
	if start > len(runes) {
		return boolean(false), nil
	}
	end := len(runes)
	if present[2] && bound[2].kind != nullKind {
		length, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, node.name, 3, "length", "?int", bound[2])
		}
		if length.i >= 0 {
			end = start + int(length.i)
		} else {
			end = len(runes) + int(length.i)
		}
	}
	if end < start {
		end = start
	}
	if end > len(runes) {
		end = len(runes)
	}
	return value{kind: stringKind, s: string(runes[start:end])}, nil
}
