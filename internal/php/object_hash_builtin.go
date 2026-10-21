package php

import (
	"fmt"
	"strings"
)

func (e *evaluator) callObjectHash(node functionCall) (value, error) {
	function := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "object"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", function + "() expects exactly 1 argument"}
	}
	if bound[0].kind != objectKind {
		return value{}, builtinTypeError(node.line, function, 1, "object", "object", bound[0])
	}
	id := e.objectIDs[bound[0].o]
	if id == 0 {
		id = int64(len(e.objectIDs) + 1)
		e.objectIDs[bound[0].o] = id
	}
	if strings.EqualFold(node.name, "spl_object_id") {
		return integer(id), nil
	}
	return value{kind: stringKind, s: fmt.Sprintf("%016x%016x", id, 0)}, nil
}
