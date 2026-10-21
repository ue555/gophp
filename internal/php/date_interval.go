package php

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var intervalISO = regexp.MustCompile(`^P(?:(\d+)Y)?(?:(\d+)M)?(?:(\d+)W)?(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)

func (e *evaluator) invokeDateInterval(object value, args []callArgument, line int) (value, error) {
	b, p, _, err := bindCallArguments(args, []parameter{{name: "duration"}}, line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !p[0] {
		return value{}, &runtimeError{line, "ArgumentCountError", "DateInterval::__construct() expects exactly 1 argument"}
	}
	v, _, ok := e.coerceType(b[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(line, "DateInterval::__construct", 1, "duration", "string", b[0])
	}
	m := intervalISO.FindStringSubmatch(v.s)
	if m == nil || v.s == "P" || v.s == "PT" || strings.HasSuffix(v.s, "T") {
		return value{}, &runtimeError{line, "DateMalformedIntervalStringException", "Unknown or bad format (" + v.s + ")"}
	}
	fields := map[string]value{"f": floating(0), "invert": integer(0), "days": boolean(false), "from_string": boolean(false)}
	for i, n := range []string{"y", "m", "weeks", "d", "h", "i", "s"} {
		x, _ := strconv.ParseInt(m[i+1], 10, 64)
		fields[n] = integer(x)
	}
	fields["d"] = integer(fields["d"].i + 7*fields["weeks"].i)
	delete(fields, "weeks")
	for _, n := range []string{"y", "m", "d", "h", "i", "s", "f", "invert", "days", "from_string"} {
		object.o.properties = object.o.properties.with(arrayKey{isString: true, text: n}, fields[n])
	}
	return value{}, nil
}
func (e *evaluator) dateIntervalChange(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	b, p, _, err := bindCallArguments(args, []parameter{{name: "interval"}}, line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !p[0] {
		return value{}, &runtimeError{line, "ArgumentCountError", method.owner.name + "::" + method.name + "() expects exactly 1 argument"}
	}
	v := b[0]
	ok := false
	if v.kind == objectKind {
		for c := v.o.class; c != nil; c = c.parent {
			if c.name == "DateInterval" {
				ok = true
				break
			}
		}
	}
	if !ok {
		return value{}, builtinTypeError(line, method.owner.name+"::"+method.name, 1, "interval", "DateInterval", v)
	}
	if object.o.dateTime == nil {
		return value{}, &runtimeError{line, "Error", "DateTime object not initialized"}
	}
	field := func(n string) value { x, _ := v.o.properties.get(arrayKey{isString: true, text: n}); return x }
	sign := int64(1)
	if field("invert").truthy() {
		sign = -sign
	}
	if method.name == "sub" {
		sign = -sign
	}
	updated := object.o.dateTime.AddDate(int(sign*field("y").i), int(sign*field("m").i), int(sign*field("d").i))
	seconds := field("h").i*3600 + field("i").i*60 + field("s").i
	updated = updated.Add(time.Duration(sign*seconds)*time.Second + time.Duration(float64(sign)*field("f").asFloat()*1e9))
	if method.owner.name == "DateTimeImmutable" {
		clone := *object.o
		clone.dateTime = &updated
		return value{kind: objectKind, o: &clone}, nil
	}
	object.o.dateTime = &updated
	return object, nil
}
