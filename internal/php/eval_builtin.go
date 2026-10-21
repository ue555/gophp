package php

import (
	"errors"
	"strconv"
)

func (e *evaluator) callEval(node functionCall) (value, error) {
	if len(node.args) != 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "eval() expects exactly 1 argument"}
	}
	code, err := e.eval(node.args[0])
	if err != nil {
		return value{}, err
	}
	text, err := e.toText(code, node.line)
	if err != nil {
		return value{}, err
	}
	filename := e.filename + "(" + strconv.Itoa(node.line) + ") : eval()'d code"
	tokens, err := lex("<?php " + text)
	if err != nil {
		var parse *syntaxError
		if errors.As(err, &parse) {
			return value{}, withSource(&runtimeError{parse.line, "ParseError", parse.message}, filename)
		}
		return value{}, withSource(err, filename)
	}
	p := parser{tokens: tokens, version: e.version, filename: filename}
	statements, err := p.program()
	if err != nil {
		var parse *syntaxError
		if errors.As(err, &parse) {
			return value{}, withSource(&runtimeError{parse.line, "ParseError", parse.message}, filename)
		}
		return value{}, withSource(err, filename)
	}
	if failure := validateStatements(statements); failure != nil {
		return value{}, withSource(failure, filename)
	}
	outer := e.filename
	e.filename = filename
	defer func() { e.filename = outer }()
	if err := e.executeProgram(statements); err != nil {
		if returned, ok := err.(*returned); ok {
			return returned.value, nil
		}
		return value{}, withSource(err, filename)
	}
	return value{kind: nullKind}, nil
}
