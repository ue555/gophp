package php

import (
	"sort"
	"strings"
)

func (e *evaluator) callSubstr(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "offset"}, {name: "length"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || !present[1] || len(args) > 3 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "substr() expects between 2 and 3 arguments"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "substr", 1, "string", "string", bound[0])
	}
	source, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	offset, err := e.countMode(bound[1], node.line)
	if err != nil {
		return value{}, err
	}
	start := int(offset)
	if offset < 0 {
		start = len(source) + int(offset)
	}
	if start < 0 {
		start = 0
	}
	if start > len(source) {
		start = len(source)
	}
	end := len(source)
	if present[2] && bound[2].kind != nullKind {
		length, err := e.countMode(bound[2], node.line)
		if err != nil {
			return value{}, err
		}
		if length >= 0 {
			end = start + int(length)
		} else {
			end = len(source) + int(length)
		}
		if end < start {
			end = start
		}
		if end > len(source) {
			end = len(source)
		}
	}
	return value{kind: stringKind, s: source[start:end]}, nil
}

func (e *evaluator) callStrrpos(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "haystack"}, {name: "needle"}, {name: "offset"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || !present[1] || len(args) > 3 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strrpos() expects between 2 and 3 arguments"}
	}
	for i := 0; i < 2; i++ {
		if bound[i].kind == arrayKind || bound[i].kind == objectKind {
			return value{}, builtinTypeError(node.line, "strrpos", i+1, []string{"haystack", "needle"}[i], "string", bound[i])
		}
	}
	haystack, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	needle, err := e.toText(bound[1], node.line)
	if err != nil {
		return value{}, err
	}
	start, stop := 0, len(haystack)
	if present[2] {
		offset, err := e.countMode(bound[2], node.line)
		if err != nil {
			return value{}, err
		}
		if offset >= 0 {
			start = int(offset)
		} else {
			stop = len(haystack) + int(offset) + len(needle)
		}
		if start > len(haystack) || stop < 0 {
			return value{}, &runtimeError{node.line, "ValueError", "strrpos(): Argument #3 ($offset) must be contained in argument #1 ($haystack)"}
		}
		if stop > len(haystack) {
			stop = len(haystack)
		}
	}
	index := strings.LastIndex(haystack[start:stop], needle)
	if index < 0 {
		return boolean(false), nil
	}
	return integer(int64(start + index)), nil
}

func (e *evaluator) callStrtr(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "from"}, {name: "to"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || !present[1] || len(args) < 2 || len(args) > 3 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "strtr() expects 2 or 3 arguments"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "strtr", 1, "string", "string", bound[0])
	}
	source, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	if !present[2] {
		if bound[1].kind != arrayKind {
			return value{}, builtinTypeError(node.line, "strtr", 2, "replace_pairs", "array", bound[1])
		}
		pairs := make(map[string]string)
		keys := make([]string, 0, len(bound[1].a.entries))
		for _, item := range bound[1].a.entries {
			key := item.key.text
			if !item.key.isString {
				key = integer(item.key.index).text()
			}
			if key == "" {
				return value{}, &runtimeError{node.line, "ValueError", "strtr(): Argument #2 ($replace_pairs) must not contain any empty strings"}
			}
			replacement, err := e.toText(item.value, node.line)
			if err != nil {
				return value{}, err
			}
			pairs[key] = replacement
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })
		var output strings.Builder
		for i := 0; i < len(source); {
			matched := false
			for _, key := range keys {
				if strings.HasPrefix(source[i:], key) {
					output.WriteString(pairs[key])
					i += len(key)
					matched = true
					break
				}
			}
			if !matched {
				output.WriteByte(source[i])
				i++
			}
		}
		return value{kind: stringKind, s: output.String()}, nil
	}
	from, err := e.toText(bound[1], node.line)
	if err != nil {
		return value{}, err
	}
	to, err := e.toText(bound[2], node.line)
	if err != nil {
		return value{}, err
	}
	mapping := make(map[byte]byte)
	for i := 0; i < len(from) && i < len(to); i++ {
		mapping[from[i]] = to[i]
	}
	result := []byte(source)
	for i, c := range result {
		if mapped, ok := mapping[c]; ok {
			result[i] = mapped
		}
	}
	return value{kind: stringKind, s: string(result)}, nil
}
