package php

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func (e *evaluator) callFilemtime(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	b, p, _, err := bindCallArguments(args, []parameter{{name: "filename"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !p[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "filemtime() expects exactly 1 argument"}
	}
	v, _, ok := e.coerceType(b[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "filemtime", 1, "filename", "string", b[0])
	}
	if strings.IndexByte(v.s, 0) >= 0 {
		return value{}, &runtimeError{node.line, "ValueError", "filemtime(): Filename contains null byte"}
	}
	info, err := os.Stat(v.s)
	if err != nil {
		if err = e.warn(node.line, "filemtime(): stat failed for "+v.s); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	return integer(info.ModTime().Unix()), nil
}

func (e *evaluator) callIsWritable(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects exactly 1 argument"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, node.name, 1, "filename", "string", bound[0])
	}
	path, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	if path == "" || strings.IndexByte(path, 0) >= 0 {
		return boolean(false), nil
	}
	// Check OS access without opening/truncating a file. This also follows
	// symlinks and accounts for read-only mounts and directory permissions.
	return boolean(syscall.Access(path, 2) == nil), nil // W_OK
}

func (e *evaluator) callFileExists(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "file_exists() expects exactly 1 argument"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "file_exists", 1, "filename", "string", bound[0])
	}
	path, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	if path == "" || strings.IndexByte(path, 0) >= 0 {
		return boolean(false), nil
	}
	_, err = os.Stat(path)
	return boolean(err == nil), nil
}

func (e *evaluator) callIsLink(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "is_link() expects exactly 1 argument"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "is_link", 1, "filename", "string", bound[0])
	}
	path, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	if path == "" || strings.IndexByte(path, 0) >= 0 {
		return boolean(false), nil
	}
	info, statErr := os.Lstat(path)
	return boolean(statErr == nil && info.Mode()&os.ModeSymlink != 0), nil
}

func (e *evaluator) callIsDir(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "is_dir() expects exactly 1 argument"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "is_dir", 1, "filename", "string", bound[0])
	}
	path, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	if path == "" || strings.IndexByte(path, 0) >= 0 {
		return boolean(false), nil
	}
	info, statErr := os.Stat(path)
	return boolean(statErr == nil && info.IsDir()), nil
}

func (e *evaluator) callIsFile(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "is_file() expects exactly 1 argument"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "is_file", 1, "filename", "string", bound[0])
	}
	path, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	if path == "" || strings.IndexByte(path, 0) >= 0 {
		return boolean(false), nil
	}
	info, statErr := os.Stat(path)
	return boolean(statErr == nil && info.Mode().IsRegular()), nil
}

func (e *evaluator) callFunctionExists(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "function"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "function_exists() expects exactly 1 argument"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "function_exists", 1, "function", "string", bound[0])
	}
	name, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	name = strings.TrimPrefix(name, "\\")
	return boolean(e.functions[symbolKey(name)] != nil || builtinNameForVersion(name, e.version)), nil
}

func (e *evaluator) callExtensionLoaded(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "extension"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "extension_loaded() expects exactly 1 argument"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "extension_loaded", 1, "extension", "string", bound[0])
	}
	if _, err := e.toText(bound[0], node.line); err != nil {
		return value{}, err
	}
	// Optional extensions are not embedded in this interpreter. In particular,
	// returning true for mbstring would skip Composer's PHP polyfill while the
	// underlying mbstring APIs remain unavailable.
	return boolean(false), nil
}

func (e *evaluator) callClassExists(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "class"}, {name: "autoload"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", node.name + "() expects 1 or 2 arguments"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, node.name, 1, "class", "string", bound[0])
	}
	name, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	name = strings.TrimPrefix(name, "\\")
	e.initClasses()
	key := symbolKey(name)
	if e.classes[key] == nil && (!present[1] || bound[1].truthy()) {
		if err := e.tryAutoload(name, node.line); err != nil {
			return value{}, err
		}
	}
	class := e.classes[key]
	if class == nil {
		return boolean(false), nil
	}
	switch strings.ToLower(node.name) {
	case "interface_exists":
		return boolean(class.kind == "interface"), nil
	case "trait_exists":
		return boolean(class.kind == "trait"), nil
	case "enum_exists":
		return boolean(class.kind == "enum"), nil
	default:
		return boolean(class.kind != "interface" && class.kind != "trait"), nil
	}
}

func (e *evaluator) callDirname(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "path"}, {name: "levels"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "dirname() expects 1 or 2 arguments"}
	}
	if bound[0].kind == arrayKind || bound[0].kind == objectKind {
		return value{}, builtinTypeError(node.line, "dirname", 1, "path", "string", bound[0])
	}
	path, err := e.toText(bound[0], node.line)
	if err != nil {
		return value{}, err
	}
	levels := int64(1)
	if present[1] {
		levels, err = e.countMode(bound[1], node.line)
		if err != nil {
			return value{}, err
		}
	}
	if levels < 1 {
		return value{}, &runtimeError{node.line, "ValueError", "dirname(): Argument #2 ($levels) must be greater than or equal to 1"}
	}
	if levels > 256 {
		levels = 256
	}
	for i := int64(0); i < levels; i++ {
		path = filepath.Dir(path)
	}
	return value{kind: stringKind, s: path}, nil
}
