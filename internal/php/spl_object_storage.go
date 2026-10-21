package php

import "strings"

func (e *evaluator) invokeObjectStorage(receiver value, method *methodDecl, args []callArgument, line int) (value, error) {
	if receiver.kind != objectKind || receiver.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method SplObjectStorage::" + method.name + "() cannot be called statically"}
	}
	o := receiver.o
	if o.objectStorage == nil {
		o.objectStorage = make(map[*phpObject]value)
	}
	name := strings.ToLower(method.name)
	if name == "count" {
		return integer(int64(len(o.objectStorage))), nil
	}
	if len(args) == 0 {
		return value{}, &runtimeError{line, "ArgumentCountError", "SplObjectStorage::" + method.name + "() expects at least 1 argument"}
	}
	key := args[0].value
	if key.kind != objectKind || key.o == nil {
		return value{}, builtinTypeError(line, "SplObjectStorage::"+method.name, 1, "object", "object", key)
	}
	switch name {
	case "contains", "offsetexists":
		_, ok := o.objectStorage[key.o]
		return boolean(ok), nil
	case "offsetget":
		v, ok := o.objectStorage[key.o]
		if !ok {
			return value{}, &runtimeError{line, "UnexpectedValueException", "Object not found"}
		}
		return v, nil
	case "attach", "offsetset":
		v := value{kind: nullKind}
		if len(args) > 1 {
			v = args[1].value
		}
		o.objectStorage[key.o] = v
		return value{kind: nullKind}, nil
	case "detach", "offsetunset":
		delete(o.objectStorage, key.o)
		return value{kind: nullKind}, nil
	}
	return value{}, &runtimeError{line, "Error", "gophp: SplObjectStorage method is not implemented"}
}
