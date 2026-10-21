package php

import (
	"regexp"
	"strings"
)

// compileBasicPCRE covers the delimiter/modifier and character-class forms
// needed by the current console path. Constructs RE2 cannot represent fail
// explicitly instead of silently changing their meaning.
func compileBasicPCRE(pattern string) (*regexp.Regexp, error) {
	pattern = strings.TrimSpace(pattern)
	if len(pattern) < 2 {
		return nil, &runtimeError{0, "Error", "gophp: unsupported regular expression delimiter"}
	}
	delimiter := pattern[0]
	closing := delimiter
	switch delimiter {
	case '{':
		closing = '}'
	case '[':
		closing = ']'
	case '(':
		closing = ')'
	case '<':
		closing = '>'
	}
	end := strings.LastIndexByte(pattern, closing)
	if end <= 0 {
		return nil, &runtimeError{0, "Error", "gophp: unsupported regular expression delimiter"}
	}
	body, flags := pattern[1:end], pattern[end+1:]
	var result strings.Builder
	insideClass, escaped := false, false
	for i := 0; i < len(body); i++ {
		c := body[i]
		if escaped {
			result.WriteByte(c)
			escaped = false
			continue
		}
		if c == '\\' {
			result.WriteByte(c)
			escaped = true
			continue
		}
		if c == '[' {
			insideClass = true
		}
		if c == ']' {
			insideClass = false
		}
		if !insideClass && strings.Contains(flags, "x") && (c == ' ' || c == '\t' || c == '\n') {
			continue
		}
		// PCRE possessive quantifiers have the same match set as greedy
		// quantifiers when no following pattern can backtrack into them.
		if !insideClass && c == '+' && i > 0 && (body[i-1] == '*' || body[i-1] == '+' || body[i-1] == '?') {
			continue
		}
		result.WriteByte(c)
	}
	body = result.String()
	modifiers := ""
	for _, flag := range flags {
		if strings.ContainsRune("ims", flag) {
			modifiers += string(flag)
		} else if flag != 'x' && flag != 'u' && flag != 'D' && flag != 'A' {
			return nil, &runtimeError{0, "Error", "gophp: unsupported regular expression modifier"}
		}
	}
	if modifiers != "" {
		body = "(?" + modifiers + ")" + body
	}
	if strings.Contains(flags, "A") && !strings.HasPrefix(body, "^") {
		body = "^" + body
	}
	body = strings.ReplaceAll(body, "(?<", "(?P<")
	re, err := regexp.Compile(body)
	if err != nil {
		return nil, &runtimeError{0, "Error", "gophp: unsupported regular expression: " + err.Error()}
	}
	return re, nil
}

func (e *evaluator) callPregMatchAll(node functionCall) (value, error) {
	if len(node.args) < 2 || len(node.args) > 5 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "preg_match_all() expects 2 to 5 arguments"}
	}
	pattern, err := e.eval(node.args[0])
	if err != nil {
		return value{}, err
	}
	subject, err := e.eval(node.args[1])
	if err != nil {
		return value{}, err
	}
	p, _, ok := e.coerceType(pattern, "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "preg_match_all", 1, "pattern", "string", pattern)
	}
	s, _, ok := e.coerceType(subject, "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "preg_match_all", 2, "subject", "string", subject)
	}
	flags := int64(0)
	if len(node.args) > 3 {
		v, err := e.eval(node.args[3])
		if err != nil {
			return value{}, err
		}
		converted, _, ok := e.coerceType(v, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_match_all", 4, "flags", "int", v)
		}
		flags = converted.i
	}
	offset := int64(0)
	if len(node.args) > 4 {
		v, err := e.eval(node.args[4])
		if err != nil {
			return value{}, err
		}
		converted, _, ok := e.coerceType(v, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_match_all", 5, "offset", "int", v)
		}
		offset = converted.i
	}
	if offset < 0 {
		offset += int64(len(s.s))
	}
	if offset < 0 || offset > int64(len(s.s)) {
		return value{}, &runtimeError{node.line, "ValueError", "preg_match_all(): Argument #5 ($offset) must be contained in argument #2 ($subject)"}
	}
	e.pregLastError = 0
	indices, names, err := pcre2Matches(p.s, s.s, int(offset), -1)
	if err != nil {
		return e.pregFailure("preg_match_all", node.line, err)
	}
	for _, index := range indices {
		for i, v := range index {
			if v >= 0 {
				index[i] = v - int(offset)
			}
		}
	}

	out := arrayValue()
	makeItem := func(start, stop int) value {
		position := int64(start) + offset
		if start < 0 {
			position = -1
			if flags&512 != 0 {
				return value{kind: nullKind}
			}
			start, stop = 0, 0
		}
		item := value{kind: stringKind, s: s.s[int(offset)+start : int(offset)+stop]}
		if flags&256 == 0 {
			return item
		}
		pair := arrayValue()
		pair.a = pair.a.with(arrayKey{index: 0}, item)
		pair.a = pair.a.with(arrayKey{index: 1}, integer(position))
		return pair
	}
	if flags&2 != 0 {
		for matchIndex, match := range indices {
			row := arrayValue()
			last := len(match)/2 - 1
			if flags&512 == 0 {
				for last > 0 && match[last*2] < 0 {
					last--
				}
			}
			for group := 0; group <= last; group++ {
				item := makeItem(match[group*2], match[group*2+1])
				row.a = row.a.with(arrayKey{index: int64(group)}, item)
				if name := names[group]; name != "" {
					row.a = row.a.with(arrayKey{isString: true, text: name}, item)
				}
			}
			out.a = out.a.with(arrayKey{index: int64(matchIndex)}, row)
		}
	} else {
		for group := 0; group < len(names); group++ {
			column := arrayValue()
			for matchIndex, match := range indices {
				column.a = column.a.with(arrayKey{index: int64(matchIndex)}, makeItem(match[group*2], match[group*2+1]))
			}
			out.a = out.a.with(arrayKey{index: int64(group)}, column)
			if name := names[group]; name != "" {
				out.a = out.a.with(arrayKey{isString: true, text: name}, column)
			}
		}
	}
	if len(node.args) >= 3 {
		if !writable(node.args[2]) {
			return value{}, &runtimeError{node.line, "Error", "preg_match_all(): Argument #3 ($matches) could not be passed by reference"}
		}
		if _, err := e.evalAssignment(assignment{target: node.args[2], source: literal{value: out, line: node.line}, line: node.line}); err != nil {
			return value{}, err
		}
	}
	return integer(int64(len(indices))), nil
}
