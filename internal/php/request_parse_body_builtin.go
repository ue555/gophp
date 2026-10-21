package php

import "strings"

func (e *evaluator) callRequestParseBody(node functionCall) (value, error) {
	if e.version == PHP83 {
		return value{}, &runtimeError{node.line, "Error", "Call to undefined function request_parse_body()"}
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) > 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "request_parse_body() expects at most 1 argument"}
	}
	server, _ := e.readVariable("_SERVER")
	contentType := ""
	if server.kind == arrayKind && server.a != nil {
		if v, ok := server.a.get(arrayKey{isString: true, text: "CONTENT_TYPE"}); ok {
			contentType = v.text()
		}
	}
	if !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") && !strings.HasPrefix(contentType, "multipart/form-data") {
		return value{}, &runtimeError{node.line, "RequestParseBodyException", "Unsupported content type"}
	}
	// The multipart and form parser requires a separate streaming upload unit.
	return value{}, &runtimeError{node.line, "Error", "gophp: request_parse_body() form parsing is not implemented"}
}
