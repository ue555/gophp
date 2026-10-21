package php

import (
	"strconv"
	"strings"
	"time"
)

func (e *evaluator) callDateFormat(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "format"}, {name: "timestamp"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects 1 to 2 arguments"}
	}
	format, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, name, 1, "format", "string", bound[0])
	}
	stamp := time.Now().Unix()
	if present[1] && bound[1].kind != nullKind {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, name, 2, "timestamp", "?int", bound[1])
		}
		stamp = v.i
	}
	t := time.Unix(stamp, 0).UTC()
	if name == "date" {
		if location, err := time.LoadLocation(e.defaultTimezone); err == nil {
			t = t.In(location)
		}
	}
	var out strings.Builder
	for i := 0; i < len(format.s); i++ {
		c := format.s[i]
		if c == '\\' && i+1 < len(format.s) {
			i++
			out.WriteByte(format.s[i])
			continue
		}
		switch c {
		case 'T':
			zone, _ := t.Zone()
			if name == "gmdate" {
				zone = "GMT"
			}
			out.WriteString(zone)
		case 'D':
			out.WriteString(t.Format("Mon"))
		case 'd':
			out.WriteString(t.Format("02"))
		case 'M':
			out.WriteString(t.Format("Jan"))
		case 'Y':
			out.WriteString(t.Format("2006"))
		case 'H':
			out.WriteString(t.Format("15"))
		case 'i':
			out.WriteString(t.Format("04"))
		case 's':
			out.WriteString(t.Format("05"))
		case 'U':
			out.WriteString(strconv.FormatInt(stamp, 10))
		default:
			out.WriteByte(c)
		}
	}
	return value{kind: stringKind, s: out.String()}, nil
}
