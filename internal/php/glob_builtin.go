package php

import (
	"os"
	"path/filepath"
	"strings"
)

func (e *evaluator) callGlob(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "pattern"}, {name: "flags"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "glob() expects 1 to 2 arguments"}
	}
	pattern, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "glob", 1, "pattern", "string", bound[0])
	}
	flags := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "glob", 2, "flags", "int", bound[1])
		}
		flags = v.i
	}
	path := pattern.s
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(e.filename), path)
	}
	matches, err := filepath.Glob(path)
	if err != nil {
		return boolean(false), nil
	}
	out := arrayValue()
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if flags&8192 != 0 && !info.IsDir() {
			continue
		}
		if flags&2 != 0 && info.IsDir() {
			match += string(os.PathSeparator)
		}
		if !filepath.IsAbs(pattern.s) {
			match = strings.TrimPrefix(match, filepath.Dir(e.filename)+string(os.PathSeparator))
		}
		out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, value{kind: stringKind, s: match})
	}
	return out, nil
}
