package php

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Explicit casts permit leading numeric strings and do not issue the implicit
// conversion deprecation for fractional values. PHP 8.5 still warns on a float
// outside the integer range, and on NaN converted to bool/string.
func (e *evaluator) cast(v value, name string, line int) (value, error) {
	switch name {
	case "bool":
		return boolean(e.toBool(v, line)), nil
	case "string":
		s, err := e.toText(v, line)
		return value{kind: stringKind, s: s}, err
	case "int", "float":
		if v.kind == objectKind {
			e.warn(line, fmt.Sprintf("Object of class %s could not be converted to %s", v.o.className(), name))
			v = integer(1)
		} else if v.kind == arrayKind {
			v = boolean(len(v.a.entries) > 0)
		}
		if name == "float" && v.kind == stringKind {
			prefix := strings.TrimLeft(numericPrefix.FindString(v.s), " \t\r\n\v\f")
			f, _ := strconv.ParseFloat(prefix, 64)
			return floating(f), nil
		}
		fromString := v.kind == stringKind
		n, _, ok := v.numeric()
		if !ok {
			n = integer(0)
		}
		if name == "float" {
			return floating(n.asFloat()), nil
		}
		if n.kind == intKind {
			return n, nil
		}
		if fromString {
			// Numeric strings saturate, while native doubles wrap modulo 2^64.
			if math.IsInf(n.f, 0) || math.IsNaN(n.f) {
				return integer(0), nil
			}
			if n.f >= math.Ldexp(1, 63) {
				return integer(math.MaxInt64), nil
			}
			if n.f < -math.Ldexp(1, 63) {
				return integer(math.MinInt64), nil
			}
		} else if e.version == PHP85 && floatOutsideInt(n.f) {
			e.warn(line, "The float "+floatText(n.f, -1)+" is not representable as an int, cast occurred")
		}
		return integer(floatToInt(n.f)), nil
	}
	panic("invalid cast")
}

func floatOutsideInt(f float64) bool {
	return math.IsNaN(f) || math.IsInf(f, 0) || f >= math.Ldexp(1, 63) || f < -math.Ldexp(1, 63)
}
