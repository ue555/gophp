package php

import (
	"strconv"
	"strings"
)

func queryArrayKey(name string, array *phpArray) arrayKey {
	if name == "" {
		key, _ := array.appendKey(0)
		return key
	}
	if i, err := strconv.ParseInt(name, 10, 64); err == nil && strconv.FormatInt(i, 10) == name {
		return arrayKey{index: i}
	}
	return arrayKey{isString: true, text: name}
}

func insertQueryValue(array *phpArray, names []string, item value) *phpArray {
	key := queryArrayKey(names[0], array)
	if len(names) == 1 {
		return array.with(key, item)
	}
	nested, ok := array.get(key)
	if !ok || nested.kind != arrayKind {
		nested = arrayValue()
	}
	nested.a = insertQueryValue(nested.a, names[1:], item)
	return array.with(key, nested)
}

func queryKeyParts(raw string) []string {
	open := strings.IndexByte(raw, '[')
	if open < 0 {
		return []string{strings.NewReplacer(".", "_", " ", "_").Replace(raw)}
	}
	parts := []string{strings.NewReplacer(".", "_", " ", "_").Replace(raw[:open])}
	for open < len(raw) && raw[open] == '[' {
		close := strings.IndexByte(raw[open+1:], ']')
		if close < 0 {
			break
		}
		parts = append(parts, raw[open+1:open+1+close])
		open += close + 2
	}
	return parts
}

func parseQueryArray(input string) value {
	out := arrayValue()
	if input == "" {
		return out
	}
	for _, pair := range strings.Split(input, "&") {
		if pair == "" {
			continue
		}
		key, text, found := strings.Cut(pair, "=")
		if !found {
			text = ""
		}
		key = decodeURLString(key, true)
		if key == "" {
			continue
		}
		parts := queryKeyParts(key)
		out.a = insertQueryValue(out.a, parts, value{kind: stringKind, s: decodeURLString(text, true)})
	}
	return out
}

func (e *evaluator) callParseStr(node functionCall) (value, error) {
	if len(node.args) != 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "parse_str() expects exactly 2 arguments"}
	}
	input, err := e.eval(node.args[0])
	if err != nil {
		return value{}, err
	}
	v, _, ok := e.coerceType(input, "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "parse_str", 1, "string", "string", input)
	}
	if !writable(node.args[1]) {
		return value{}, &runtimeError{node.line, "Error", "parse_str(): Argument #2 ($result) could not be passed by reference"}
	}
	out := parseQueryArray(v.s)
	if _, err := e.evalAssignment(assignment{target: node.args[1], source: literal{value: out, line: node.line}, line: node.line}); err != nil {
		return value{}, err
	}
	return value{kind: nullKind}, nil
}
