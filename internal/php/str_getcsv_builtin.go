package php

import "strings"

func (e *evaluator) callStrGetCSV(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "string"}, {name: "separator"}, {name: "enclosure"}, {name: "escape"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 4 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "str_getcsv() expects 1 to 4 arguments"}
	}
	text, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "str_getcsv", 1, "string", "string", bound[0])
	}
	delim, enclosure, escape := byte(','), byte('"'), byte('\\')
	for index, target := range []*byte{&delim, &enclosure, &escape} {
		if !present[index+1] {
			continue
		}
		v, _, ok := e.coerceType(bound[index+1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "str_getcsv", index+2, []string{"separator", "enclosure", "escape"}[index], "string", bound[index+1])
		}
		if len(v.s) != 1 {
			return value{}, &runtimeError{node.line, "ValueError", "str_getcsv(): Argument must be a single character"}
		}
		*target = v.s[0]
	}
	out := arrayValue()
	var item strings.Builder
	inQuotes := false
	for i := 0; i < len(text.s); i++ {
		c := text.s[i]
		if inQuotes {
			if c == escape && i+1 < len(text.s) && text.s[i+1] == enclosure {
				item.WriteByte(enclosure)
				i++
				continue
			}
			if c == enclosure {
				if i+1 < len(text.s) && text.s[i+1] == enclosure {
					item.WriteByte(enclosure)
					i++
				} else {
					inQuotes = false
				}
				continue
			}
			item.WriteByte(c)
			continue
		}
		if c == enclosure && item.Len() == 0 {
			inQuotes = true
			continue
		}
		if c == delim {
			out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, value{kind: stringKind, s: item.String()})
			item.Reset()
			continue
		}
		item.WriteByte(c)
	}
	last := value{kind: stringKind, s: item.String()}
	if text.s == "" {
		last = value{kind: nullKind}
	}
	out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, last)
	return out, nil
}
