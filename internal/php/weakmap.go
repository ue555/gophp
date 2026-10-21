package php

import "strings"

// WeakMap's object-keyed storage is currently valid while both objects stay
// reachable. Reclaiming a key with no other owner requires a GC-aware pass.
func (e *evaluator) invokeWeakMap(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method WeakMap::" + method.name + "() cannot be called statically"}
	}
	if object.o.weakMap == nil {
		object.o.weakMap = make(map[*phpObject]value)
	}
	switch strings.ToLower(method.name) {
	case "count":
		if len(args) != 0 {
			return value{}, &runtimeError{line, "ArgumentCountError", "WeakMap::count() expects exactly 0 arguments"}
		}
		return integer(int64(len(object.o.weakMap))), nil
	case "getiterator":
		return value{}, &runtimeError{line, "Error", "gophp: WeakMap iteration is not implemented"}
	case "offsetexists", "offsetget", "offsetset", "offsetunset":
		need := 1
		if strings.EqualFold(method.name, "offsetSet") {
			need = 2
		}
		if len(args) != need {
			return value{}, &runtimeError{line, "ArgumentCountError", "WeakMap::" + method.name + "() expects exactly " + map[int]string{1: "1", 2: "2"}[need] + " argument(s)"}
		}
		key := args[0].value
		if key.kind != objectKind || key.o == nil {
			return value{}, builtinTypeError(line, "WeakMap::"+method.name, 1, "object", "object", key)
		}
		stored, exists := object.o.weakMap[key.o]
		switch strings.ToLower(method.name) {
		case "offsetexists":
			return boolean(exists && stored.kind != nullKind), nil
		case "offsetget":
			if !exists {
				return value{}, &runtimeError{line, "Error", "Object " + key.o.className() + " not contained in WeakMap"}
			}
			return stored, nil
		case "offsetset":
			object.o.weakMap[key.o] = args[1].value
		case "offsetunset":
			delete(object.o.weakMap, key.o)
		}
		return value{kind: nullKind}, nil
	}
	return value{}, &runtimeError{line, "Error", "gophp: unsupported WeakMap method"}
}
