package php

import (
	"sort"
	"strings"
)

func copyArrayObjectData(data *phpArray) *phpArray {
	if data == nil {
		return &phpArray{}
	}
	out := *data
	out.entries = append([]arrayEntry(nil), data.entries...)
	return &out
}

func (e *evaluator) invokeArrayObject(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method ArrayObject::" + method.name + "() cannot be called statically"}
	}
	name := strings.ToLower(method.name)
	if name == "__construct" {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "flags"}, {name: "iteratorClass"}}, line)
		if err != nil {
			return value{}, err
		}
		if len(args) > 3 {
			return value{}, &runtimeError{line, "ArgumentCountError", "ArrayObject::__construct() expects at most 3 arguments"}
		}
		items := arrayValue()
		if present[0] {
			if bound[0].kind != arrayKind {
				return value{}, builtinTypeError(line, "ArrayObject::__construct", 1, "array", "array|object", bound[0])
			}
			items = bound[0]
		}
		object.o.iteratorArray = copyArrayObjectData(items.a)
		return value{kind: nullKind}, nil
	}
	if object.o.iteratorArray == nil {
		object.o.iteratorArray = &phpArray{}
	}
	data := object.o.iteratorArray
	switch name {
	case "rewind":
		object.o.iteratorIndex = 0
		return value{kind: nullKind}, nil
	case "valid":
		return boolean(object.o.iteratorIndex < len(data.entries)), nil
	case "key":
		if object.o.iteratorIndex >= len(data.entries) {
			return value{kind: nullKind}, nil
		}
		key := data.entries[object.o.iteratorIndex].key
		if key.isString {
			return value{kind: stringKind, s: key.text}, nil
		}
		return integer(key.index), nil
	case "current":
		if object.o.iteratorIndex >= len(data.entries) {
			return value{kind: nullKind}, nil
		}
		return dereference(data.entries[object.o.iteratorIndex].value), nil
	case "next":
		object.o.iteratorIndex++
		return value{kind: nullKind}, nil
	case "count":
		return integer(int64(len(data.entries))), nil
	case "getarraycopy":
		return value{kind: arrayKind, a: copyArrayObjectData(data)}, nil
	case "exchangearray":
		if len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "ArrayObject::exchangeArray() expects exactly 1 argument"}
		}
		if args[0].value.kind != arrayKind {
			return value{}, builtinTypeError(line, "ArrayObject::exchangeArray", 1, "array", "array|object", args[0].value)
		}
		old := value{kind: arrayKind, a: copyArrayObjectData(data)}
		object.o.iteratorArray = copyArrayObjectData(args[0].value.a)
		return old, nil
	case "ksort":
		if len(args) > 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "ArrayObject::ksort() expects at most 1 argument"}
		}
		items := copyArrayObjectData(data)
		sort.SliceStable(items.entries, func(i, j int) bool {
			a, b := items.entries[i].key, items.entries[j].key
			if !a.isString && !b.isString {
				return a.index < b.index
			}
			if a.isString != b.isString {
				return !a.isString
			}
			return a.text < b.text
		})
		object.o.iteratorArray = items
		return boolean(true), nil
	case "getiterator":
		class, err := e.resolveClassName("ArrayIterator", line)
		if err != nil {
			return value{}, err
		}
		return value{kind: objectKind, o: &phpObject{class: class, properties: &phpArray{}, iteratorArray: copyArrayObjectData(data)}}, nil
	case "append":
		if len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "ArrayObject::append() expects exactly 1 argument"}
		}
		key, err := data.appendKey(line)
		if err != nil {
			return value{}, err
		}
		object.o.iteratorArray = data.with(key, args[0].value)
		return value{kind: nullKind}, nil
	case "offsetexists", "offsetget", "offsetunset":
		if len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "ArrayObject::" + method.name + "() expects exactly 1 argument"}
		}
		key, err := e.arrayKey(args[0].value, line)
		if err != nil {
			return value{}, err
		}
		v, exists := data.get(key)
		if name == "offsetexists" {
			return boolean(exists && v.kind != nullKind), nil
		}
		if name == "offsetget" {
			if exists {
				return dereference(v), nil
			}
			return value{kind: nullKind}, nil
		}
		object.o.iteratorArray = data.without(key)
		return value{kind: nullKind}, nil
	case "offsetset":
		if len(args) != 2 {
			return value{}, &runtimeError{line, "ArgumentCountError", "ArrayObject::offsetSet() expects exactly 2 arguments"}
		}
		var key arrayKey
		var err error
		if args[0].value.kind == nullKind {
			key, err = data.appendKey(line)
		} else {
			key, err = e.arrayKey(args[0].value, line)
		}
		if err != nil {
			return value{}, err
		}
		object.o.iteratorArray = data.with(key, args[1].value)
		return value{kind: nullKind}, nil
	}
	return value{}, &runtimeError{line, "Error", "Call to undefined method ArrayObject::" + method.name + "()"}
}
