package php

import "strings"

func (e *evaluator) filterAdvance(object value, line int) error {
	for {
		valid, err := e.iteratorMethod(object.o.filterInner, "valid", line)
		if err != nil {
			return err
		}
		if !valid.truthy() {
			object.o.filterValid = false
			return nil
		}
		accepted, err := e.iteratorMethod(object, "accept", line)
		if err != nil {
			return err
		}
		if accepted.truthy() {
			object.o.filterValid = true
			return nil
		}
		if _, err := e.iteratorMethod(object.o.filterInner, "next", line); err != nil {
			return err
		}
	}
}

func (e *evaluator) invokeFilterIterator(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method FilterIterator::" + method.name + "() cannot be called statically"}
	}
	switch strings.ToLower(method.name) {
	case "__construct":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "iterator"}}, line)
		if err != nil {
			return value{}, err
		}
		if len(args) != 1 || !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", "FilterIterator::__construct() expects exactly 1 argument"}
		}
		v := bound[0]
		if v.kind != objectKind || v.o == nil || !classImplements(v.o.class, "Iterator") {
			return value{}, builtinTypeError(line, "FilterIterator::__construct", 1, "iterator", "Iterator", v)
		}
		object.o.filterInner = v
		return value{kind: nullKind}, nil
	case "getinneriterator":
		if object.o.filterInner.kind != objectKind {
			return value{}, &runtimeError{line, "Error", "FilterIterator is not initialized"}
		}
		return object.o.filterInner, nil
	case "rewind":
		if object.o.filterInner.kind != objectKind {
			return value{}, &runtimeError{line, "Error", "FilterIterator is not initialized"}
		}
		if _, err := e.iteratorMethod(object.o.filterInner, "rewind", line); err != nil {
			return value{}, err
		}
		return value{kind: nullKind}, e.filterAdvance(object, line)
	case "valid":
		return boolean(object.o.filterValid), nil
	case "current", "key":
		return e.iteratorMethod(object.o.filterInner, strings.ToLower(method.name), line)
	case "next":
		if _, err := e.iteratorMethod(object.o.filterInner, "next", line); err != nil {
			return value{}, err
		}
		return value{kind: nullKind}, e.filterAdvance(object, line)
	}
	return value{}, &runtimeError{line, "Error", "gophp: FilterIterator method is not implemented"}
}
