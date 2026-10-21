package php

import (
	"fmt"
	"math"
	"strings"
)

func (e *evaluator) callNumberFormat(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "num"}, {name: "decimals"}, {name: "decimal_separator"}, {name: "thousands_separator"}}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 4 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "number_format() expects 1 to 4 arguments"}
	}
	n, _, ok := e.coerceType(bound[0], "float", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "number_format", 1, "num", "float", bound[0])
	}
	decimals := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "number_format", 2, "decimals", "int", bound[1])
		}
		decimals = v.i
	}
	if decimals > 100 || decimals < -100 {
		return value{}, &runtimeError{node.line, "ValueError", "number_format(): Argument #2 ($decimals) must be between -100 and 100"}
	}
	decimalSep, thousandsSep := ".", ","
	if present[2] {
		if bound[2].kind == nullKind {
			decimalSep = ""
		} else {
			v, _, ok := e.coerceType(bound[2], "string", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "number_format", 3, "decimal_separator", "?string", bound[2])
			}
			decimalSep = v.s
		}
	}
	if present[3] {
		if bound[3].kind == nullKind {
			thousandsSep = ""
		} else {
			v, _, ok := e.coerceType(bound[3], "string", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "number_format", 4, "thousands_separator", "?string", bound[3])
			}
			thousandsSep = v.s
		}
	}
	f := n.f
	if n.kind == intKind {
		f = float64(n.i)
	}
	if decimals < 0 {
		factor := math.Pow10(int(-decimals))
		f = math.Round(f/factor) * factor
		decimals = 0
	}
	factor := math.Pow10(int(decimals))
	f = math.Round(f*factor) / factor
	if f == 0 { f = 0 } // PHP 7.2+ never formats negative zero.
	formatted := fmt.Sprintf("%.*f", int(decimals), f)
	parts := strings.SplitN(formatted, ".", 2)
	whole := parts[0]
	sign := ""
	if strings.HasPrefix(whole, "-") {
		sign = "-"
		whole = whole[1:]
	}
	groups := make([]string, 0, (len(whole)+2)/3)
	for len(whole) > 3 {
		groups = append([]string{whole[len(whole)-3:]}, groups...)
		whole = whole[:len(whole)-3]
	}
	groups = append([]string{whole}, groups...)
	out := sign + strings.Join(groups, thousandsSep)
	if len(parts) == 2 {
		out += decimalSep + parts[1]
	}
	return value{kind: stringKind, s: out}, nil
}
