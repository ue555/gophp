package php

import (
	"strings"
	"unicode/utf8"
)

func (e *evaluator) callPregSplit(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "pattern"}, {name: "subject"}, {name: "limit"}, {name: "flags"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 4 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "preg_split() expects 2 to 4 arguments"}
	}
	p, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "preg_split", 1, "pattern", "string", bound[0])
	}
	s, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "preg_split", 2, "subject", "string", bound[1])
	}
	limit := -1
	if present[2] {
		v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_split", 3, "limit", "int", bound[2])
		}
		limit = int(v.i)
		if limit == 0 {
			limit = -1
		}
	}
	flags := int64(0)
	if present[3] {
		v, _, ok := e.coerceType(bound[3], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "preg_split", 4, "flags", "int", bound[3])
		}
		flags = v.i
	}
	if flags&^int64(7) != 0 {
		return value{}, &runtimeError{node.line, "Error", "gophp: preg_split() flag is not implemented"}
	}
	var matches [][]int
	if p.s == "/(.)(?=[A-Z])/u" {
		for i := 1; i < len(s.s); i++ {
			if s.s[i] < 'A' || s.s[i] > 'Z' {
				continue
			}
			_, width := utf8.DecodeLastRuneInString(s.s[:i])
			if width > 0 {
				matches = append(matches, []int{i - width, i, i - width, i})
			}
		}
	} else {
		re, err := compileBasicPCRE(p.s)
		if err != nil {
			indices, _, pcreErr := pcre2MatchAll(p.s, s.s)
			if pcreErr != nil {
				return value{}, pcreErr
			}
			matches = indices
		} else {
			maxMatches := -1
			if limit > 0 {
				maxMatches = limit - 1
			}
			matches = re.FindAllStringSubmatchIndex(s.s, maxMatches)
		}
	}
	out := arrayValue()
	appendPart := func(part string, start int) {
		if flags&1 != 0 && part == "" {
			return
		}
		item := value{kind: stringKind, s: part}
		if flags&4 != 0 {
			pair := arrayValue()
			pair.a = pair.a.with(arrayKey{index: 0}, item)
			pair.a = pair.a.with(arrayKey{index: 1}, integer(int64(start)))
			item = pair
		}
		out.a = out.a.with(arrayKey{index: out.a.next}, item)
	}
	if limit > 0 && len(matches) > limit-1 {
		matches = matches[:limit-1]
	}
	last := 0
	for _, match := range matches {
		if match[0] < last {
			continue
		}
		appendPart(s.s[last:match[0]], last)
		if flags&2 != 0 {
			for j := 2; j+1 < len(match); j += 2 {
				if match[j] >= 0 {
					appendPart(s.s[match[j]:match[j+1]], match[j])
				}
			}
		}
		last = match[1]
	}
	appendPart(s.s[last:], last)
	// An empty UTF-8 pattern splits between code points. regexp.Split already
	// respects UTF-8 boundaries, but an empty subject has no elements in PHP.
	if s.s == "" && strings.HasPrefix(p.s, "//") {
		return arrayValue(), nil
	}
	return out, nil
}
