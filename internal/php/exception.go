package php

import (
	"errors"
	"fmt"
	"strings"
)

type exceptionData struct {
	message string
	code    int64
	file    string
	line    int
	trace   string
}

type thrown struct{ object value }

func (t *thrown) Error() string { return t.object.o.className() }

func isThrowableClass(class *classDecl) bool {
	for current := class; current != nil; current = current.parent {
		if equalName(current.name, "Exception") || equalName(current.name, "Error") {
			return true
		}
	}
	return false
}

func (e *evaluator) initializeException(object value, expressions []expr, line int) error {
	args, err := e.evalArguments(expressions)
	if err != nil {
		return err
	}
	params := []parameter{{name: "message"}, {name: "code"}, {name: "previous"}}
	bound, present, _, err := bindCallArguments(args, params, line)
	if err != nil {
		return err
	}
	if len(args) > len(params) {
		return &runtimeError{line, "ArgumentCountError", fmt.Sprintf("%s::__construct() expects at most 3 arguments, %d given", object.o.className(), len(args))}
	}
	if present[0] {
		object.o.exception.message, err = e.toText(bound[0], line)
		if err != nil {
			return err
		}
	}
	if present[1] {
		if bound[1].kind != intKind {
			return &runtimeError{line, "TypeError", object.o.className() + "::__construct(): Argument #2 ($code) must be of type int, " + bound[1].diagnosticName() + " given"}
		}
		object.o.exception.code = bound[1].i
	}
	if present[2] && bound[2].kind != nullKind && (bound[2].kind != objectKind || bound[2].o.exception == nil) {
		return &runtimeError{line, "TypeError", object.o.className() + "::__construct(): Argument #3 ($previous) must be of type ?Throwable, " + bound[2].diagnosticName() + " given"}
	}
	return nil
}

func (e *evaluator) callExceptionMethod(object value, node methodCall) (value, bool, error) {
	name := strings.ToLower(node.name)
	if name != "getmessage" && name != "getcode" && name != "getfile" && name != "getline" && name != "gettraceasstring" {
		return value{}, false, nil
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, true, err
	}
	if len(args) > 0 {
		// PHP's built-in accessors ignore extra positional arguments.
		if _, _, _, err := bindCallArguments(args, nil, node.line); err != nil {
			return value{}, true, err
		}
	}
	data := object.o.exception
	switch name {
	case "getmessage":
		return value{kind: stringKind, s: data.message}, true, nil
	case "getcode":
		return integer(data.code), true, nil
	case "getfile":
		return value{kind: stringKind, s: data.file}, true, nil
	case "getline":
		return integer(int64(data.line)), true, nil
	default:
		return value{kind: stringKind, s: strings.TrimSuffix(data.trace, "\n")}, true, nil
	}
}

func (e *evaluator) throwValue(v value, line int) error {
	if v.kind != objectKind || v.o == nil || !isThrowableClass(v.o.class) {
		return &runtimeError{line, "Error", "Can only throw objects"}
	}
	if v.o.exception == nil {
		v.o.exception = &exceptionData{file: e.filename, line: line}
	}
	if v.o.exception.trace == "" {
		e.captureFailure()
		v.o.exception.trace = e.stackTrace()
	}
	return &thrown{object: v}
}

func (e *evaluator) catchable(err error) (*thrown, bool) {
	var raised *thrown
	if errors.As(err, &raised) {
		return raised, true
	}
	var failure *runtimeError
	if !errors.As(err, &failure) || failure.class == "" {
		return nil, false
	}
	e.initClasses()
	class := e.classes[strings.ToLower(failure.class)]
	if class == nil || !isThrowableClass(class) {
		return nil, false
	}
	filename := e.filename
	var sourced *sourceError
	if errors.As(err, &sourced) {
		filename = sourced.filename
	}
	e.captureFailure()
	object := value{kind: objectKind, o: &phpObject{class: class, properties: &phpArray{}, exception: &exceptionData{
		message: failure.message, file: filename, line: failure.line, trace: e.stackTrace(),
	}}}
	return &thrown{object: object}, true
}

func throwableMatches(class *classDecl, name string) bool {
	for current := class; current != nil; current = current.parent {
		if equalName(current.name, name) || classImplements(current, name) {
			return true
		}
	}
	return false
}

func (e *evaluator) executeTry(s statement) error {
	pending := e.execute(s.thenBranch)
	if pending != nil {
		if raised, ok := e.catchable(pending); ok {
			for _, clause := range s.catches {
				matches := false
				for _, name := range clause.types {
					if throwableMatches(raised.object.o.class, name) {
						matches = true
						break
					}
				}
				if !matches {
					continue
				}
				e.failureFrames = nil
				if clause.name != "" {
					if err := e.variableSlot(clause.name).set(raised.object); err != nil {
						return err
					}
				}
				pending = e.execute(clause.body)
				break
			}
		}
	}
	if s.finally != nil {
		if final := e.execute(s.finally); final != nil {
			return final
		}
	}
	return pending
}
