package php

import "strings"

func (e *evaluator) callArrayFilterOrUnique(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "array"}, {name: "callback"}, {name: "mode"}}
	if name == "array_unique" {
		params = []parameter{{name: "array"}, {name: "flags"}}
	}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > len(params) || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects 1 to 3 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, name, 1, "array", "array", bound[0])
	}
	out := arrayValue()
	seen := make(map[string]bool)
	flags := int64(2) // array_unique defaults to SORT_STRING.
	if name == "array_unique" && present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, name, 2, "flags", "int", bound[1])
		}
		flags = v.i
		if flags != 0 && flags != 2 {
			return value{}, &runtimeError{node.line, "Error", "gophp: array_unique() supports SORT_REGULAR and SORT_STRING"}
		}
	}
	var callback resolvedCallable
	if name == "array_filter" && present[1] && bound[1].kind != nullKind {
		var ok bool
		callback, ok = e.resolveCallable(bound[1], node.line)
		if !ok {
			return value{}, &runtimeError{node.line, "TypeError", "array_filter(): Argument #2 ($callback) must be a valid callback or null"}
		}
	}
	for _, entry := range bound[0].a.entries {
		item := dereference(entry.value)
		keep := true
		if name == "array_unique" {
			if flags == 0 {
				for _, previous := range out.a.entries {
					order, err := e.valueOrder(item, dereference(previous.value), node.line, make(map[*phpArray]bool))
					if err != nil {
						return value{}, err
					}
					if order == 0 {
						keep = false
						break
					}
				}
			} else {
				key, convertErr := e.toText(item, node.line)
				if convertErr != nil {
					return value{}, convertErr
				}
				keep = !seen[key]
				seen[key] = true
			}
		} else if present[1] && bound[1].kind != nullKind {
			callArgs := []callArgument{{value: item}}
			if present[2] && bound[2].i != 0 {
				key := integer(entry.key.index)
				if entry.key.isString {
					key = value{kind: stringKind, s: entry.key.text}
				}
				if bound[2].i == 2 {
					callArgs = []callArgument{{value: key}}
				} else {
					callArgs = append(callArgs, callArgument{value: key})
				}
			}
			result, callErr := e.invokeInternalCallback(callback, callArgs, node.line)
			if callErr != nil {
				return value{}, callErr
			}
			keep = result.truthy()
		} else {
			keep = item.truthy()
		}
		if keep {
			out.a = out.a.with(entry.key, item)
		}
	}
	return out, nil
}
