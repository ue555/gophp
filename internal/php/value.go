package php

import (
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

type valueKind uint8

const (
	nullKind valueKind = iota
	intKind
	floatKind
	stringKind
	boolKind
	arrayKind
	objectKind
	// refKind is internal to variable and array-element bindings; evaluated PHP values never expose it.
	refKind
	resourceKind
)

type value struct {
	kind   valueKind
	i      int64
	f      float64
	s      string
	a      *phpArray
	o      *phpObject
	ref    *referenceCell
	stream *phpStream
}

func integer(i int64) value    { return value{kind: intKind, i: i} }
func floating(f float64) value { return value{kind: floatKind, f: f} }

func boolean(b bool) value {
	v := value{kind: boolKind}
	if b {
		v.i = 1
	}
	return v
}

// truthy implements PHP's boolean conversion for all supported value types.
// Numeric-looking strings are not parsed: only "" and "0" are false strings.
func (v value) truthy() bool {
	switch v.kind {
	case nullKind:
		return false
	case boolKind, intKind:
		return v.i != 0
	case floatKind:
		return v.f != 0 // Both signed zeros are false; NaN is true.
	case stringKind:
		return v.s != "" && v.s != "0"
	case arrayKind:
		return len(v.a.entries) != 0
	case objectKind:
		return true
	case resourceKind:
		return true
	}
	panic("invalid value kind")
}

// identical compares the PHP type first, then the value without coercion.
func (v value) identical(other value) bool {
	if v.kind != other.kind {
		return false
	}
	switch v.kind {
	case nullKind:
		return true
	case boolKind, intKind:
		return v.i == other.i
	case floatKind:
		return v.f == other.f // NaN is never identical, including to itself.
	case stringKind:
		return v.s == other.s
	case objectKind:
		return v.o == other.o
	case resourceKind:
		return v.stream == other.stream
	case arrayKind:
		if v.a == other.a {
			return true
		}
		if len(v.a.entries) != len(other.a.entries) {
			return false
		}
		for i, entry := range v.a.entries {
			otherEntry := other.a.entries[i]
			if entry.key != otherEntry.key || !dereference(entry.value).identical(dereference(otherEntry.value)) {
				return false
			}
		}
		return true
	}
	panic("invalid value kind")
}
func (v value) typeName() string {
	if v.kind == resourceKind {
		return "resource"
	}
	if v.kind == objectKind {
		return v.o.className()
	}
	return [...]string{"null", "int", "float", "string", "bool", "array", "stdClass"}[v.kind]
}

// Diagnostics describing an actual argument distinguish true and false;
// gettype() and operator type descriptions still report bool/boolean.
func (v value) diagnosticName() string {
	if v.kind == boolKind {
		if v.i == 0 {
			return "false"
		}
		return "true"
	}
	return v.typeName()
}
func (v value) asFloat() float64 {
	if v.kind == floatKind {
		return v.f
	}
	return float64(v.i)
}

func parseNumber(s string) value {
	if !strings.ContainsAny(s, ".eE") {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return integer(n)
		}
	}
	f, _ := strconv.ParseFloat(s, 64)
	return floating(f)
}

func parseLiteralNumber(s string) value {
	s = strings.ReplaceAll(s, "_", "")
	if len(s) > 1 && s[0] == '0' && !strings.ContainsAny(s, ".eE") {
		base := 8
		digits := s[1:]
		if len(s) > 2 && (s[1] == 'x' || s[1] == 'X') {
			base, digits = 16, s[2:]
		} else if len(s) > 2 && (s[1] == 'b' || s[1] == 'B') {
			base, digits = 2, s[2:]
		} else if len(s) > 2 && (s[1] == 'o' || s[1] == 'O') {
			digits = s[2:]
		}
		if n, err := strconv.ParseInt(digits, base, 64); err == nil {
			return integer(n)
		}
	}
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		if n, err := strconv.ParseInt(s[2:], 16, 64); err == nil {
			return integer(n)
		}
	}
	if !strings.ContainsAny(s, ".eE") {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return integer(n)
		}
	}
	// The caller validates syntax. Range errors intentionally produce +/-Inf,
	// matching PHP's conversion to a double for out-of-range decimal values.
	f, _ := strconv.ParseFloat(s, 64)
	return floating(f)
}

