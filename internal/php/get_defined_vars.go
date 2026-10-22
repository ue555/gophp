package php

import (
	"fmt"
	"sort"
)

func (e *evaluator) callGetDefinedVars(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if _, _, _, err := bindCallArguments(args, nil, node.line); err != nil {
		return value{}, err
	}
	if len(args) != 0 {
		e.captureFailure(callFrame{method: "get_defined_vars", line: node.line, filename: e.filename, args: argumentValues(args)})
		return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("get_defined_vars() expects exactly 0 arguments, %d given", len(args))}
	}
	names := make(map[string]bool)
	for name := range e.variables {
		if name != "this" {
			names[name] = true
		}
	}
	for name := range e.bindings {
		names[name] = true
	}
	keys := make([]string, 0, len(names))
	for name := range names {
		keys = append(keys, name)
	}
	// The interpreter's variable frames are maps (as for $GLOBALS). Preserve
	// values/reference cells and scope isolation, with deterministic key order.
	sort.Strings(keys)
	out := arrayValue()
	for _, name := range keys {
		if e.bindings[name] != 0 {
			e.ensureVariableReference(name)
		}
		if v, exists := e.variableOwner(name)[name]; exists {
			out.a = out.a.with(arrayKey{text: name, isString: true}, v)
		}
	}
	return out, nil
}
