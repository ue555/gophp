package php

import (
	"strings"
)

func (e *evaluator) callPregReplace(node functionCall) (value, error) {
	limit := int64(-1)
	if len(node.args) < 3 || len(node.args) > 5 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "preg_replace() expects 3 to 5 arguments"}
	}
	args, err := e.evalArguments(node.args[:3])
	if err != nil {
		return value{}, err
	}
	if len(node.args) >= 4 {
		limitValue, err := e.eval(node.args[3])
		if err != nil {
			return value{}, err
		}
		converted, _, ok := e.coerceType(limitValue, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_replace", 4, "limit", "int", limitValue)
		}
		limit = converted.i
	}
	pattern, err := e.toText(args[0].value, node.line)
	if err != nil {
		return value{}, err
	}
	replacement, err := e.toText(args[1].value, node.line)
	if err != nil {
		return value{}, err
	}
	e.pregLastError = 0
	count := int64(0)
	convert := func(v value) (value, error) {
		subject, err := e.toText(v, node.line)
		if err != nil {
			return value{}, err
		}
		matches, _, err := pcre2Matches(pattern, subject, 0, int(limit))
		if err != nil {
			return value{}, err
		}
		var out strings.Builder
		last := 0
		for _, match := range matches {
			out.WriteString(subject[last:match[0]])
			out.WriteString(expandPCREReplacement(replacement, subject, match))
			last = match[1]
			count++
		}
		out.WriteString(subject[last:])
		return value{kind: stringKind, s: out.String()}, nil
	}

	var result value
	if args[2].value.kind == arrayKind {
		result = arrayValue()
		for _, entry := range args[2].value.a.entries {
			item, err := convert(dereference(entry.value))
			if err != nil {
				return e.pregFailure("preg_replace", node.line, err)
			}
			result.a = result.a.with(entry.key, item)
		}
	} else {
		result, err = convert(args[2].value)
		if err != nil {
			return e.pregFailure("preg_replace", node.line, err)
		}
	}
	if len(node.args) >= 5 {
		if !writable(node.args[4]) {
			return value{}, &runtimeError{node.line, "Error", "preg_replace(): Argument #5 ($count) could not be passed by reference"}
		}
		_, err = e.evalAssignment(assignment{target: node.args[4], source: literal{value: integer(count), line: node.line}, line: node.line})
		if err != nil {
			return value{}, err
		}
	}
	return result, nil
}

func expandPCREReplacement(replacement, subject string, match []int) string {
	var out strings.Builder
	for i := 0; i < len(replacement); i++ {
		if replacement[i] == '$' || replacement[i] == '\\' {
			j := i + 1
			braced := false
			if j < len(replacement) && replacement[j] == '{' && replacement[i] == '$' {
				braced = true
				j++
			}
			start := j
			for j < len(replacement) && replacement[j] >= '0' && replacement[j] <= '9' {
				j++
			}
			if j > start && (!braced || (j < len(replacement) && replacement[j] == '}')) {
				group := 0
				for _, digit := range replacement[start:j] {
					group = group*10 + int(digit-'0')
				}
				if group*2+1 < len(match) && match[group*2] >= 0 {
					out.WriteString(subject[match[group*2]:match[group*2+1]])
				}
				if braced {
					j++
				}
				i = j - 1
				continue
			}
		}
		out.WriteByte(replacement[i])
	}
	return out.String()
}
