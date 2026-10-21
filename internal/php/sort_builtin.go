package php

import (
	"sort"
	"strings"
)

func (e *evaluator) callSort(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	if len(node.args) < 1 || len(node.args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects 1 to 2 arguments"}
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
	if len(node.args) == 2 {
		if _, err := e.eval(node.args[1]); err != nil {
			return value{}, err
		}
	}
	items := append([]arrayEntry(nil), input.a.entries...)
	sort.SliceStable(items, func(i, j int) bool {
		a, b := dereference(items[i].value), dereference(items[j].value)
		if a.kind == intKind && b.kind == intKind {
			if name == "arsort" {
				return a.i > b.i
			}
			return a.i < b.i
		}
		if name == "arsort" {
			return a.text() > b.text()
		}
		return a.text() < b.text()
	})
	out := arrayValue()
	for i, item := range items {
		key := arrayKey{index: int64(i)}
		if name != "sort" {
			key = item.key
		}
		out.a = out.a.with(key, item.value)
	}
	if _, err := e.evalAssignment(assignment{target: node.args[0], source: literal{value: out, line: node.line}, line: node.line}); err != nil {
		return value{}, err
	}
	return boolean(true), nil
}
