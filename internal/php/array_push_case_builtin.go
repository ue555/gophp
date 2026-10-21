package php

import "strings"

func (e *evaluator) callArrayPush(node functionCall) (value, error) {
	if len(node.args) < 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_push() expects at least 1 argument"}
	}
	target := node.args[0]
	if named, ok := target.(namedArgument); ok {
		if named.name != "array" {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + named.name}
		}
		target = named.operand
	}
	if !writable(target) {
		return value{}, &runtimeError{node.line, "Error", "array_push(): Argument #1 ($array) cannot be passed by reference"}
	}
	loc, err := e.locate(target)
	if err != nil {
		return value{}, err
	}
	current := loc.get()
	if current.kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_push", 1, "array", "array", current)
	}
	items, err := e.evalArguments(node.args[1:])
	if err != nil {
		return value{}, err
	}
	out := current
	for _, item := range items {
		if item.named {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + item.name}
		}
		key, err := out.a.appendKey(node.line)
		if err != nil {
			return value{}, err
		}
		out.a = e.arrayWith(out.a, key, item.value)
	}
	if err := loc.set(out); err != nil {
		return value{}, err
	}
	return integer(int64(len(out.a.entries))), nil
}

func (e *evaluator) callArrayChangeKeyCase(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "case"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_change_key_case() expects 1 to 2 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_change_key_case", 1, "array", "array", bound[0])
	}
	mode := int64(0)
	if present[1] {
		mode = bound[1].i
	}
	if mode != 0 && mode != 1 {
		return value{}, &runtimeError{node.line, "ValueError", "array_change_key_case(): Argument #2 ($case) must be either CASE_UPPER or CASE_LOWER"}
	}
	out := arrayValue()
	for _, entry := range bound[0].a.entries {
		key := entry.key
		if key.isString {
			if mode == 1 {
				key.text = strings.ToUpper(key.text)
			} else {
				key.text = strings.ToLower(key.text)
			}
		}
		out.a = out.a.with(key, entry.value)
	}
	return out, nil
}
