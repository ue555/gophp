package php

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type phpProcess struct {
	command  string
	cmd      *exec.Cmd
	done     chan struct{}
	exitCode int
	closed   bool
	pipes    []*phpStream
}

func processBuiltin(name string) bool {
	switch strings.ToLower(name) {
	case "proc_open", "proc_get_status", "proc_close", "proc_terminate":
		return true
	}
	return false
}

func (e *evaluator) callProcessBuiltin(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	if name == "proc_open" {
		return e.callProcOpen(node)
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || args[0].value.kind != resourceKind || args[0].value.stream == nil || args[0].value.stream.process == nil {
		return value{}, &runtimeError{node.line, "TypeError", name + "(): Argument #1 ($process) must be of type resource"}
	}
	p := args[0].value.stream.process
	switch name {
	case "proc_get_status":
		if len(args) != 1 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "proc_get_status() expects exactly 1 argument"}
		}
		running := true
		select {
		case <-p.done:
			running = false
		default:
		}
		out := arrayValue()
		set := func(k string, v value) { out.a = out.a.with(arrayKey{isString: true, text: k}, v) }
		set("command", value{kind: stringKind, s: p.command})
		set("pid", integer(int64(p.cmd.Process.Pid)))
		set("running", boolean(running))
		set("signaled", boolean(false))
		set("stopped", boolean(false))
		code := int64(-1)
		if !running {
			code = int64(p.exitCode)
		}
		set("exitcode", integer(code))
		set("cached", boolean(!running))
		set("termsig", integer(0))
		set("stopsig", integer(0))
		return out, nil
	case "proc_close":
		if len(args) != 1 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "proc_close() expects exactly 1 argument"}
		}
		for _, pipe := range p.pipes {
			if !pipe.closed && pipe.file != nil {
				_ = pipe.file.Close()
				pipe.closed = true
			}
		}
		<-p.done
		p.closed = true
		args[0].value.stream.closed = true
		return integer(int64(p.exitCode)), nil
	case "proc_terminate":
		if len(args) > 2 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "proc_terminate() expects 1 to 2 arguments"}
		}
		if p.closed {
			return boolean(false), nil
		}
		return boolean(p.cmd.Process.Kill() == nil), nil
	}
	return value{}, nil
}

