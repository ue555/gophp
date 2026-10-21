package php

import (
	"fmt"
	"strings"
)

const (
	phpEWarning        = 2
	phpENotice         = 8
	phpEUserError      = 256
	phpEUserWarning    = 512
	phpEUserNotice     = 1024
	phpEDeprecated     = 8192
	phpEUserDeprecated = 16384
)

type errorHandlerState struct {
	callback value
	mask     int64
}

func allErrorLevels(version Version) int64 {
	if version == PHP83 {
		return 32767
	}
	return 30719
}

func errorLevelName(level int64) string {
	switch level {
	case phpENotice, phpEUserNotice:
		return "Notice"
	case phpEDeprecated, phpEUserDeprecated:
		return "Deprecated"
	default:
		return "Warning"
	}
}

func diagnosticLevel(name string) int64 {
	switch name {
	case "Notice":
		return phpENotice
	case "Deprecated":
		return phpEDeprecated
	default:
		return phpEWarning
	}
}

func (e *evaluator) handleDiagnostic(name string, line int, message string) error {
	return e.emitError(diagnosticLevel(name), name, line, message)
}

func (e *evaluator) emitError(level int64, name string, line int, message string) error {
	filename := e.filename
	last := arrayValue()
	last.a = last.a.with(arrayKey{isString: true, text: "type"}, integer(level))
	last.a = last.a.with(arrayKey{isString: true, text: "message"}, value{kind: stringKind, s: message})
	last.a = last.a.with(arrayKey{isString: true, text: "file"}, value{kind: stringKind, s: filename})
	last.a = last.a.with(arrayKey{isString: true, text: "line"}, integer(int64(line)))
	e.lastError = last
	if e.errorHandler.kind != nullKind && e.errorHandlerMask&level != 0 {
		result, err := e.callRegisteredErrorHandler(level, message, filename, line)
		if err != nil {
			return err
		}
		if result.kind != boolKind || result.i != 0 {
			return nil
		}
	}
	if e.suppressedDiagnostics > 0 || e.errorReporting&level == 0 {
		return nil
	}
	fmt.Fprintf(e.stderr, "%s: %s in %s on line %d\n", name, message, filename, line)
	return nil
}

// PHP suspends the currently executing handler while it calls user code. A
// warning raised inside that code uses the default handler unless user code
// registers a new one. A registration must also observe no current handler.
// Restore the suspended handler if no handler remains when the call ends,
// including after an explicit set_error_handler(null).
func (e *evaluator) callRegisteredErrorHandler(level int64, message, filename string, line int) (value, error) {
	active := errorHandlerState{e.errorHandler, e.errorHandlerMask}
	e.errorHandler = value{}
	e.errorHandlerMask = 0
	result, err := e.callErrorHandler(active.callback, []value{integer(level), {kind: stringKind, s: message}, {kind: stringKind, s: filename}, integer(int64(line))}, line)
	if e.errorHandler.kind == nullKind {
		e.errorHandler, e.errorHandlerMask = active.callback, active.mask
	}
	return result, err
}

func (e *evaluator) validErrorCallback(v value) bool {
	_, ok := e.resolveCallable(v, 0)
	return ok
}

func (e *evaluator) callErrorHandler(callback value, values []value, line int) (value, error) {
	args := make([]callArgument, len(values))
	for i, v := range values {
		args[i] = callArgument{value: v, line: line}
	}
	target, ok := e.resolveCallable(callback, line)
	if !ok {
		return value{}, &runtimeError{line, "TypeError", "Invalid error handler callback"}
	}
	return e.invokeResolved(target, args, line)
}

func (e *evaluator) callErrorBuiltin(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	fail := func(err error) (value, error) {
		if e.failureFrames == nil {
			e.captureFailure(callFrame{method: name, line: node.line, filename: e.filename, args: argumentValues(args)})
		}
		return value{}, err
	}
	params := []parameter{{name: "callback"}, {name: "error_levels"}}
	if name == "restore_error_handler" {
		params = nil
	} else if name == "trigger_error" {
		params = []parameter{{name: "message"}, {name: "error_level"}}
	}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return fail(err)
	}
	if name == "restore_error_handler" {
		if len(args) > 0 {
			return fail(&runtimeError{node.line, "ArgumentCountError", "restore_error_handler() expects exactly 0 arguments, " + fmt.Sprint(len(args)) + " given"})
		}
		if len(e.errorHandlerStack) > 0 {
			previous := e.errorHandlerStack[len(e.errorHandlerStack)-1]
			e.errorHandlerStack = e.errorHandlerStack[:len(e.errorHandlerStack)-1]
			e.errorHandler, e.errorHandlerMask = previous.callback, previous.mask
		}
		return boolean(true), nil
	}
	if !present[0] || len(args) > 2 {
		return fail(&runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("%s() expects %s %d arguments, %d given", name, map[bool]string{true: "at most", false: "at least"}[len(args) > 2], map[bool]int{true: 2, false: 1}[len(args) > 2], len(args))})
	}
	if name == "set_error_handler" {
		callback := bound[0]
		if callback.kind != nullKind && !e.validErrorCallback(callback) {
			return fail(&runtimeError{node.line, "TypeError", "set_error_handler(): Argument #1 ($callback) must be a valid callback or null, function \"" + callback.text() + "\" not found or invalid function name"})
		}
		mask := allErrorLevels(e.version)
		if present[1] {
			coerced, ok := e.coerceDeclaredType(bound[1], "int", false, e.strictTypes)
			if !ok {
				return fail(builtinTypeError(node.line, name, 2, "error_levels", "int", bound[1]))
			}
			mask = coerced.i
		}
		previous := e.errorHandler
		e.errorHandlerStack = append(e.errorHandlerStack, errorHandlerState{previous, e.errorHandlerMask})
		e.errorHandler, e.errorHandlerMask = callback, mask
		return previous, nil
	}
	message, ok := e.coerceDeclaredType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return fail(builtinTypeError(node.line, name, 1, "message", "string", bound[0]))
	}
	level := int64(phpEUserNotice)
	if present[1] {
		coerced, ok := e.coerceDeclaredType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return fail(builtinTypeError(node.line, name, 2, "error_level", "int", bound[1]))
		}
		level = coerced.i
	}
	if level != phpEUserError && level != phpEUserWarning && level != phpEUserNotice && level != phpEUserDeprecated {
		return fail(&runtimeError{node.line, "ValueError", "trigger_error(): Argument #2 ($error_level) must be one of E_USER_ERROR, E_USER_WARNING, E_USER_NOTICE, or E_USER_DEPRECATED"})
	}
	if level == phpEUserError && e.version != PHP83 {
		if err := e.emitError(phpEDeprecated, "Deprecated", node.line, "Passing E_USER_ERROR to trigger_error() is deprecated since 8.4, throw an exception or call exit with a string message instead"); err != nil {
			return value{}, err
		}
	}
	if level == phpEUserError {
		if e.errorHandler.kind != nullKind && e.errorHandlerMask&level != 0 {
			result, err := e.callRegisteredErrorHandler(level, message.s, e.filename, node.line)
			if err != nil {
				return value{}, err
			}
			if result.kind != boolKind || result.i != 0 {
				return boolean(true), nil
			}
		}
		if e.version == PHP85 && e.failureFrames == nil {
			e.captureFailure(callFrame{method: name, line: node.line, filename: e.filename, args: argumentValues(args)})
		}
		return value{}, &fatalDiagnostic{node.line, message.s}
	}
	if err := e.emitError(level, errorLevelName(level), node.line, message.s); err != nil {
		return value{}, err
	}
	return boolean(true), nil
}
