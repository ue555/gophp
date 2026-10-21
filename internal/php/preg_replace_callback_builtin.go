package php

import "strings"

func (e *evaluator) callPregReplaceCallback(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "pattern"}, {name: "callback"}, {name: "subject"}, {name: "limit"}, {name: "count"}, {name: "flags"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 3 || len(args) > 6 || !present[0] || !present[1] || !present[2] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "preg_replace_callback() expects between 3 and 6 arguments"}
	}
	pattern, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "preg_replace_callback", 1, "pattern", "string|array", bound[0])
	}
	subject, _, ok := e.coerceType(bound[2], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "preg_replace_callback", 3, "subject", "string|array", bound[2])
	}
	target, ok := e.resolveCallable(bound[1], node.line)
	if !ok {
		return value{}, &runtimeError{node.line, "TypeError", "preg_replace_callback(): Argument #2 ($callback) must be a valid callback"}
	}
	limit := -1
	if present[3] {
		v, _, ok := e.coerceType(bound[3], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_replace_callback", 4, "limit", "int", bound[3])
		}
		limit = int(v.i)
	}
	e.pregLastError = 0
	indices, _, err := pcre2Matches(pattern.s, subject.s, 0, limit)
	if err != nil {
		return e.pregFailure("preg_replace_callback", node.line, err)
	}

	var out strings.Builder
	last := 0
	for _, index := range indices {
		out.WriteString(subject.s[last:index[0]])
		matches := arrayValue()
		for group := 0; group < len(index)/2; group++ {
			start, end := index[2*group], index[2*group+1]
			part := ""
			if start >= 0 {
				part = subject.s[start:end]
			}
			matches.a = matches.a.with(arrayKey{index: int64(group)}, value{kind: stringKind, s: part})
		}
		result, err := e.invokeInternalCallback(target, []callArgument{{value: matches}}, node.line)
		if err != nil {
			return value{}, err
		}
		replacement, err := e.toText(result, node.line)
		if err != nil {
			return value{}, err
		}
		out.WriteString(replacement)
		last = index[1]
	}
	out.WriteString(subject.s[last:])
	return value{kind: stringKind, s: out.String()}, nil
}

func (e *evaluator) callLcfirst(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "lcfirst() expects exactly 1 argument"}
	}
	s, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "lcfirst", 1, "string", "string", bound[0])
	}
	if len(s.s) > 0 && s.s[0] >= 'A' && s.s[0] <= 'Z' {
		s.s = string(s.s[0]+'a'-'A') + s.s[1:]
	}
	return s, nil
}
