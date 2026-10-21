package php

import (
	"fmt"
	"strconv"
	"strings"
)

func encodeQueryPart(source string, rfc3986 bool) string {
	var out strings.Builder
	for i := 0; i < len(source); i++ {
		b := source[i]
		if b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b == '_' || b == '.' || b == '~' && rfc3986 {
			out.WriteByte(b)
		} else if b == ' ' && !rfc3986 {
			out.WriteByte('+')
		} else {
			out.WriteString(fmt.Sprintf("%%%02X", b))
		}
	}
	return out.String()
}

func (e *evaluator) callHTTPBuildQuery(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "data"}, {name: "numeric_prefix"}, {name: "arg_separator"}, {name: "encoding_type"}}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 4 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "http_build_query() expects 1 to 4 arguments"}
	}
	if bound[0].kind != arrayKind && bound[0].kind != objectKind {
		return value{}, builtinTypeError(node.line, "http_build_query", 1, "data", "array|object", bound[0])
	}
	prefix := ""
	if present[1] && bound[1].kind != nullKind {
		v, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "http_build_query", 2, "numeric_prefix", "string", bound[1])
		}
		prefix = v.s
	}
	separator := "&"
	if present[2] && bound[2].kind != nullKind {
		v, _, ok := e.coerceType(bound[2], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "http_build_query", 3, "arg_separator", "?string", bound[2])
		}
		separator = v.s
	}
	rfc3986 := false
	if present[3] {
		v, _, ok := e.coerceType(bound[3], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "http_build_query", 4, "encoding_type", "int", bound[3])
		}
		if v.i != 1 && v.i != 2 {
			return value{}, &runtimeError{node.line, "ValueError", "http_build_query(): Argument #4 ($encoding_type) must be either PHP_QUERY_RFC1738 or PHP_QUERY_RFC3986"}
		}
		rfc3986 = v.i == 2
	}
	var items []string
	visited := map[*phpArray]bool{}
	var walk func(value, string, int)
	walk = func(v value, key string, depth int) {
		if depth > 64 {
			return
		}
		v = dereference(v)
		if v.kind == nullKind {
			return
		}
		var entries []arrayEntry
		if v.kind == arrayKind && v.a != nil {
			if visited[v.a] {
				return
			}
			visited[v.a] = true
			defer delete(visited, v.a)
			entries = v.a.entries
		}
		if v.kind == objectKind && v.o != nil && v.o.properties != nil {
			entries = v.o.properties.entries
		}
		if v.kind == arrayKind || v.kind == objectKind {
			for _, entry := range entries {
				part := entry.key.text
				if !entry.key.isString {
					part = strconv.FormatInt(entry.key.index, 10)
					if key == "" {
						part = prefix + part
					}
				}
				child := part
				if key != "" {
					child = key + "[" + part + "]"
				}
				walk(entry.value, child, depth+1)
			}
			return
		}
		if v.kind == boolKind && v.i == 0 {
			items = append(items, encodeQueryPart(key, rfc3986)+"=0")
			return
		}
		items = append(items, encodeQueryPart(key, rfc3986)+"="+encodeQueryPart(v.text(), rfc3986))
	}
	walk(bound[0], "", 0)
	return value{kind: stringKind, s: strings.Join(items, separator)}, nil
}
