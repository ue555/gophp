package php

import (
	"sort"
	"strings"
)

func (e *evaluator) callGetDeclaredTypes(node functionCall) (value, error) {
	if len(node.args) != 0 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects exactly 0 arguments"}
	}
	e.initClasses()
	wantInterface := strings.EqualFold(node.name, "get_declared_interfaces")
	names := make([]string, 0, len(e.classes))
	seen := make(map[string]bool)
	for _, class := range e.classes {
		if class == nil || (class.kind == "interface") != wantInterface {
			continue
		}
		if !seen[class.name] {
			names = append(names, class.name)
			seen[class.name] = true
		}
	}
	// Registration order is not retained in the current class map. The fixed
	// scanner sorts this result before using it; keep this result deterministic.
	sort.Strings(names)
	out := arrayValue()
	for i, name := range names {
		out.a = out.a.with(arrayKey{index: int64(i)}, value{kind: stringKind, s: name})
	}
	return out, nil
}
