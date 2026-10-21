package php

import "strings"

func (e *evaluator) callStrReplace(node functionCall) (value, error) {
	if len(node.args) < 3 || len(node.args) > 4 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "str_replace() expects 3 to 4 arguments"}
	}
	args, err := e.evalArguments(node.args[:3])
	if err != nil {
		return value{}, err
	}
	search := []value{args[0].value}
	if args[0].value.kind == arrayKind {
		search = nil
		for _, entry := range args[0].value.a.entries {
			search = append(search, dereference(entry.value))
		}
	}
	replace := []value{args[1].value}
	if args[1].value.kind == arrayKind {
		replace = nil
		for _, entry := range args[1].value.a.entries {
			replace = append(replace, dereference(entry.value))
		}
	}
	count := int64(0)
	replaceText := func(input value) (value, error) {
		text, err := e.toText(input, node.line)
		if err != nil {
			return value{}, err
		}
		for i, item := range search {
			old, err := e.toText(item, node.line)
			if err != nil {
				return value{}, err
			}
			newValue := value{kind: stringKind}
			if i < len(replace) {
				newValue = replace[i]
			}
			newText, err := e.toText(newValue, node.line)
			if err != nil {
				return value{}, err
			}
			if old != "" {
				count += int64(strings.Count(text, old))
				text = strings.ReplaceAll(text, old, newText)
			}
		}
		return value{kind: stringKind, s: text}, nil
	}
	var result value
	if args[2].value.kind == arrayKind {
		result = arrayValue()
		for _, entry := range args[2].value.a.entries {
			converted, err := replaceText(dereference(entry.value))
			if err != nil {
				return value{}, err
			}
			result.a = result.a.with(entry.key, converted)
		}
	} else {
		result, err = replaceText(args[2].value)
		if err != nil {
			return value{}, err
		}
	}
	if len(node.args) == 4 {
		if !writable(node.args[3]) {
			return value{}, &runtimeError{node.line, "Error", "str_replace(): Argument #4 ($count) could not be passed by reference"}
		}
		_, err = e.evalAssignment(assignment{target: node.args[3], source: literal{value: integer(count), line: node.line}, line: node.line})
		if err != nil {
			return value{}, err
		}
	}
	return result, nil
}
