package php

import (
	"io"
	"os"
	"strings"
	"syscall"
)

type phpStream struct {
	uri       string
	writer    io.Writer
	file      *os.File
	closed    bool
	chunkSize int64
	process   *phpProcess
	eof       bool
	tempPath  string
}

func streamBuiltin(name string) bool {
	switch strings.ToLower(name) {
	case "fopen", "fwrite", "fflush", "fclose", "flock", "stream_get_contents", "get_resource_type", "stream_isatty", "stream_get_meta_data", "stream_set_chunk_size", "fseek", "fread", "feof", "stream_set_blocking":
		return true
	}
	return false
}

func (e *evaluator) callStreamBuiltin(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "stream"}}
	switch name {
	case "fopen":
		params = []parameter{{name: "filename"}, {name: "mode"}, {name: "use_include_path"}, {name: "context"}}
	case "fwrite":
		params = []parameter{{name: "stream"}, {name: "data"}, {name: "length"}}
	case "stream_set_chunk_size":
		params = []parameter{{name: "stream"}, {name: "size"}}
	case "flock":
		params = []parameter{{name: "stream"}, {name: "operation"}, {name: "would_block"}}
	case "stream_get_contents":
		params = []parameter{{name: "stream"}, {name: "length"}, {name: "offset"}}
	case "fseek":
		params = []parameter{{name: "stream"}, {name: "offset"}, {name: "whence"}}
	case "fread":
		params = []parameter{{name: "stream"}, {name: "length"}}
	case "stream_set_blocking":
		params = []parameter{{name: "stream"}, {name: "enable"}}
	}
	bound, present, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if name == "fopen" {
		if len(args) < 2 || len(args) > 4 || !present[0] || !present[1] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "fopen() expects 2 to 4 arguments"}
		}
		path, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "fopen", 1, "filename", "string", bound[0])
		}
		mode, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "fopen", 2, "mode", "string", bound[1])
		}
		stream := &phpStream{uri: path.s, chunkSize: 8192}
		switch path.s {
		case "php://stdout", "php://output":
			stream.writer = e.stdout
		case "php://stderr":
			stream.writer = e.stderr
		case "php://stdin":
			stream.file = os.Stdin
		case "php://temp":
			fallthrough
		default:
			if strings.HasPrefix(path.s, "php://temp/") || path.s == "php://temp" {
				file, openErr := os.CreateTemp("", "gophp-stream-*")
				if openErr != nil {
					return boolean(false), nil
				}
				stream.file, stream.writer, stream.tempPath = file, file, file.Name()
				e.openStreams = append(e.openStreams, stream)
				break
			}
			flags := os.O_RDONLY
			if strings.ContainsAny(mode.s, "wa+") {
				flags = os.O_CREATE | os.O_WRONLY
				if strings.Contains(mode.s, "a") {
					flags |= os.O_APPEND
				} else {
					flags |= os.O_TRUNC
				}
			}
			file, openErr := os.OpenFile(path.s, flags, 0666)
			if openErr != nil {
				if warnErr := e.warn(node.line, "fopen("+path.s+"): Failed to open stream: "+openErr.Error()); warnErr != nil {
					return value{}, warnErr
				}
				return boolean(false), nil
			}
			stream.file = file
			stream.writer = file
		}
		return value{kind: resourceKind, stream: stream}, nil
	}
	if len(args) < 1 || len(args) > len(params) || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects a stream argument"}
	}
	streamValue := bound[0]
	if streamValue.kind != resourceKind || streamValue.stream == nil || streamValue.stream.closed {
		return value{}, builtinTypeError(node.line, name, 1, "stream", "resource", streamValue)
	}
	stream := streamValue.stream
	switch name {
	case "fseek":
		if len(args) < 2 || len(args) > 3 || !present[1] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "fseek() expects 2 to 3 arguments"}
		}
		offset, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "fseek", 2, "offset", "int", bound[1])
		}
		whence := int64(io.SeekStart)
		if present[2] {
			whence = bound[2].i
		}
		if stream.file == nil {
			return integer(-1), nil
		}
		if _, err := stream.file.Seek(offset.i, int(whence)); err != nil {
			return integer(-1), nil
		}
		stream.eof = false
		return integer(0), nil
	case "fread":
		if len(args) != 2 || !present[1] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "fread() expects exactly 2 arguments"}
		}
		length, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "fread", 2, "length", "int", bound[1])
		}
		if length.i <= 0 {
			return value{}, &runtimeError{node.line, "ValueError", "fread(): Argument #2 ($length) must be greater than 0"}
		}
		if stream.file == nil {
			return boolean(false), nil
		}
		if length.i > 8192 {
			length.i = 8192
		}
		buffer := make([]byte, length.i)
		n, readErr := stream.file.Read(buffer)
		if readErr == io.EOF {
			stream.eof = true
		}
		if readErr != nil && readErr != io.EOF && readErr != syscall.EAGAIN {
			return boolean(false), nil
		}
		return value{kind: stringKind, s: string(buffer[:n])}, nil
	case "feof":
		if len(args) != 1 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "feof() expects exactly 1 argument"}
		}
		return boolean(stream.eof), nil
	case "stream_set_blocking":
		if len(args) != 2 || !present[1] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "stream_set_blocking() expects exactly 2 arguments"}
		}
		if stream.file == nil {
			return boolean(false), nil
		}
		return boolean(syscall.SetNonblock(int(stream.file.Fd()), !bound[1].truthy()) == nil), nil
	case "flock":
		if !present[1] || stream.file == nil {
			return boolean(false), nil
		}
		operation, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "flock", 2, "operation", "int", bound[1])
		}
		operationCode := int(operation.i)
		if operationCode&3 == 3 {
			operationCode = (operationCode &^ 3) | syscall.LOCK_UN
		}
		return boolean(syscall.Flock(int(stream.file.Fd()), operationCode) == nil), nil
	case "stream_get_contents":
		if stream.file == nil {
			return boolean(false), nil
		}
		if present[2] {
			offset, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "stream_get_contents", 3, "offset", "int", bound[2])
			}
			if _, err := stream.file.Seek(offset.i, io.SeekStart); err != nil {
				return boolean(false), nil
			}
		}
		reader := io.Reader(stream.file)
		if present[1] {
			length, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "stream_get_contents", 2, "length", "?int", bound[1])
			}
			if length.i >= 0 {
				reader = io.LimitReader(reader, length.i)
			}
		}
		data, err := io.ReadAll(reader)
		if err != nil {
			return boolean(false), nil
		}
		if !present[1] || bound[1].i < 0 {
			stream.eof = true
		}
		return value{kind: stringKind, s: string(data)}, nil
	case "stream_set_chunk_size":
		if len(args) != 2 || !present[1] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "stream_set_chunk_size() expects exactly 2 arguments"}
		}
		size, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, name, 2, "size", "int", bound[1])
		}
		if size.i < 1 {
			return value{}, &runtimeError{node.line, "ValueError", "stream_set_chunk_size(): Argument #2 ($size) must be greater than 0"}
		}
		previous := stream.chunkSize
		stream.chunkSize = size.i
		return integer(previous), nil
	case "get_resource_type":
		if stream.process != nil {
			return value{kind: stringKind, s: "process"}, nil
		}
		return value{kind: stringKind, s: "stream"}, nil
	case "stream_isatty":
		return boolean(false), nil
	case "stream_get_meta_data":
		out := arrayValue()
		out.a = out.a.with(arrayKey{isString: true, text: "uri"}, value{kind: stringKind, s: stream.uri})
		out.a = out.a.with(arrayKey{isString: true, text: "blocked"}, boolean(true))
		out.a = out.a.with(arrayKey{isString: true, text: "seekable"}, boolean(stream.file != nil && stream.file != os.Stdin))
		return out, nil
	case "fflush":
		return boolean(true), nil
	case "fclose":
		stream.closed = true
		if stream.file != nil && stream.file != os.Stdin {
			_ = stream.file.Close()
		}
		if stream.tempPath != "" {
			_ = os.Remove(stream.tempPath)
		}
		return boolean(true), nil
	case "fwrite":
		if len(args) < 2 || !present[1] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "fwrite() expects at least 2 arguments"}
		}
		data, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "fwrite", 2, "data", "string", bound[1])
		}
		if stream.writer == nil {
			return boolean(false), nil
		}
		bytes := []byte(data.s)
		if present[2] {
			length, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "fwrite", 3, "length", "?int", bound[2])
			}
			if length.i >= 0 && length.i < int64(len(bytes)) {
				bytes = bytes[:length.i]
			}
		}
		n, writeErr := stream.writer.Write(bytes)
		if writeErr != nil {
			if warnErr := e.warn(node.line, "fwrite(): "+writeErr.Error()); warnErr != nil {
				return value{}, warnErr
			}
			return boolean(false), nil
		}
		return integer(int64(n)), nil
	}
	return value{}, nil
}

func (e *evaluator) closeOpenStreams() {
	for _, stream := range e.openStreams {
		if !stream.closed && stream.file != nil {
			_ = stream.file.Close()
			stream.closed = true
		}
		if stream.tempPath != "" {
			_ = os.Remove(stream.tempPath)
		}
	}
}
