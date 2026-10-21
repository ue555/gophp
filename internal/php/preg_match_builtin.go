package php

import (
	"regexp"
	"strconv"
	"unicode/utf8"
)

var largeDotPattern = regexp.MustCompile(`^/\.\{1,([0-9]+)\}/([us]*)$`)

// RE2 rejects counted repetitions above 1000. A plain dot run has a
// direct linear implementation, including UTF-8 byte offsets.
func largeDotMatch(pattern, subject string) ([]int, bool) {
	parts := largeDotPattern.FindStringSubmatch(pattern)
	if parts == nil {
		return nil, false
	}
	limit, err := strconv.Atoi(parts[1])
	if err != nil || limit < 1 {
		return nil, false
	}
	start := 0
	if !containsByte(parts[2], 's') {
		for start < len(subject) && subject[start] == '\n' {
			start++
		}
	}
	length, count := start, 0
	for length < len(subject) && count < limit {
		if subject[length] == '\n' && !containsByte(parts[2], 's') {
			break
		}
		step := 1
		if containsByte(parts[2], 'u') {
			_, step = utf8.DecodeRuneInString(subject[length:])
		}
		length += step
		count++
	}
	if count == 0 {
		return nil, true
	}
	return []int{start, length}, true
}

func containsByte(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

func (e *evaluator) callPregMatch(node functionCall) (value, error) {
	if len(node.args) < 2 || len(node.args) > 5 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "preg_match() expects 2 to 5 arguments"}
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
		return value{}, builtinTypeError(node.line, "preg_match", 1, "pattern", "string", pattern)
	}
	s, _, ok := e.coerceType(subject, "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "preg_match", 2, "subject", "string", subject)
	}
	flags := int64(0)
	if len(node.args) >= 4 {
		v, err := e.eval(node.args[3])
		if err != nil {
			return value{}, err
		}
		converted, _, ok := e.coerceType(v, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_match", 4, "flags", "int", v)
		}
		flags = converted.i
	}
	offset := int64(0)
	if len(node.args) >= 5 {
		v, err := e.eval(node.args[4])
		if err != nil {
			return value{}, err
		}
		converted, _, ok := e.coerceType(v, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_match", 5, "offset", "int", v)
		}
		offset = converted.i
	}
	if offset < 0 {
		offset += int64(len(s.s))
	}
	if offset < 0 || offset > int64(len(s.s)) {
		return value{}, &runtimeError{node.line, "ValueError", "preg_match(): Argument #5 ($offset) must be contained in argument #2 ($subject)"}
	}
	e.pregLastError = 0
	matches, names, err := pcre2Matches(p.s, s.s, int(offset), 1)
	if err != nil {
		return e.pregFailure("preg_match", node.line, err)
	}
	var index []int
	if len(matches) > 0 {
		index = matches[0]
		for i, v := range index {
			if v >= 0 {
				index[i] = v - int(offset)
			}
		}
	}

	if len(node.args) >= 3 {
		out := arrayValue()
		if index != nil {
			for i := 0; i < len(index)/2; i++ {
				start, stop := index[2*i], index[2*i+1]
				if start < 0 && flags&512 == 0 {
					continue
				}
				item := value{kind: nullKind}
				if start >= 0 {
					item = value{kind: stringKind, s: s.s[int(offset)+start : int(offset)+stop]}
				}
				if flags&256 != 0 && start >= 0 {
					pair := arrayValue()
					pair.a = pair.a.with(arrayKey{index: 0}, item)
					pair.a = pair.a.with(arrayKey{index: 1}, integer(int64(start)+offset))
					item = pair
				}
				out.a = out.a.with(arrayKey{index: int64(i)}, item)
				if i < len(names) && names[i] != "" {
					name := names[i]
					out.a = out.a.with(arrayKey{isString: true, text: name}, item)
				}
			}
		}
		if !writable(node.args[2]) {
			return value{}, &runtimeError{node.line, "Error", "preg_match(): Argument #3 ($matches) could not be passed by reference"}
		}
		if _, err := e.evalAssignment(assignment{target: node.args[2], source: literal{value: out, line: node.line}, line: node.line}); err != nil {
			return value{}, err
		}
	}
	if index != nil {
		return integer(1), nil
	}
	return integer(0), nil
}
