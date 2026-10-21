package php

import (
	"bytes"
	"encoding/json"
	"strings"
)

type orderedJSONObject []struct {
	key   string
	value any
}

func (entries orderedJSONObject) MarshalJSON() ([]byte, error) {
	var out bytes.Buffer
	out.WriteByte('{')
	for index, entry := range entries {
		if index > 0 {
			out.WriteByte(',')
		}
		key, err := json.Marshal(entry.key)
		if err != nil {
			return nil, err
		}
		item, err := json.Marshal(entry.value)
		if err != nil {
			return nil, err
		}
		out.Write(key)
		out.WriteByte(':')
		out.Write(item)
	}
	out.WriteByte('}')
	return out.Bytes(), nil
}

func (e *evaluator) jsonValue(v value, seen map[*phpArray]bool, depth int, line int) (any, error) {
	if depth > 512 {
		return nil, &runtimeError{line, "JsonException", "Maximum stack depth exceeded"}
	}
	v = dereference(v)
	switch v.kind {
	case nullKind:
		return nil, nil
	case boolKind:
		return v.i != 0, nil
	case intKind:
		return v.i, nil
	case floatKind:
		return v.f, nil
	case stringKind:
		return v.s, nil
	case arrayKind:
		if v.a == nil {
			return []any{}, nil
		}
		if seen[v.a] {
			return nil, &runtimeError{line, "JsonException", "Recursion detected"}
		}
		seen[v.a] = true
		defer delete(seen, v.a)
		list := true
		for i, entry := range v.a.entries {
			if entry.key.isString || entry.key.index != int64(i) {
				list = false
				break
			}
		}
		if list {
			out := make([]any, 0, len(v.a.entries))
			for _, entry := range v.a.entries {
				item, err := e.jsonValue(entry.value, seen, depth+1, line)
				if err != nil {
					return nil, err
				}
				out = append(out, item)
			}
			return out, nil
		}
		out := make(orderedJSONObject, 0, len(v.a.entries))
		for _, entry := range v.a.entries {
			key := entry.key.text
			if !entry.key.isString {
				key = integer(entry.key.index).text()
			}
			item, err := e.jsonValue(entry.value, seen, depth+1, line)
			if err != nil {
				return nil, err
			}
			out = append(out, struct {
				key   string
				value any
			}{key, item})
		}
		return out, nil
	case objectKind:
		if v.o == nil {
			return nil, nil
		}
		if v.o.enumCase != "" {
			if v.o.class.enumBacking == "" {
				return nil, &runtimeError{line, "JsonException", "Non-backed enums have no default serialization"}
			}
			if item, ok := v.o.properties.get(arrayKey{isString: true, text: "value"}); ok {
				return e.jsonValue(item, seen, depth+1, line)
			}
		}
		if method := e.lookupMethod(v.o.class, "jsonSerialize"); method != nil && classImplements(v.o.class, "JsonSerializable") {
			result, err := e.invoke(v, method, nil, line)
			if err != nil {
				return nil, err
			}
			return e.jsonValue(result, seen, depth+1, line)
		}
		out := make(orderedJSONObject, 0)
		if v.o.properties != nil {
			for _, entry := range v.o.properties.entries {
				if !entry.key.isString || strings.HasPrefix(entry.key.text, "\x00") {
					continue
				}
				item, err := e.jsonValue(entry.value, seen, depth+1, line)
				if err != nil {
					return nil, err
				}
				out = append(out, struct {
					key   string
					value any
				}{entry.key.text, item})
			}
		}
		return out, nil
	default:
		return nil, &runtimeError{line, "JsonException", "Type is not supported"}
	}
}

func (e *evaluator) callJSONEncode(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}, {name: "flags"}, {name: "depth"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 3 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "json_encode() expects 1 to 3 arguments"}
	}
	flags := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "json_encode", 2, "flags", "int", bound[1])
		}
		flags = v.i
	}
	converted, err := e.jsonValue(bound[0], make(map[*phpArray]bool), 0, node.line)
	if err != nil {
		if flags&4194304 != 0 {
			return value{}, err
		}
		return boolean(false), nil
	}
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetEscapeHTML(false)
	if flags&128 != 0 {
		encoder.SetIndent("", "    ")
	}
	if err := encoder.Encode(converted); err != nil {
		if flags&4194304 != 0 {
			return value{}, &runtimeError{node.line, "JsonException", err.Error()}
		}
		return boolean(false), nil
	}
	out := strings.TrimSuffix(b.String(), "\n")
	if flags&64 == 0 {
		out = strings.ReplaceAll(out, "/", "\\/")
	}
	return value{kind: stringKind, s: out}, nil
}
