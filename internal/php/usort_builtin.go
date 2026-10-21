package php

import "strings"

func (e *evaluator) callUsort(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	if len(node.args) != 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects exactly 2 arguments"}
	}
	if !writable(node.args[0]) {
		return value{}, &runtimeError{node.line, "Error", name + "(): Argument #1 ($array) cannot be passed by reference"}
	}
	loc, err := e.locate(node.args[0])
	if err != nil {
		return value{}, err
	}
	input := loc.get()
	if input.kind != arrayKind {
		return value{}, builtinTypeError(node.line, name, 1, "array", "array", input)
	}
	callback, err := e.eval(node.args[1])
	if err != nil {
		return value{}, err
	}
	target, ok := e.resolveCallable(callback, node.line)
	if !ok {
		return value{}, &runtimeError{node.line, "TypeError", name + "(): Argument #2 ($callback) must be a valid callback"}
	}
	items := append([]arrayEntry(nil), input.a.entries...)
	// Use an error-aware stable insertion sort. A callback exception must stop
	// immediately, before the final array assignment.
	for i := 1; i < len(items); i++ {
		for j := i; j > 0; j-- {
			left, right := dereference(items[j-1].value), dereference(items[j].value)
			if name == "uksort" {
				toValue := func(k arrayKey) value {
					if k.isString {
						return value{kind: stringKind, s: k.text}
					}
					return integer(k.index)
				}
				left, right = toValue(items[j-1].key), toValue(items[j].key)
			}
			comparison, err := e.invokeInternalCallback(target, []callArgument{{value: left}, {value: right}}, node.line)
			if err != nil {
				return value{}, err
			}
			if comparison.kind != intKind {
				converted, _, valid := e.coerceType(comparison, "int", false, false)
				if !valid {
					return value{}, builtinTypeError(node.line, name, 2, "callback", "int", comparison)
				}
				comparison = converted
			}
			if comparison.i <= 0 {
				break
			}
			items[j-1], items[j] = items[j], items[j-1]
		}
	}
	out := arrayValue()
	for i, item := range items {
		key := arrayKey{index: int64(i)}
		if name == "uasort" || name == "uksort" {
			key = item.key
		}
		out.a = out.a.with(key, item.value)
	}
	if err := loc.set(out); err != nil {
		return value{}, err
	}
	return boolean(true), nil
}
