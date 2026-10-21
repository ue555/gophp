package php

import (
	"fmt"
	"strings"
)

func (e *evaluator) callImplode(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("implode() expects 1 to 2 arguments, %d given", len(args))}
	}
	separator := ""
	var elements value
	if len(args) == 1 {
		if args[0].named && args[0].name != "array" {
			return value{}, &runtimeError{node.line, "Error", "Unknown named parameter $" + args[0].name}
		}
		elements = args[0].value
	} else {
		bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "separator"}, {name: "array"}}, node.line)
		if bindErr != nil {
			return value{}, bindErr
		}
		if !present[0] || !present[1] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "implode() expects 2 arguments"}
		}
		sep, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, &runtimeError{node.line, "TypeError", "implode(): Argument #1 ($separator) must be of type string, " + bound[0].diagnosticName() + " given"}
		}
		separator, elements = sep.s, bound[1]
	}
	if elements.kind != arrayKind {
		return value{}, &runtimeError{node.line, "TypeError", fmt.Sprintf("implode(): Argument #%d ($array) must be of type array, %s given", len(args), elements.diagnosticName())}
	}
	parts := make([]string, 0, len(elements.a.entries))
	for _, entry := range elements.a.entries {
		part, convertErr := e.toText(dereference(entry.value), node.line)
		if convertErr != nil {
			return value{}, convertErr
		}
		parts = append(parts, part)
	}
	return value{kind: stringKind, s: strings.Join(parts, separator)}, nil
}
