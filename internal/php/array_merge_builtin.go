package php

import "fmt"

func (e *evaluator) callArrayMerge(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	out := arrayValue()
	for i, arg := range args {
		if arg.named {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + arg.name}
		}
		if arg.value.kind != arrayKind {
			return value{}, &runtimeError{node.line, "TypeError", fmt.Sprintf("array_merge(): Argument #%d must be of type array, %s given", i+1, arg.value.diagnosticName())}
		}
		for _, entry := range arg.value.a.entries {
			key := entry.key
			if !key.isString {
				key, err = out.a.appendKey(node.line)
				if err != nil {
					return value{}, err
				}
			}
			out.a = out.a.with(key, dereference(entry.value))
		}
	}
	return out, nil
}

func (e *evaluator) callArrayValues(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_values() expects exactly 1 argument"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_values", 1, "array", "array", bound[0])
	}
	out := arrayValue()
	for _, entry := range bound[0].a.entries {
		key, keyErr := out.a.appendKey(node.line)
		if keyErr != nil {
			return value{}, keyErr
		}
		out.a = out.a.with(key, dereference(entry.value))
	}
	return out, nil
}

func (e *evaluator) callArrayShift(node functionCall) (value, error) {
	if len(node.args) != 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_shift() expects exactly 1 argument"}
	}
	target := node.args[0]
	if named, ok := target.(namedArgument); ok {
		if named.name != "array" {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + named.name}
		}
		target = named.operand
	}
	if !writable(target) {
		return value{}, &runtimeError{node.line, "Error", "array_shift(): Argument #1 ($array) cannot be passed by reference"}
	}
	loc, err := e.locate(target)
	if err != nil {
		return value{}, err
	}
	array := loc.get()
	if array.kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_shift", 1, "array", "array", array)
	}
	if len(array.a.entries) == 0 {
		return value{kind: nullKind}, nil
	}
	first := dereference(array.a.entries[0].value)
	out := arrayValue()
	for _, entry := range array.a.entries[1:] {
		key := entry.key
		if !key.isString {
			key, err = out.a.appendKey(node.line)
			if err != nil {
				return value{}, err
			}
		}
		out.a = out.a.with(key, entry.value)
	}
	if err := loc.set(out); err != nil {
		return value{}, err
	}
	return first, nil
}

func (e *evaluator) callArrayPop(node functionCall) (value, error) {
	if len(node.args) != 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_pop() expects exactly 1 argument"}
	}
	target := node.args[0]
	if named, ok := target.(namedArgument); ok {
		if named.name != "array" {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + named.name}
		}
		target = named.operand
	}
	if !writable(target) {
		return value{}, &runtimeError{node.line, "Error", "array_pop(): Argument #1 ($array) cannot be passed by reference"}
	}
	loc, err := e.locate(target)
	if err != nil {
		return value{}, err
	}
	array := loc.get()
	if array.kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_pop", 1, "array", "array", array)
	}
	if len(array.a.entries) == 0 {
		return value{kind: nullKind}, nil
	}
	last := array.a.entries[len(array.a.entries)-1]
	copyArray := *array.a
	copyArray.entries = append([]arrayEntry(nil), array.a.entries[:len(array.a.entries)-1]...)
	if !last.key.isString && last.key.index >= 0 {
		copyArray.next = last.key.index
		copyArray.full = false
	}
	if err := loc.set(value{kind: arrayKind, a: &copyArray}); err != nil {
		return value{}, err
	}
	return dereference(last.value), nil
}

func (e *evaluator) callInArray(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "needle"}, {name: "haystack"}, {name: "strict"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "in_array() expects 2 to 3 arguments"}
	}
	if bound[1].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "in_array", 2, "haystack", "array", bound[1])
	}
	strict := present[2] && bound[2].truthy()
	for _, entry := range bound[1].a.entries {
		candidate := dereference(entry.value)
		if strict {
			if bound[0].identical(candidate) {
				return boolean(true), nil
			}
		} else {
			result, cmpErr := e.comparisonValue("==", bound[0], candidate, node.line)
			if cmpErr != nil {
				return value{}, cmpErr
			}
			if result.truthy() {
				return boolean(true), nil
			}
		}
	}
	return boolean(false), nil
}

func (e *evaluator) callArrayEnd(node functionCall) (value, error) {
	if len(node.args) != 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "end() expects exactly 1 argument"}
	}
	target := node.args[0]
	if named, ok := target.(namedArgument); ok {
		if named.name != "array" {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + named.name}
		}
		target = named.operand
	}
	if !writable(target) {
		return value{}, &runtimeError{node.line, "Error", "end(): Argument #1 ($array) cannot be passed by reference"}
	}
	loc, err := e.locate(target)
	if err != nil {
		return value{}, err
	}
	array := loc.get()
	if array.kind != arrayKind {
		return value{}, builtinTypeError(node.line, "end", 1, "array", "array|object", array)
	}
	if len(array.a.entries) == 0 {
		return boolean(false), nil
	}
	return dereference(array.a.entries[len(array.a.entries)-1].value), nil
}

func (e *evaluator) callArrayCurrent(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "current() expects exactly 1 argument"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "current", 1, "array", "array|object", bound[0])
	}
	if len(bound[0].a.entries) == 0 {
		return boolean(false), nil
	}
	return dereference(bound[0].a.entries[0].value), nil
}

func (e *evaluator) callArraySlice(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "offset"}, {name: "length"}, {name: "preserve_keys"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 4 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_slice() expects 2 to 4 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_slice", 1, "array", "array", bound[0])
	}
	offset, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "array_slice", 2, "offset", "int", bound[1])
	}
	n := int64(len(bound[0].a.entries))
	start := offset.i
	if start < 0 {
		start += n
	}
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}
	end := n
	if present[2] && bound[2].kind != nullKind {
		length, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "array_slice", 3, "length", "?int", bound[2])
		}
		if length.i >= 0 {
			end = start + length.i
			if end > n {
				end = n
			}
		} else {
			end = n + length.i
			if end < start {
				end = start
			}
		}
	}
	preserve := present[3] && bound[3].truthy()
	out := arrayValue()
	for _, entry := range bound[0].a.entries[start:end] {
		key := entry.key
		if !preserve && !key.isString {
			key, err = out.a.appendKey(node.line)
			if err != nil {
				return value{}, err
			}
		}
		out.a = out.a.with(key, dereference(entry.value))
	}
	return out, nil
}

func (e *evaluator) callArrayKeyExists(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "key"}, {name: "array"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 2 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_key_exists() expects exactly 2 arguments"}
	}
	if bound[1].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_key_exists", 2, "array", "array", bound[1])
	}
	key, err := e.arrayKey(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	return boolean(bound[1].a.find(key) >= 0), nil
}
