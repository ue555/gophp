package php

import (
	"regexp"
	"strings"
	"unicode/utf8"
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
	if len(pattern) < 2 {
		return value{}, &runtimeError{node.line, "Error", "gophp: invalid regular expression"}
	}
	end := strings.LastIndexByte(pattern, pattern[0])
	if end <= 0 {
		return value{}, &runtimeError{node.line, "Error", "gophp: invalid regular expression"}
	}
	body := pattern[1:end]
	// A one-character negative lookbehind before an escaped r/n sequence is
	// supported without changing the preceding character in the result.
	newlineEscape := strings.Contains(body, "(?<!") && strings.HasSuffix(body, "[rn]")
	camelBoundary := pattern == "/(.)(?=[A-Z])/u" && strings.HasPrefix(replacement, "$1")
	tableBoundary := pattern == `~(?<=\w)([A-Z])~u` && replacement == "_$1"
	if camelBoundary || tableBoundary {
		// The lookahead is zero-width. Preserve every uppercase boundary,
		// including adjacent capitals, without consuming the next character.
	} else if newlineEscape {
		body = `(^|[^\\])\\[rn]`
		replacement = "${1}" + replacement
	}
	var re *regexp.Regexp
	if camelBoundary || tableBoundary {
		re = nil
	} else if newlineEscape {
		re, err = regexp.Compile(body)
	} else {
		re, err = compileBasicPCRE(pattern)
	}
	usePCRE2 := err != nil && !newlineEscape
	if err != nil && !usePCRE2 {
		return value{}, &runtimeError{node.line, "Error", "gophp: unsupported regular expression: " + err.Error()}
	}
	count := int64(0)
	convert := func(v value) (value, error) {
		subject, err := e.toText(v, node.line)
		if err != nil {
			return value{}, err
		}
		if limit == 0 {
			return value{kind: stringKind, s: subject}, nil
		}
		if camelBoundary {
			var out strings.Builder
			for i, r := range subject {
				if i > 0 && r >= 'A' && r <= 'Z' && utf8.RuneCountInString(subject[:i]) > 0 {
					out.WriteString(replacement[2:])
					count++
				}
				out.WriteRune(r)
			}
			return value{kind: stringKind, s: out.String()}, nil
		}
		if tableBoundary {
			var out strings.Builder
			previous := rune(0)
			for _, r := range subject {
				if previous != 0 && r >= 'A' && r <= 'Z' && (previous >= 'a' && previous <= 'z' || previous >= 'A' && previous <= 'Z' || previous >= '0' && previous <= '9' || previous == '_') {
					out.WriteByte('_')
					count++
				}
				out.WriteRune(r)
				previous = r
			}
			return value{kind: stringKind, s: out.String()}, nil
		}
		if usePCRE2 {
			matches, _, err := pcre2MatchAll(pattern, subject)
			if err != nil {
				return value{}, err
			}
			if limit >= 0 && int64(len(matches)) > limit {
				matches = matches[:limit]
			}
			var out strings.Builder
			last := 0
			for _, match := range matches {
				if match[0] < last {
					continue
				}
				out.WriteString(subject[last:match[0]])
				out.WriteString(expandPCREReplacement(replacement, subject, match))
				last = match[1]
				count++
			}
			out.WriteString(subject[last:])
			return value{kind: stringKind, s: out.String()}, nil
		}
		if limit >= 0 {
			matches := re.FindAllStringSubmatchIndex(subject, int(limit))
			var out strings.Builder
			last := 0
			for _, match := range matches {
				out.WriteString(subject[last:match[0]])
				out.Write(re.ExpandString(nil, replacement, subject, match))
				last = match[1]
				count++
			}
			out.WriteString(subject[last:])
			return value{kind: stringKind, s: out.String()}, nil
		}
		count += int64(len(re.FindAllStringIndex(subject, -1)))
		return value{kind: stringKind, s: re.ReplaceAllString(subject, replacement)}, nil
	}
	var result value
	if args[2].value.kind == arrayKind {
		result = arrayValue()
		for _, entry := range args[2].value.a.entries {
			item, err := convert(dereference(entry.value))
			if err != nil {
				return value{}, err
			}
			result.a = result.a.with(entry.key, item)
		}
	} else {
		result, err = convert(args[2].value)
		if err != nil {
			return value{}, err
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
