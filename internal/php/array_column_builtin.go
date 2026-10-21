package php

func (e *evaluator) callArrayColumn(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "column_key"}, {name: "index_key"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_column() expects 2 to 3 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_column", 1, "array", "array", bound[0])
	}
	out := arrayValue()
	for _, entry := range bound[0].a.entries {
		row := dereference(entry.value)
		var source *phpArray
		if row.kind == arrayKind {
			source = row.a
		} else if row.kind == objectKind && row.o != nil {
			source = row.o.properties
		} else {
			continue
		}
		if source == nil {
			continue
		}
		item := row
		if bound[1].kind != nullKind {
			key, err := e.arrayKey(bound[1], node.line)
			if err != nil {
				return value{}, err
			}
			found, ok := source.get(key)
			if !ok {
				continue
			}
			item = found
		}
		var key arrayKey
		if present[2] && bound[2].kind != nullKind {
			lookup, err := e.arrayKey(bound[2], node.line)
			if err != nil {
				return value{}, err
			}
			index, ok := source.get(lookup)
			if ok {
				key, err = e.arrayKey(index, node.line)
				if err != nil {
					return value{}, err
				}
			} else {
				key, err = out.a.appendKey(node.line)
				if err != nil {
					return value{}, err
				}
			}
		} else {
			key, err = out.a.appendKey(node.line)
			if err != nil {
				return value{}, err
			}
		}
		out.a = out.a.with(key, item)
	}
	return out, nil
}
