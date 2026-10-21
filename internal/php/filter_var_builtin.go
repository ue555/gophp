package php

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

func (e *evaluator) callFilterVar(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}, {name: "filter"}, {name: "options"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 3 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "filter_var() expects 1 to 3 arguments"}
	}
	filter := int64(516) // FILTER_DEFAULT / FILTER_UNSAFE_RAW
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "filter_var", 2, "filter", "int", bound[1])
		}
		filter = v.i
	}
	flags := int64(0)
	if present[2] {
		if bound[2].kind == arrayKind {
			if entry, found := bound[2].a.getRaw(arrayKey{isString: true, text: "flags"}); found {
				v, _, ok := e.coerceType(entry, "int", false, e.strictTypes)
				if !ok {
					return value{}, builtinTypeError(node.line, "filter_var", 3, "options", "array|int", bound[2])
				}
				flags = v.i
			}
		} else if bound[2].kind == intKind {
			flags = bound[2].i
		} else {
			return value{}, builtinTypeError(node.line, "filter_var", 3, "options", "array|int", bound[2])
		}
	}
	if filter == 516 {
		return bound[0], nil
	}
	if filter == 257 || filter == 258 {
		s, _, ok := e.coerceType(bound[0], "string", false, false)
		if ok {
			if filter == 257 {
				trimmed := strings.TrimSpace(s.s)
				if n, err := strconv.ParseInt(trimmed, 10, 64); err == nil && (trimmed == "0" || !strings.HasPrefix(strings.TrimPrefix(trimmed, "-"), "0")) {
					return integer(n), nil
				}
			} else {
				switch strings.ToLower(strings.TrimSpace(s.s)) {
				case "1", "true", "on", "yes":
					return boolean(true), nil
				case "0", "false", "off", "no", "":
					return boolean(false), nil
				}
			}
		}
		if flags&134217728 != 0 {
			return value{kind: nullKind}, nil
		}
		return boolean(false), nil
	}
	if filter == 273 {
		s, _, ok := e.coerceType(bound[0], "string", false, false)
		if ok {
			u, parseErr := url.Parse(s.s)
			if parseErr == nil && u.Scheme != "" && u.Host != "" && !strings.ContainsAny(s.s, " \t\r\n") {
				return s, nil
			}
		}
		if flags&134217728 != 0 {
			return value{kind: nullKind}, nil
		}
		return boolean(false), nil
	}
	if filter != 275 {
		return value{}, &runtimeError{node.line, "Error", "gophp: filter_var() filter is not implemented"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, false)
	if ok {
		ip := net.ParseIP(s.s)
		if ip != nil && (flags&1048576 == 0 || ip.To4() != nil) && (flags&2097152 == 0 || ip.To4() == nil) {
			return s, nil
		}
	}
	if flags&134217728 != 0 {
		return value{kind: nullKind}, nil
	}
	return boolean(false), nil
}
