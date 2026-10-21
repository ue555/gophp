package php

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var manifestUmask = readProcessUmask()

// Read the process mask once; changing it is deliberately unsupported because
// the HTTP interpreter shares a process across independent PHP requests.
func readProcessUmask() int {
	mask := syscall.Umask(0)
	syscall.Umask(mask)
	return mask
}

func manifestFilesystemBuiltin(name string) bool {
	switch strings.ToLower(name) {
	case "tempnam", "sys_get_temp_dir", "chmod", "umask", "rename", "unlink":
		return true
	}
	return false
}

func (e *evaluator) callManifestFilesystem(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{}
	required := 0
	switch name {
	case "tempnam":
		params = []parameter{{name: "directory"}, {name: "prefix"}}
		required = 2
	case "chmod":
		params = []parameter{{name: "filename"}, {name: "permissions"}}
		required = 2
	case "umask":
		params = []parameter{{name: "mask"}}
	case "rename":
		params = []parameter{{name: "from"}, {name: "to"}, {name: "context"}}
		required = 2
	case "unlink":
		params = []parameter{{name: "filename"}, {name: "context"}}
		required = 1
	}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < required || len(args) > len(params) {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "(): invalid number of arguments"}
	}
	for i := 0; i < required; i++ {
		if !present[i] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", name + "(): missing required argument"}
		}
	}
	if name == "sys_get_temp_dir" {
		return value{kind: stringKind, s: os.TempDir()}, nil
	}
	if name == "umask" {
		if present[0] && bound[0].kind != nullKind {
			return value{}, &runtimeError{node.line, "Error", "gophp: setting the process umask is not implemented"}
		}
		return integer(int64(manifestUmask)), nil
	}
	paths := make([]string, required)
	for i := 0; i < required; i++ {
		if name == "chmod" && i == 1 {
			continue
		}
		v, _, ok := e.coerceType(bound[i], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, name, i+1, params[i].name, "string", bound[i])
		}
		if strings.IndexByte(v.s, 0) >= 0 {
			return value{}, &runtimeError{node.line, "ValueError", name + "(): path must not contain any null bytes"}
		}
		paths[i] = v.s
	}
	if len(params) > required && present[required] && bound[required].kind != nullKind {
		return value{}, &runtimeError{node.line, "Error", "gophp: filesystem stream context is not implemented"}
	}
	switch name {
	case "tempnam":
		prefix := filepath.Base(paths[1])
		if paths[1] == "" {
			prefix = ""
		}
		if len(prefix) > 63 {
			prefix = prefix[:63]
		}
		// CreateTemp gives each call a distinct file and mode 0600. A non-writable
		// directory falls back to the system temporary directory, with PHP's notice.
		dir := paths[0]
		info, statErr := os.Stat(dir)
		if statErr != nil || !info.IsDir() || syscall.Access(dir, 2) != nil {
			if err := e.diagnostic("Notice", node.line, "tempnam(): file created in the system's temporary directory"); err != nil {
				return value{}, err
			}
			dir = os.TempDir()
		}
		dir, err = filepath.Abs(dir)
		if err != nil {
			break
		}
		var file *os.File
		file, err = os.CreateTemp(dir, prefix+"*")
		if err == nil {
			path := file.Name()
			err = file.Close()
			if err == nil {
				return value{kind: stringKind, s: path}, nil
			}
		}
	case "chmod":
		mode, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, name, 2, "permissions", "int", bound[1])
		}
		permissions := os.FileMode(mode.i & 0777)
		if mode.i&04000 != 0 {
			permissions |= os.ModeSetuid
		}
		if mode.i&02000 != 0 {
			permissions |= os.ModeSetgid
		}
		if mode.i&01000 != 0 {
			permissions |= os.ModeSticky
		}
		err = os.Chmod(paths[0], permissions)
	case "rename":
		err = os.Rename(paths[0], paths[1])
	case "unlink":
		err = syscall.Unlink(paths[0])
	}
	if err != nil {
		if warnErr := e.warn(node.line, name+"(): "+err.Error()); warnErr != nil {
			return value{}, warnErr
		}
		return boolean(false), nil
	}
	return boolean(true), nil
}
