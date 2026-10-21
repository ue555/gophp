package php

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// The native DateTime entry points used by Carbon's startup path. Date/time
// parsing and formatting outside these forms still needs the extension work.
func (e *evaluator) invokeDateTime(object value, called *classDecl, method *methodDecl, args []callArgument, line int) (value, error) {
	if strings.EqualFold(method.name, "createFromFormat") {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "format"}, {name: "datetime"}, {name: "timezone"}}, line)
		if err != nil {
			return value{}, err
		}
		if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
			return value{}, &runtimeError{line, "ArgumentCountError", "DateTime::createFromFormat() expects 2 to 3 arguments"}
		}
		layout := ""
		switch bound[0].text() {
		case "Y-m-d H:i:s":
			layout = "2006-01-02 15:04:05"
		case "Y-m-d":
			layout = "2006-01-02"
		case "U.u":
			layout = "U.u"
		default:
			return value{}, &runtimeError{line, "Error", "gophp: unsupported DateTime::createFromFormat pattern"}
		}
		loc := time.UTC
		if present[2] && bound[2].kind == objectKind && bound[2].o != nil && bound[2].o.timezone != nil {
			loc = bound[2].o.timezone
		}
		var parsed time.Time
		var parseErr error
		if bound[0].text() == "U.u" {
			parts := strings.SplitN(bound[1].text(), ".", 2)
			var seconds int64
			seconds, parseErr = strconv.ParseInt(parts[0], 10, 64)
			if parseErr == nil {
				fraction := "000000"
				if len(parts) == 2 {
					fraction = (parts[1] + fraction)[:6]
				}
				var micros int64
				micros, parseErr = strconv.ParseInt(fraction, 10, 64)
				if parseErr == nil {
					parsed = time.Unix(seconds, micros*1000).In(loc)
				}
			}
		} else {
			parsed, parseErr = time.ParseInLocation(layout, bound[1].text(), loc)
		}
		if parseErr != nil {
			return boolean(false), nil
		}
		e.initClasses()
		if called == nil {
			called = method.owner
		}
		instance := value{kind: objectKind, o: &phpObject{class: called, properties: &phpArray{}, dateTime: &parsed}}
		for cls := called; cls != nil; cls = cls.parent {
			for _, property := range cls.properties {
				if property.static || property.initial == nil {
					continue
				}
				previousScope := e.scope
				e.scope = property.owner
				initial, evalErr := e.eval(property.initial)
				e.scope = previousScope
				if evalErr != nil {
					return value{}, evalErr
				}
				instance.o.properties = instance.o.properties.with(propertyStorageKey(property), initial)
			}
		}
		return instance, nil
	}
	if strings.EqualFold(method.name, "getLastErrors") {
		if len(args) != 0 {
			return value{}, &runtimeError{line, "ArgumentCountError", "DateTime::getLastErrors() expects exactly 0 arguments"}
		}
		return boolean(false), nil
	}
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method " + method.owner.name + "::" + method.name + "() cannot be called statically"}
	}
	switch strings.ToLower(method.name) {
	case "add", "sub":
		return e.dateIntervalChange(object, method, args, line)
	case "__construct":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "datetime"}, {name: "timezone"}}, line)
		if err != nil {
			return value{}, err
		}
		if len(args) > 2 {
			return value{}, &runtimeError{line, "ArgumentCountError", "DateTime::__construct() expects at most 2 arguments"}
		}
		loc := time.UTC
		if present[1] && bound[1].kind != nullKind {
			if bound[1].kind != objectKind || bound[1].o == nil || !strings.EqualFold(bound[1].o.className(), "DateTimeZone") {
				return value{}, builtinTypeError(line, "DateTime::__construct", 2, "timezone", "?DateTimeZone", bound[1])
			}
			if bound[1].o.timezone != nil {
				loc = bound[1].o.timezone
			}
		}
		now := time.Now().In(loc)
		if present[0] && bound[0].kind != nullKind && bound[0].text() != "now" && bound[0].text() != "" {
			parsed, parseErr := time.ParseInLocation("2006-01-02 15:04:05", bound[0].text(), loc)
			if parseErr != nil {
				return value{}, &runtimeError{line, "Error", "gophp: unsupported DateTime input " + bound[0].text()}
			}
			now = parsed
		}
		object.o.dateTime = &now
		return value{kind: nullKind}, nil
	case "format":
		if len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "DateTime::format() expects exactly 1 argument"}
		}
		if object.o.dateTime == nil {
			return value{}, &runtimeError{line, "Error", "DateTime object not initialized"}
		}
		formatted, err := formatDateTime(*object.o.dateTime, args[0].value.text(), line)
		return value{kind: stringKind, s: formatted}, err

	case "gettimestamp":
		if len(args) != 0 {
			return value{}, &runtimeError{line, "ArgumentCountError", "DateTime::getTimestamp() expects exactly 0 arguments"}
		}
		if object.o.dateTime == nil {
			return value{}, &runtimeError{line, "Error", "DateTime object not initialized"}
		}
		return integer(object.o.dateTime.Unix()), nil
	case "gettimezone":
		e.initClasses()
		return value{kind: objectKind, o: &phpObject{class: e.classes["datetimezone"], properties: &phpArray{}, timezone: object.o.dateTime.Location()}}, nil
	case "settimezone":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "timezone"}}, line)
		if err != nil {
			return value{}, err
		}
		if len(args) != 1 || !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", "DateTime::setTimezone() expects exactly 1 argument"}
		}
		zone := bound[0]
		if zone.kind != objectKind || zone.o == nil || zone.o.timezone == nil {
			return value{}, builtinTypeError(line, "DateTime::setTimezone", 1, "timezone", "DateTimeZone", zone)
		}
		if object.o.dateTime == nil {
			return value{}, &runtimeError{line, "Error", "DateTime object not initialized"}
		}
		updated := object.o.dateTime.In(zone.o.timezone)
		if method.owner.name == "DateTimeImmutable" {
			cloned := *object.o
			cloned.dateTime = &updated
			return value{kind: objectKind, o: &cloned}, nil
		}
		object.o.dateTime = &updated
		return object, nil
	}
	return value{}, &runtimeError{line, "Error", "gophp: unsupported DateTime method"}
}

