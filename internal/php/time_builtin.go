package php

import (
	"fmt"
	"time"
)

func (e *evaluator) callTime(node functionCall) (value, error) {
	if len(node.args) != 0 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "time() expects exactly 0 arguments"}
	}
	return integer(time.Now().Unix()), nil
}

func (e *evaluator) callMicrotime(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "as_float"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) > 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "microtime() expects at most 1 argument"}
	}
	now := time.Now()
	if present[0] && bound[0].truthy() {
		return value{kind: floatKind, f: float64(now.Unix()) + float64(now.Nanosecond())/1e9}, nil
	}
	return value{kind: stringKind, s: fmt.Sprintf("0.%08d %d", now.Nanosecond()/10, now.Unix())}, nil
}