// The output parameter is a place, not an evaluated value. In particular,
// reading an undefined $pipes here would emit a warning that PHP does not.
func (e *evaluator) callProcOpen(node functionCall) (value, error) {
	if len(node.args) < 3 || len(node.args) > 6 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("proc_open() expects 3 to 6 arguments, %d given", len(node.args))}
	}
	args := make([]value, len(node.args))
	for i, arg := range node.args {
		if _, named := arg.(namedArgument); named {
			return value{}, &runtimeError{node.line, "Error", "gophp: named proc_open() arguments are not supported"}
		}
		if i == 2 {
			continue
		}
		var err error
		args[i], err = e.eval(arg)
		if err != nil {
			return value{}, err
		}
	}
	if args[1].kind != arrayKind || args[1].a == nil {
		return value{}, builtinTypeError(node.line, "proc_open", 2, "descriptor_spec", "array", args[1])
	}
	place, err := e.referencePlace(node.args[2], false)
	if err != nil {
		return value{}, err
	}
	var cmd *exec.Cmd
	commandText := ""
	switch args[0].kind {
	case stringKind:
		commandText = args[0].s
		cmd = exec.Command("/bin/sh", "-c", commandText)
	case arrayKind:
		if args[0].a == nil || len(args[0].a.entries) == 0 {
			return value{}, &runtimeError{node.line, "ValueError", "proc_open(): Argument #1 ($command) must have at least one element"}
		}
		argv := []string{}
		for _, entry := range args[0].a.entries {
			v, _, ok := e.coerceType(dereference(entry.value), "string", false, e.strictTypes)
			if !ok {
				return value{}, &runtimeError{node.line, "TypeError", "proc_open(): Command arguments must be strings"}
			}
			argv = append(argv, v.s)
		}
		commandText = strings.Join(argv, " ")
		cmd = exec.Command(argv[0], argv[1:]...)
	default:
		return value{}, builtinTypeError(node.line, "proc_open", 1, "command", "array|string", args[0])
	}
	if len(args) > 3 && args[3].kind != nullKind {
		cwd, _, ok := e.coerceType(args[3], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "proc_open", 4, "cwd", "?string", args[3])
		}
		cmd.Dir = cwd.s
	}
	if len(args) > 4 && args[4].kind != nullKind {
		if args[4].kind != arrayKind {
			return value{}, builtinTypeError(node.line, "proc_open", 5, "env_vars", "?array", args[4])
		}
		env := append([]string(nil), os.Environ()...)
		for _, entry := range args[4].a.entries {
			if !entry.key.isString {
				continue
			}
			v, _, ok := e.coerceType(dereference(entry.value), "string", false, false)
			if ok {
				env = append(env, entry.key.text+"="+v.s)
			}
		}
		cmd.Env = env
	}
	pipes := arrayValue()
	childFiles := []*os.File{}
	parentStreams := []*phpStream{}
	closeFiles := func() {
		for _, file := range childFiles {
			_ = file.Close()
		}
		for _, stream := range parentStreams {
			_ = stream.file.Close()
		}
	}
	for _, entry := range args[1].a.entries {
		if entry.key.isString || entry.key.index < 0 || entry.key.index > 2 || entry.value.kind != arrayKind || entry.value.a == nil || len(entry.value.a.entries) < 2 {
			closeFiles()
			return value{}, &runtimeError{node.line, "Error", "gophp: proc_open() supports only pipe descriptors 0, 1 and 2"}
		}
		kind := dereference(entry.value.a.entries[0].value)
		mode := dereference(entry.value.a.entries[1].value)
		if kind.kind != stringKind || kind.s != "pipe" || mode.kind != stringKind || (mode.s != "r" && mode.s != "w") {
			closeFiles()
			return value{}, &runtimeError{node.line, "Error", "gophp: proc_open() supports only pipe descriptors"}
		}
		readEnd, writeEnd, pipeErr := os.Pipe()
		if pipeErr != nil {
			closeFiles()
			return boolean(false), nil
		}
		var child, parent *os.File
		if mode.s == "r" {
			child, parent = readEnd, writeEnd
		} else {
			child, parent = writeEnd, readEnd
		}
		switch entry.key.index {
		case 0:
			cmd.Stdin = child
		case 1:
			cmd.Stdout = child
		case 2:
			cmd.Stderr = child
		}
		childFiles = append(childFiles, child)
		stream := &phpStream{uri: "pipe", file: parent, chunkSize: 8192}
		if mode.s == "r" {
			stream.writer = parent
		}
		parentStreams = append(parentStreams, stream)
		pipes.a = pipes.a.with(arrayKey{index: entry.key.index}, value{kind: resourceKind, stream: stream})
	}
	if err := cmd.Start(); err != nil {
		closeFiles()
		if warning := e.warn(node.line, "proc_open(): "+err.Error()); warning != nil {
			return value{}, warning
		}
		return boolean(false), nil
	}
	for _, file := range childFiles {
		_ = file.Close()
	}
	p := &phpProcess{command: commandText, cmd: cmd, done: make(chan struct{}), pipes: parentStreams, exitCode: -1}
	go func() {
		err := cmd.Wait()
		if err == nil {
			p.exitCode = 0
		} else if exit, ok := err.(*exec.ExitError); ok {
			p.exitCode = exit.ExitCode()
		}
		close(p.done)
	}()
	e.processes = append(e.processes, p)
	if err := place.assign(pipes); err != nil {
		_ = cmd.Process.Kill()
		return value{}, err
	}
	return value{kind: resourceKind, stream: &phpStream{uri: "process", process: p}}, nil
}

func (e *evaluator) closeProcesses() {
	for _, p := range e.processes {
		for _, pipe := range p.pipes {
			if !pipe.closed && pipe.file != nil {
				_ = pipe.file.Close()
				pipe.closed = true
			}
		}
		if !p.closed {
			select {
			case <-p.done:
			default:
				_ = p.cmd.Process.Kill()
				<-p.done
			}
		}
	}
}