func (e *evaluator) invokeDateTimeZone(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method DateTimeZone::__construct() cannot be called statically"}
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "timezone"}}, line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{line, "ArgumentCountError", "DateTimeZone::__construct() expects exactly 1 argument"}
	}
	name, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(line, "DateTimeZone::__construct", 1, "timezone", "string", bound[0])
	}
	location, loadErr := time.LoadLocation(name.s)
	if loadErr != nil {
		return value{}, &runtimeError{line, "DateInvalidTimeZoneException", "DateTimeZone::__construct(): Unknown or bad timezone (" + name.s + ")"}
	}
	object.o.timezone = location
	return value{kind: nullKind}, nil
}

// Format the DateTime tokens used by Carbon's Eloquent JSON serialization.
// Escaped bytes are literal; unknown alphabetic tokens stay explicit errors.
func formatDateTime(t time.Time, format string, line int) (string, error) {
	var result strings.Builder
	for i := 0; i < len(format); i++ {
		ch := format[i]
		if ch == '\\' {
			i++
			if i < len(format) {
				result.WriteByte(format[i])
			} else {
				result.WriteByte(0)
			}
			continue
		}
		layouts := map[byte]string{'Y': "2006", 'y': "06", 'm': "01", 'n': "1", 'd': "02", 'j': "2", 'H': "15", 'h': "03", 'g': "3", 'i': "04", 's': "05", 'a': "pm", 'A': "PM", 'e': "", 'T': "MST", 'O': "-0700", 'P': "-07:00", 'p': "Z07:00", 'D': "Mon", 'l': "Monday", 'M': "Jan", 'F': "January", 'c': "2006-01-02T15:04:05-07:00", 'r': "Mon, 02 Jan 2006 15:04:05 -0700"}
		if layout, ok := layouts[ch]; ok {
			if ch == 'e' {
				result.WriteString(t.Location().String())
			} else {
				result.WriteString(t.Format(layout))
			}
			continue
		}
		switch ch {
		case 'u':
			fmt.Fprintf(&result, "%06d", t.Nanosecond()/1000)
		case 'v':
			fmt.Fprintf(&result, "%03d", t.Nanosecond()/1000000)
		case 'U':
			fmt.Fprintf(&result, "%d", t.Unix())
		case 'G':
			fmt.Fprintf(&result, "%d", t.Hour())
		case 'w':
			fmt.Fprintf(&result, "%d", t.Weekday())
		case 'N':
			day := int(t.Weekday())
			if day == 0 {
				day = 7
			}
			fmt.Fprintf(&result, "%d", day)
		case 'z':
			fmt.Fprintf(&result, "%d", t.YearDay()-1)
		case 't':
			fmt.Fprintf(&result, "%d", time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day())
		case 'L':
			if t.Year()%4 == 0 && (t.Year()%100 != 0 || t.Year()%400 == 0) {
				result.WriteByte('1')
			} else {
				result.WriteByte('0')
			}
		case 'Z':
			_, offset := t.Zone()
			fmt.Fprintf(&result, "%d", offset)
		default:
			if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' {
				return "", &runtimeError{line, "Error", "gophp: unsupported DateTime::format token " + string(ch)}
			}
			result.WriteByte(ch)
		}
	}
	return result.String(), nil
}
