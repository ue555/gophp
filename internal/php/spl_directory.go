package php

import (
	"os"
	"path/filepath"
	"strings"
)

func (e *evaluator) invokeSPLDirectory(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method " + method.owner.name + "::" + method.name + "() cannot be called statically"}
	}
	o := object.o
	name := strings.ToLower(method.name)
	if name == "__construct" {
		if method.owner.name == "SplFileInfo" {
			bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}}, line)
			if err != nil {
				return value{}, err
			}
			if len(args) != 1 || !present[0] {
				return value{}, &runtimeError{line, "ArgumentCountError", "SplFileInfo::__construct() expects exactly 1 argument"}
			}
			path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(line, "SplFileInfo::__construct", 1, "filename", "string", bound[0])
			}
			o.filePath = path.s
			return value{kind: nullKind}, nil
		}
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "directory"}, {name: "flags"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", method.owner.name + "::__construct() expects at least 1 argument"}
		}
		path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, method.owner.name+"::__construct", 1, "directory", "string", bound[0])
		}
		flags := int64(4096)
		if present[1] {
			v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(line, method.owner.name+"::__construct", 2, "flags", "int", bound[1])
			}
			flags = v.i
		}
		if err := initDirectoryIterator(o, path.s, flags); err != nil {
			return value{}, &runtimeError{line, "UnexpectedValueException", "RecursiveDirectoryIterator::__construct(" + path.s + "): Failed to open directory"}
		}
		return value{kind: nullKind}, nil
	}
	currentPath := o.filePath
	if o.directoryPath != "" {
		if o.directoryIndex >= 0 && o.directoryIndex < len(o.directoryEntries) {
			currentPath = filepath.Join(o.directoryPath, o.directoryEntries[o.directoryIndex])
		} else {
			currentPath = o.directoryPath
		}
	}
	switch name {
	case "getpathname", "__tostring":
		return value{kind: stringKind, s: currentPath}, nil
	case "getfilename":
		return value{kind: stringKind, s: filepath.Base(currentPath)}, nil
	case "getextension":
		base := filepath.Base(currentPath)
		dot := strings.LastIndexByte(base, '.')
		if dot < 0 {
			return value{kind: stringKind}, nil
		}
		return value{kind: stringKind, s: base[dot+1:]}, nil
	case "getbasename":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "suffix"}}, line)
		if err != nil {
			return value{}, err
		}
		base := filepath.Base(currentPath)
		if present[0] {
			suffix, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(line, "SplFileInfo::getBasename", 1, "suffix", "string", bound[0])
			}
			if suffix.s != "" {
				base = strings.TrimSuffix(base, suffix.s)
			}
		}
		return value{kind: stringKind, s: base}, nil
	case "getpath":
		return value{kind: stringKind, s: filepath.Dir(currentPath)}, nil
	case "getrealpath":
		path, err := filepath.EvalSymlinks(currentPath)
		if err != nil {
			return boolean(false), nil
		}
		path, err = filepath.Abs(path)
		if err != nil {
			return boolean(false), nil
		}
		return value{kind: stringKind, s: path}, nil
	case "isfile", "isdir", "islink":
		info, err := os.Lstat(currentPath)
		if err != nil {
			return boolean(false), nil
		}
		if name == "islink" {
			return boolean(info.Mode()&os.ModeSymlink != 0), nil
		}
		if name == "isfile" {
			return boolean(info.Mode().IsRegular()), nil
		}
		if info.IsDir() {
			return boolean(true), nil
		}
		if info.Mode()&os.ModeSymlink != 0 && o.directoryFlags&512 != 0 {
			info, err = os.Stat(currentPath)
			return boolean(err == nil && info.IsDir()), nil
		}
		return boolean(false), nil
	case "rewind":
		o.directoryIndex = 0
		return value{kind: nullKind}, nil
	case "next":
		o.directoryIndex++
		return value{kind: nullKind}, nil
	case "valid":
		return boolean(o.directoryIndex >= 0 && o.directoryIndex < len(o.directoryEntries)), nil
	case "key":
		if o.directoryIndex < 0 || o.directoryIndex >= len(o.directoryEntries) {
			return value{kind: nullKind}, nil
		}
		if o.directoryFlags&256 != 0 {
			return value{kind: stringKind, s: o.directoryEntries[o.directoryIndex]}, nil
		}
		return value{kind: stringKind, s: currentPath}, nil
	case "current":
		if o.directoryIndex < 0 || o.directoryIndex >= len(o.directoryEntries) {
			return value{kind: nullKind}, nil
		}
		if o.directoryFlags&32 != 0 {
			return value{kind: stringKind, s: currentPath}, nil
		}
		if o.directoryFlags&16 != 0 {
			return object, nil
		}
		e.initClasses()
		return value{kind: objectKind, o: &phpObject{class: e.classes["splfileinfo"], properties: &phpArray{}, filePath: currentPath}}, nil
	case "getflags":
		return integer(o.directoryFlags), nil
	case "setflags":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "flags"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", "FilesystemIterator::setFlags() expects exactly 1 argument"}
		}
		v, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "FilesystemIterator::setFlags", 1, "flags", "int", bound[0])
		}
		o.directoryFlags = v.i
		return value{kind: nullKind}, nil
	case "getsubpath":
		return value{kind: stringKind, s: o.directorySubPath}, nil
	case "getsubpathname":
		if o.directoryIndex < 0 || o.directoryIndex >= len(o.directoryEntries) {
			return value{kind: stringKind}, nil
		}
		name := o.directoryEntries[o.directoryIndex]
		if o.directorySubPath != "" {
			name = filepath.Join(o.directorySubPath, name)
		}
		return value{kind: stringKind, s: name}, nil
	case "haschildren":
		if o.directoryIndex < 0 || o.directoryIndex >= len(o.directoryEntries) {
			return boolean(false), nil
		}
		info, err := os.Lstat(currentPath)
		if err != nil {
			return boolean(false), nil
		}
		if info.Mode()&os.ModeSymlink != 0 && o.directoryFlags&512 == 0 {
			return boolean(false), nil
		}
		info, err = os.Stat(currentPath)
		return boolean(err == nil && info.IsDir()), nil
	case "getchildren":
		if o.directoryIndex < 0 || o.directoryIndex >= len(o.directoryEntries) {
			return value{kind: nullKind}, nil
		}
		// SPL creates children as instances of the concrete iterator class.
		// Constructing a bare phpObject loses subclass property defaults and
		// leaves typed properties uninitialized (Symfony Finder relies on them).
		childValue, err := e.newInstance(newObject{name: o.class.name, line: line, args: []expr{
			literal{value: value{kind: stringKind, s: currentPath}, line: line},
			literal{value: integer(o.directoryFlags), line: line},
		}})
		if err != nil {
			return value{}, err
		}
		child := childValue.o
		child.directorySubPath = o.directoryEntries[o.directoryIndex]
		if o.directorySubPath != "" {
			child.directorySubPath = filepath.Join(o.directorySubPath, child.directorySubPath)
		}
		return childValue, nil
	}
	return value{}, &runtimeError{line, "Error", "gophp: Spl directory method is not implemented"}
}

func initDirectoryIterator(o *phpObject, path string, flags int64) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return err
	}
	o.directoryPath = path
	o.filePath = path
	o.directoryEntries = names
	o.directoryIndex = 0
	o.directoryFlags = flags
	return nil
}
