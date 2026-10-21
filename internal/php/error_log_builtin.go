package php

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// error_log's CLI destinations are kept separate from PHP diagnostics. The
// application logger uses type 0; type 3 is the file-appending form.
func (e *evaluator) callErrorLog(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "message"}, {name: "message_type"}, {name: "destination"}, {name: "additional_headers"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 4 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "error_log() expects between 1 and 4 arguments"}
	}
	message, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "error_log", 1, "message", "string", bound[0])
	}
	mode := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "error_log", 2, "message_type", "int", bound[1])
		}
		mode = v.i
	}
	msg := strings.SplitN(message.s, "\x00", 2)[0]
	switch mode {
	case 0, 4:
		_, _ = fmt.Fprintln(e.stderr, msg)
		return boolean(true), nil
	case 3:
		if !present[2] || bound[2].kind == nullKind {
			return boolean(false), nil
		}
		destination, _, ok := e.coerceType(bound[2], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "error_log", 3, "destination", "?string", bound[2])
		}
		path := destination.s
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(e.filename), path)
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return boolean(false), nil
		}
		_, err = f.WriteString(msg)
		closeErr := f.Close()
		return boolean(err == nil && closeErr == nil), nil
	default:
		return boolean(false), nil
	}
}
