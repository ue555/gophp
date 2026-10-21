package php

import (
	"strings"
	"time"
)

func (e *evaluator) callDefaultTimezone(node functionCall) (value, error) {
	if strings.EqualFold(node.name, "date_default_timezone_get") {
		if len(node.args) != 0 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "date_default_timezone_get() expects exactly 0 arguments"}
		}
		return value{kind: stringKind, s: e.defaultTimezone}, nil
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "timezoneId"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "date_default_timezone_set() expects exactly 1 argument"}
	}
	zone, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "date_default_timezone_set", 1, "timezoneId", "string", bound[0])
	}
	if _, err := time.LoadLocation(zone.s); err != nil {
		return value{}, &runtimeError{node.line, "ValueError", "date_default_timezone_set(): Argument #1 ($timezoneId) must be a valid timezone identifier"}
	}
	e.defaultTimezone = zone.s
	return boolean(true), nil
}
