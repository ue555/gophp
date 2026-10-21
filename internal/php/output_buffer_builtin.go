package php

import (
	"bytes"
	"io"
	"strings"
)

type outputBuffer struct {
	buffer   *bytes.Buffer
	previous io.Writer
}

func outputBufferBuiltin(name string) bool {
	switch strings.ToLower(name) {
	case "ob_start", "ob_get_contents", "ob_get_clean", "ob_end_clean", "ob_end_flush", "ob_flush", "ob_get_level":
		return true
	}
	return false
}

func (e *evaluator) callOutputBufferBuiltin(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	if len(node.args) != 0 {
		return value{}, &runtimeError{node.line, "Error", "gophp: output buffer callbacks and arguments are not supported"}
	}
	if name == "ob_get_level" {
		return integer(int64(len(e.outputBuffers))), nil
	}
	if name == "ob_start" {
		buffer := &bytes.Buffer{}
		e.outputBuffers = append(e.outputBuffers, outputBuffer{buffer: buffer, previous: e.stdout})
		e.stdout = buffer
		return boolean(true), nil
	}
	if len(e.outputBuffers) == 0 {
		if name == "ob_get_contents" || name == "ob_get_clean" {
			return boolean(false), nil
		}
		return boolean(false), nil
	}
	current := e.outputBuffers[len(e.outputBuffers)-1]
	switch name {
	case "ob_get_contents":
		return value{kind: stringKind, s: current.buffer.String()}, nil
	case "ob_flush":
		_, _ = current.previous.Write(current.buffer.Bytes())
		current.buffer.Reset()
		return boolean(true), nil
	case "ob_get_clean", "ob_end_clean", "ob_end_flush":
		e.stdout = current.previous
		e.outputBuffers = e.outputBuffers[:len(e.outputBuffers)-1]
		if name == "ob_end_flush" {
			_, _ = current.previous.Write(current.buffer.Bytes())
		}
		if name == "ob_get_clean" {
			return value{kind: stringKind, s: current.buffer.String()}, nil
		}
		return boolean(true), nil
	}
	return value{}, nil
}

func (e *evaluator) flushOutputBuffers() {
	for len(e.outputBuffers) > 0 {
		current := e.outputBuffers[len(e.outputBuffers)-1]
		e.outputBuffers = e.outputBuffers[:len(e.outputBuffers)-1]
		_, _ = current.previous.Write(current.buffer.Bytes())
		e.stdout = current.previous
	}
}
