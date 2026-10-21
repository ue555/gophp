package php

import (
	"strconv"
	"strings"
)

// IteratorIterator and LimitIterator delegate to the same iterator protocol
// used by foreach. The wrapped object retains its own cursor and identity.
func (e *evaluator) invokeIteratorWrapper(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method " + method.owner.name + "::" + method.name + "() cannot be called statically"}
	}
	name := strings.ToLower(method.name)
	if name == "__construct" {
		limited := strings.EqualFold(method.owner.name, "LimitIterator")
		parameters := []parameter{{name: "iterator"}, {name: "class"}}
		if limited {
			parameters = []parameter{{name: "iterator"}, {name: "offset"}, {name: "limit"}}
		}
		bound, present, _, err := bindCallArguments(args, parameters, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", method.owner.name + "::__construct() expects at least 1 argument"}
		}
		if len(args) > len(parameters) {
			return value{}, &runtimeError{line, "ArgumentCountError", method.owner.name + "::__construct() expects at most " + strconv.Itoa(len(parameters)) + " arguments"}
		}
		if !limited && present[1] && bound[1].kind != nullKind {
			return value{}, &runtimeError{line, "Error", "gophp: IteratorIterator custom inner class is not implemented"}
		}
		inner := bound[0]
		expected := "Traversable"
		if limited {
			expected = "Iterator"
		}
		if inner.kind != objectKind || inner.o == nil || !classImplements(inner.o.class, expected) {
			return value{}, builtinTypeError(line, method.owner.name+"::__construct", 1, "iterator", expected, inner)
		}
		if !limited && classImplements(inner.o.class, "IteratorAggregate") {
			inner, err = e.iteratorMethod(inner, "getIterator", line)
			if err != nil {
				return value{}, err
			}
			if inner.kind != objectKind || inner.o == nil || !classImplements(inner.o.class, "Iterator") {
				return value{}, &runtimeError{line, "TypeError", "IteratorAggregate::getIterator() must return Traversable"}
			}
		}
		o := object.o
		o.iteratorInner = inner
		o.iteratorOffset, o.iteratorLimit = 0, -1
		o.iteratorPosition, o.iteratorReady = 0, false
		if limited {
			if present[1] {
				offset, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
				if !ok {
					return value{}, builtinTypeError(line, "LimitIterator::__construct", 2, "offset", "int", bound[1])
				}
				o.iteratorOffset = offset.i
			}
			if present[2] {
				limit, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
				if !ok {
					return value{}, builtinTypeError(line, "LimitIterator::__construct", 3, "limit", "int", bound[2])
				}
				o.iteratorLimit = limit.i
			}
			if o.iteratorOffset < 0 {
				return value{}, &runtimeError{line, "ValueError", "LimitIterator::__construct(): Argument #2 ($offset) must be greater than or equal to 0"}
			}
			if o.iteratorLimit < -1 {
				return value{}, &runtimeError{line, "ValueError", "LimitIterator::__construct(): Argument #3 ($limit) must be greater than or equal to -1"}
			}
		}
		return value{kind: nullKind}, nil
	}
	if len(args) != 0 && name != "seek" {
		return value{}, &runtimeError{line, "ArgumentCountError", method.owner.name + "::" + method.name + "() expects exactly 0 arguments"}
	}
	if name == "seek" {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "offset"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", "LimitIterator::seek() expects exactly 1 argument"}
		}
		position, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "LimitIterator::seek", 1, "offset", "int", bound[0])
		}
		positionText := strconv.FormatInt(position.i, 10)
		if position.i < object.o.iteratorOffset {
			return value{}, &runtimeError{line, "OutOfBoundsException", "Cannot seek to " + positionText + " which is below the offset"}
		}
		if object.o.iteratorLimit >= 0 && position.i >= object.o.iteratorOffset+object.o.iteratorLimit {
			return value{}, &runtimeError{line, "OutOfBoundsException", "Cannot seek to " + positionText + " which is beyond the limit"}
		}
		if _, err := e.iteratorWrapperMethod(object, "rewind", nil, line); err != nil {
			return value{}, err
		}
		for object.o.iteratorPosition < position.i {
			valid, err := e.iteratorWrapperMethod(object, "valid", nil, line)
			if err != nil {
				return value{}, err
			}
			if !valid.truthy() {
				return value{}, &runtimeError{line, "OutOfBoundsException", "Cannot seek to " + positionText + " which is behind the end of the iterator"}
			}
			if _, err := e.iteratorWrapperMethod(object, "next", nil, line); err != nil {
				return value{}, err
			}
		}
		return integer(object.o.iteratorPosition), nil
	}
	return e.iteratorWrapperMethod(object, method.name, nil, line)
}

func (e *evaluator) iteratorWrapperMethod(object value, name string, _ []callArgument, line int) (value, error) {
	o := object.o
	inner := o.iteratorInner
	if inner.kind != objectKind {
		return value{}, &runtimeError{line, "Error", "Object is not initialized"}
	}
	switch strings.ToLower(name) {
	case "getinneriterator":
		return inner, nil
	case "getposition":
		return integer(o.iteratorPosition), nil
	case "rewind":
		if o.iteratorLimit == 0 {
			return value{}, &runtimeError{line, "OutOfBoundsException", "Cannot seek to 0 which is behind offset " + strconv.FormatInt(o.iteratorOffset, 10) + " plus count 0"}
		}
		if _, err := e.iteratorMethod(inner, "rewind", line); err != nil {
			return value{}, err
		}
		o.iteratorPosition = 0
		for o.iteratorPosition < o.iteratorOffset {
			valid, err := e.iteratorMethod(inner, "valid", line)
			if err != nil {
				return value{}, err
			}
			if !valid.truthy() {
				break
			}
			if _, err := e.iteratorMethod(inner, "next", line); err != nil {
				return value{}, err
			}
			o.iteratorPosition++
		}
		o.iteratorReady = true
		return value{kind: nullKind}, nil
	case "valid":
		if !o.iteratorReady || o.iteratorLimit >= 0 && o.iteratorPosition-o.iteratorOffset >= o.iteratorLimit {
			return boolean(false), nil
		}
		return e.iteratorMethod(inner, "valid", line)
	case "key", "current":
		valid, err := e.iteratorWrapperMethod(object, "valid", nil, line)
		if err != nil || !valid.truthy() {
			return value{kind: nullKind}, err
		}
		return e.iteratorMethod(inner, strings.ToLower(name), line)
	case "next":
		if _, err := e.iteratorMethod(inner, "next", line); err != nil {
			return value{}, err
		}
		o.iteratorPosition++
		o.iteratorReady = true
		return value{kind: nullKind}, nil
	}
	return value{}, &runtimeError{line, "Error", "gophp: iterator wrapper method is not implemented"}
}
