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
	case "object":
		if v.kind == objectKind {
			return v, nil
		}
		e.initClasses()
		object := value{kind: objectKind, o: &phpObject{class: e.classes["stdclass"], properties: &phpArray{}}}
		if v.kind == arrayKind {
			for _, entry := range v.a.entries {
				key := entry.key.text
				if !entry.key.isString {
					key = strconv.FormatInt(entry.key.index, 10)
				}
				object.o.properties = object.o.properties.with(arrayKey{isString: true, text: key}, entry.value)
			}
		} else if v.kind != nullKind {
			object.o.properties = object.o.properties.with(arrayKey{isString: true, text: "scalar"}, v)
		}
		return object, nil
	case "array":
		switch v.kind {
		case arrayKind:
			return v, nil
		case nullKind:
			return arrayValue(), nil
		case objectKind:
			result := arrayValue()
			if v.o.properties != nil {
				for _, item := range v.o.properties.entries {
					result.a = result.a.with(item.key, item.value)
				}
			}
			return result, nil
		default:
			result := arrayValue()
			result.a = result.a.with(arrayKey{index: 0}, v)
			return result, nil
		}
	case "bool":
		return boolean(e.toBool(v, line)), nil
	case "string":
		s, err := e.toText(v, line)
		return value{kind: stringKind, s: s}, err
	case "int", "float":
		if v.kind == objectKind {
			if err := e.warn(line, fmt.Sprintf("Object of class %s could not be converted to %s", v.o.className(), name)); err != nil {
				return value{}, err
			}
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
			if err := e.warn(line, "The float "+floatText(n.f, -1)+" is not representable as an int, cast occurred"); err != nil {
				return value{}, err
			}
		}
		return integer(floatToInt(n.f)), nil
	}
	panic("invalid cast")
}

func floatOutsideInt(f float64) bool {
	return math.IsNaN(f) || math.IsInf(f, 0) || f >= math.Ldexp(1, 63) || f < -math.Ldexp(1, 63)
}
