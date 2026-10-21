package php

import (
	"net/url"
	"strconv"
	"strings"
)

func (e *evaluator) callParseURL(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "url"}, {name: "component"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "parse_url() expects 1 to 2 arguments"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "parse_url", 1, "url", "string", bound[0])
	}
	component := int64(-1)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "parse_url", 2, "component", "int", bound[1])
		}
		component = v.i
	}
	if component < -1 || component > 7 {
		return value{}, &runtimeError{node.line, "ValueError", "parse_url(): Argument #2 ($component) must be a valid URL component identifier, -1 given"}
	}
	u, parseErr := url.Parse(s.s)
	if parseErr != nil {
		return boolean(false), nil
	}
	fields := make(map[int64]value)
	put := func(id int64, text string) { fields[id] = value{kind: stringKind, s: text} }
	if u.Scheme != "" {
		put(0, u.Scheme)
	}
	if u.Hostname() != "" {
		put(1, u.Hostname())
	}
	if u.Port() != "" {
		p, err := strconv.ParseInt(u.Port(), 10, 64)
		if err != nil || p < 0 || p > 65535 {
			return boolean(false), nil
		}
		fields[2] = integer(p)
	}
	if u.User != nil {
		put(3, u.User.Username())
		if password, has := u.User.Password(); has {
			put(4, password)
		}
	}
	if u.Path != "" {
		put(5, u.EscapedPath())
	}
	if strings.Contains(s.s, "?") {
		put(6, u.RawQuery)
	}
	if strings.Contains(s.s, "#") {
		put(7, u.Fragment)
	}
	if component >= 0 {
		if v, ok := fields[component]; ok {
			return v, nil
		}
		return value{kind: nullKind}, nil
	}
	out := arrayValue()
	for id, name := range []string{"scheme", "host", "port", "user", "pass", "path", "query", "fragment"} {
		if v, ok := fields[int64(id)]; ok {
			out.a = out.a.with(arrayKey{isString: true, text: name}, v)
		}
	}
	return out, nil
}
