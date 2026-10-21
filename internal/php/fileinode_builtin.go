package php

import (
	"os"
	"syscall"
)

func (e *evaluator) callFileInode(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "filename"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "fileinode() expects exactly 1 argument"}
	}
	path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "fileinode", 1, "filename", "string", bound[0])
	}
	info, statErr := os.Stat(path.s)
	if statErr != nil {
		if err := e.warn(node.line, "fileinode(): stat failed for "+path.s); err != nil {
			return value{}, err
		}
		return boolean(false), nil
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return integer(int64(stat.Ino)), nil
	}
	return value{}, &runtimeError{node.line, "Error", "gophp: fileinode is unavailable on this platform"}
}
