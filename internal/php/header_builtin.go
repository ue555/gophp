package php

import "strings"

func (e *evaluator) callHeaderBuiltin(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	switch name {
	case "headers_sent":
		if len(args) > 2 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "headers_sent() expects at most 2 arguments"}
		}
		if len(args) > 0 {
			return value{}, &runtimeError{node.line, "Error", "gophp: headers_sent() output references are not implemented"}
		}
		return boolean(false), nil
	case "headers_list":
		if len(args) != 0 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "headers_list() expects exactly 0 arguments"}
		}
		out := arrayValue()
		for _, h := range e.responseHeaders {
			out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, value{kind: stringKind, s: h})
		}
		return out, nil
	case "header":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "header"}, {name: "replace"}, {name: "response_code"}}, node.line)
		if err != nil {
			return value{}, err
		}
		if len(args) < 1 || len(args) > 3 || !present[0] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "header() expects 1 to 3 arguments"}
		}
		h, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "header", 1, "header", "string", bound[0])
		}
		if strings.ContainsAny(h.s, "\r\n") {
			return value{}, &runtimeError{node.line, "ValueError", "Header may not contain more than a single header, new line detected"}
		}
		replace := !present[1] || bound[1].truthy()
		if present[2] {
			status, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "header", 3, "response_code", "int", bound[2])
			}
			if status.i > 0 {
				e.responseStatus = status.i
			}
		}
		if replace {
			prefix := strings.ToLower(strings.SplitN(h.s, ":", 2)[0])
			var kept []string
			for _, old := range e.responseHeaders {
				if strings.ToLower(strings.SplitN(old, ":", 2)[0]) != prefix {
					kept = append(kept, old)
				}
			}
			e.responseHeaders = kept
		}
		e.responseHeaders = append(e.responseHeaders, h.s)
		return value{kind: nullKind}, nil
	case "header_remove":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "name"}}, node.line)
		if err != nil {
			return value{}, err
		}
		if len(args) > 1 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "header_remove() expects at most 1 argument"}
		}
		if !present[0] {
			e.responseHeaders = nil
			return value{kind: nullKind}, nil
		}
		key, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "header_remove", 1, "name", "?string", bound[0])
		}
		var kept []string
		for _, old := range e.responseHeaders {
			if !strings.EqualFold(strings.SplitN(old, ":", 2)[0], key.s) {
				kept = append(kept, old)
			}
		}
		e.responseHeaders = kept
		return value{kind: nullKind}, nil
	case "http_response_code":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "response_code"}}, node.line)
		if err != nil {
			return value{}, err
		}
		if len(args) > 1 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "http_response_code() expects at most 1 argument"}
		}
		previous := e.responseStatus
		if previous == 0 {
			previous = 200
		}
		if present[0] {
			v, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "http_response_code", 1, "response_code", "int", bound[0])
			}
			e.responseStatus = v.i
		}
		return integer(previous), nil
	}
	return value{}, &runtimeError{node.line, "Error", "Unknown header function"}
}