func (v value) text() string {
	switch v.kind {
	case nullKind:
		return ""
	case boolKind:
		if v.i != 0 {
			return "1"
		}
		return ""
	case intKind:
		return strconv.FormatInt(v.i, 10)
	case stringKind:
		return v.s
	case arrayKind:
		return "Array"
	case floatKind:
		return floatText(v.f, 14)
	case resourceKind:
		return "Resource id #1"
	}
	panic("invalid value kind")
}

func floatText(f float64, precision int) string {
	if math.IsNaN(f) {
		return "NAN"
	}
	if math.IsInf(f, 1) {
		return "INF"
	}
	if math.IsInf(f, -1) {
		return "-INF"
	}
	// Echo uses precision 14; diagnostic floats use shortest round-trip text.
	s := strconv.FormatFloat(f, 'G', precision, 64)
	if index := strings.IndexByte(s, 'E'); index >= 0 {
		mantissa := s[:index]
		if !strings.Contains(mantissa, ".") {
			mantissa += ".0"
		}
		exponent, _ := strconv.Atoi(s[index+1:])
		sign := "+"
		if exponent < 0 {
			sign = "-"
			exponent = -exponent
		}
		return mantissa + "E" + sign + strconv.Itoa(exponent)
	}
	return s
}

var numericPrefix = regexp.MustCompile(`^[ \t\r\n\v\f]*[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?`)

// numeric returns the converted value, whether a trailing-junk warning is
// needed, and whether conversion is possible. Only PHP's ASCII whitespace is
// accepted; Go's Unicode whitespace rules differ.
func (v value) numeric() (value, bool, bool) {
	switch v.kind {
	case intKind, floatKind:
		return v, false, true
	case nullKind, boolKind:
		return integer(v.i), false, true
	case stringKind:
		prefix := numericPrefix.FindString(v.s)
		if prefix == "" {
			return value{}, false, false
		}
		trailing := strings.Trim(v.s[len(prefix):], " \t\r\n\v\f") != ""
		return parseNumber(strings.TrimLeft(prefix, " \t\r\n\v\f")), trailing, true
	case arrayKind, objectKind:
		return value{}, false, false
	}
	panic("invalid value kind")
}

func arithmetic(op string, a, b value) value {
	if a.kind == intKind && b.kind == intKind {
		x, y := a.i, b.i
		switch op {
		case "**":
			if y >= 0 && y <= 63 {
				result := new(big.Int).Exp(big.NewInt(x), big.NewInt(y), nil)
				if result.IsInt64() {
					return integer(result.Int64())
				}
			}
		case "+":
			z := x + y
			if !((y > 0 && z < x) || (y < 0 && z > x)) {
				return integer(z)
			}
		case "-":
			z := x - y
			if !((y < 0 && z < x) || (y > 0 && z > x)) {
				return integer(z)
			}
		case "*":
			if x == 0 || y == 0 {
				return integer(0)
			}
			z := x * y
			if !(x == math.MinInt64 && y == -1 || y == math.MinInt64 && x == -1) && z/y == x {
				return integer(z)
			}
		case "/":
			if !(x == math.MinInt64 && y == -1) && x%y == 0 {
				return integer(x / y)
			}
		}
	}
	x, y := a.asFloat(), b.asFloat()
	switch op {
	case "**":
		return floating(math.Pow(x, y))
	case "+":
		return floating(x + y)
	case "-":
		return floating(x - y)
	case "*":
		return floating(x * y)
	case "/":
		return floating(x / y)
	}
	panic("invalid arithmetic operator")
}
