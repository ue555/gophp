package php

import (
	"fmt"
	"strings"
)

// The comparison CLI uses php -n on the pinned Unix PHP images. Keep this
// request-local; changing it must not affect another Run invocation.
const defaultIncludePath = ".:/usr/local/lib/php"

func includePathBuiltin(name string) bool {
	switch strings.ToLower(name) {
	case "get_include_path", "set_include_path", "ini_restore", "ini_get", "ini_set":
		return true
	}
	return false
}

func (e *evaluator) callIncludePathBuiltin(node functionCall) (result value, err error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter(nil)
	if name == "set_include_path" {
		params = []parameter{{name: "include_path"}}
	} else if name == "ini_restore" || name == "ini_get" {
		params = []parameter{{name: "option"}}
	} else if name == "ini_set" {
		params = []parameter{{name: "option"}, {name: "value"}}
	}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	defer func() {
		if err != nil && e.failureFrames == nil {
			e.captureFailure(callFrame{method: name, line: node.line, filename: e.filename, args: argumentValues(args)})
		}
	}()
	want := len(params)
	if len(args) != want || want == 1 && !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("%s() expects exactly %d argument%s, %d given", name, want, map[bool]string{true: "", false: "s"}[want == 1], len(args))}
	}
	if name == "get_include_path" {
		return value{kind: stringKind, s: e.includePath}, nil
	}
	if name == "ini_get" {
		option, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "ini_get", 1, "option", "string", bound[0])
		}
		switch option.s {
		case "include_path":
			return value{kind: stringKind, s: e.includePath}, nil
		case "memory_limit":
			return value{kind: stringKind, s: "128M"}, nil
		case "default_charset":
			return value{kind: stringKind, s: "UTF-8"}, nil
		case "date.timezone":
			return value{kind: stringKind, s: "UTC"}, nil
		case "default_mimetype":
			return value{kind: stringKind, s: e.defaultMimetype}, nil
		default:
			return boolean(false), nil
		}
	}
	if name == "ini_set" {
		option, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "ini_set", 1, "option", "string", bound[0])
		}
		input := bound[1]
		if input.kind == nullKind {
			input = value{kind: stringKind, s: ""}
		}
		converted, _, ok := e.coerceType(input, "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "ini_set", 2, "value", "string", bound[1])
		}
		switch option.s {
		case "default_mimetype":
			previous := e.defaultMimetype
			e.defaultMimetype = converted.s
			return value{kind: stringKind, s: previous}, nil
		case "include_path":
			if converted.s == "" {
				return boolean(false), nil
			}
			previous := e.includePath
			e.includePath = converted.s
			return value{kind: stringKind, s: previous}, nil
		default:
			return boolean(false), nil
		}
	}
	argument := bound[0]
	if argument.kind == nullKind {
		if e.strictTypes {
			return value{}, builtinTypeError(node.line, name, 1, params[0].name, "string", argument)
		}
		if err := e.diagnostic("Deprecated", node.line, name+"(): Passing null to parameter #1 ($"+params[0].name+") of type string is deprecated"); err != nil {
			return value{}, err
		}
		argument = value{kind: stringKind, s: ""}
	}
	converted, ok := e.coerceDeclaredType(argument, "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, name, 1, params[0].name, "string", bound[0])
	}
	if strings.IndexByte(converted.s, 0) >= 0 {
		return value{}, &runtimeError{node.line, "ValueError", name + "(): Argument #1 ($" + params[0].name + ") must not contain any null bytes"}
	}
	if name == "set_include_path" {
		if converted.s == "" {
			return value{kind: boolKind}, nil
		}
		previous := e.includePath
		e.includePath = converted.s
		return value{kind: stringKind, s: previous}, nil
	}
	if converted.s != "include_path" {
		return value{}, &runtimeError{node.line, "Error", "gophp: ini_restore() currently supports only include_path"}
	}
	e.includePath = defaultIncludePath
	return value{}, nil
}
