package php

import "strings"

func (e *evaluator) invokeAppendIterator(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method AppendIterator::" + method.name + "() cannot be called statically"}
	}
	name := strings.ToLower(method.name)
	if name == "__construct" {
		if len(args) != 0 {
			return value{}, &runtimeError{line, "ArgumentCountError", "AppendIterator::__construct() expects exactly 0 arguments"}
		}
		object.o.appendIterators = make([]value, 0)
		object.o.appendIndex = 0
		object.o.appendReady = false
		return value{kind: nullKind}, nil
	}
	if name == "append" {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "iterator"}}, line)
		if err != nil {
			return value{}, err
		}
		if len(args) != 1 || !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", "AppendIterator::append() expects exactly 1 argument"}
		}
		inner := bound[0]
		if inner.kind != objectKind || inner.o == nil || !classImplements(inner.o.class, "Iterator") {
			return value{}, builtinTypeError(line, "AppendIterator::append", 1, "iterator", "Iterator", inner)
		}
		object.o.appendIterators = append(object.o.appendIterators, inner)
		return value{kind: nullKind}, nil
	}
	if len(args) != 0 {
		return value{}, &runtimeError{line, "ArgumentCountError", "AppendIterator::" + method.name + "() expects exactly 0 arguments"}
	}
	return e.appendIteratorMethod(object, name, line)
}

// Move past exhausted inner iterators, rewinding each one when it first
// becomes active. Keys and cursor state belong to the original iterators.
func (e *evaluator) advanceAppendIterator(object *phpObject, line int) error {
	for object.appendIndex < len(object.appendIterators) {
		valid, err := e.iteratorMethod(object.appendIterators[object.appendIndex], "valid", line)
		if err != nil {
			return err
		}
		if valid.truthy() {
			return nil
		}
		object.appendIndex++
		if object.appendIndex < len(object.appendIterators) {
			if _, err := e.iteratorMethod(object.appendIterators[object.appendIndex], "rewind", line); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *evaluator) appendIteratorMethod(object value, name string, line int) (value, error) {
	o := object.o
	switch strings.ToLower(name) {
	case "getarrayiterator":
		items := &phpArray{}
		for index, inner := range o.appendIterators {
			items = items.with(arrayKey{index: int64(index)}, inner)
		}
		return value{kind: objectKind, o: &phpObject{class: e.classes["arrayiterator"], properties: &phpArray{}, iteratorArray: items}}, nil
	case "getinneriterator", "getiteratorindex":
		if !o.appendReady || o.appendIndex >= len(o.appendIterators) {
			return value{kind: nullKind}, nil
		}
		if strings.EqualFold(name, "getiteratorindex") {
			return integer(int64(o.appendIndex)), nil
		}
		return o.appendIterators[o.appendIndex], nil
	case "rewind":
		o.appendIndex = 0
		o.appendReady = true
		if len(o.appendIterators) == 0 {
			return value{kind: nullKind}, nil
		}
		if _, err := e.iteratorMethod(o.appendIterators[0], "rewind", line); err != nil {
			return value{}, err
		}
		return value{kind: nullKind}, e.advanceAppendIterator(o, line)
	case "valid":
		return boolean(o.appendReady && o.appendIndex < len(o.appendIterators)), nil
	case "key", "current":
		if !o.appendReady || o.appendIndex >= len(o.appendIterators) {
			return value{kind: nullKind}, nil
		}
		return e.iteratorMethod(o.appendIterators[o.appendIndex], strings.ToLower(name), line)
	case "next":
		if !o.appendReady || o.appendIndex >= len(o.appendIterators) {
			return value{kind: nullKind}, nil
		}
		if _, err := e.iteratorMethod(o.appendIterators[o.appendIndex], "next", line); err != nil {
			return value{}, err
		}
		return value{kind: nullKind}, e.advanceAppendIterator(o, line)
	}
	return value{}, &runtimeError{line, "Error", "gophp: AppendIterator method is not implemented"}
}
