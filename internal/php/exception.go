package php

import (
	"errors"
	"fmt"
	"strings"
)

type exceptionData struct {
	message  string
	code     int64
	codeText string // PDO errors use SQLSTATE strings rather than integer codes.
	severity int64
	file     string
	line     int
	frames   []callFrame
	previous value
	captured bool
}

type thrown struct{ object value }

func (t *thrown) Error() string { return t.object.o.className() }

func isThrowableClass(class *classDecl) bool {
	for current := class; current != nil; current = current.parent {
		if equalName(current.name, "Exception") || equalName(current.name, "Error") {
			return true
		}
	}
	return false
}

func (e *evaluator) initializeException(object value, expressions []expr, line int) error {
	args, err := e.evalArguments(expressions)
	if err != nil {
		return err
	}
	object.o.exception.frames = e.snapshotCallFrames()
	return e.initializeExceptionArguments(object, args, line)
}

func (e *evaluator) initializeExceptionArguments(object value, args []callArgument, line int) error {
	if equalName(object.o.className(), "ErrorException") {
		params := []parameter{{name: "message"}, {name: "code"}, {name: "severity"}, {name: "filename"}, {name: "line"}, {name: "previous"}}
		bound, present, _, err := bindCallArguments(args, params, line)
		if err != nil {
			return err
		}
		if len(args) > 6 {
			return &runtimeError{line, "ArgumentCountError", fmt.Sprintf("ErrorException::__construct() expects at most 6 arguments, %d given", len(args))}
		}
		if present[0] {
			v, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
			if !ok {
				return builtinTypeError(line, "ErrorException::__construct", 1, "message", "string", bound[0])
			}
			object.o.exception.message = v.s
		}
		if present[1] {
			v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
			if !ok {
				return builtinTypeError(line, "ErrorException::__construct", 2, "code", "int", bound[1])
			}
			object.o.exception.code = v.i
		}
		object.o.exception.severity = 1
		if present[2] {
			v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
			if !ok {
				return builtinTypeError(line, "ErrorException::__construct", 3, "severity", "int", bound[2])
			}
			object.o.exception.severity = v.i
		}
		if present[3] && bound[3].kind != nullKind {
			v, _, ok := e.coerceType(bound[3], "string", false, e.strictTypes)
			if !ok {
				return builtinTypeError(line, "ErrorException::__construct", 4, "filename", "?string", bound[3])
			}
			object.o.exception.file = v.s
		}
		if present[4] && bound[4].kind != nullKind {
			v, _, ok := e.coerceType(bound[4], "int", false, e.strictTypes)
			if !ok {
				return builtinTypeError(line, "ErrorException::__construct", 5, "line", "?int", bound[4])
			}
			object.o.exception.line = int(v.i)
		}
		if present[5] {
			if bound[5].kind != nullKind && (bound[5].kind != objectKind || bound[5].o.exception == nil) {
				return builtinTypeError(line, "ErrorException::__construct", 6, "previous", "?Throwable", bound[5])
			}
			object.o.exception.previous = bound[5]
		}
		return nil
	}
	params := []parameter{{name: "message"}, {name: "code"}, {name: "previous"}}
	bound, present, _, err := bindCallArguments(args, params, line)
	if err != nil {
		return err
	}
	if len(args) > len(params) {
		return &runtimeError{line, "ArgumentCountError", fmt.Sprintf("%s::__construct() expects at most 3 arguments, %d given", object.o.className(), len(args))}
	}
	if present[0] {
		coerced, ok := e.coerceDeclaredType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return &runtimeError{line, "TypeError", object.o.className() + "::__construct(): Argument #1 ($message) must be of type string, " + bound[0].diagnosticName() + " given"}
		}
		object.o.exception.message = coerced.s
	}
	if present[1] {
		coerced, ok := e.coerceDeclaredType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return &runtimeError{line, "TypeError", object.o.className() + "::__construct(): Argument #2 ($code) must be of type int, " + bound[1].diagnosticName() + " given"}
		}
		if err := e.warnLossyIntTypeCoercion(bound[1], "int", e.strictTypes, line); err != nil {
			return err
		}
		object.o.exception.code = coerced.i
	}
	if present[2] && bound[2].kind != nullKind && (bound[2].kind != objectKind || bound[2].o.exception == nil) {
		return &runtimeError{line, "TypeError", object.o.className() + "::__construct(): Argument #3 ($previous) must be of type ?Throwable, " + bound[2].diagnosticName() + " given"}
	}
	if present[2] {
		object.o.exception.previous = bound[2]
	}
	return nil
}

func (e *evaluator) callExceptionMethod(object value, node methodCall) (value, bool, error) {
	return e.callExceptionAccessor(object, node.name, node.args, node.line)
}

func (e *evaluator) callExceptionAccessor(object value, methodName string, expressions []expr, line int) (value, bool, error) {
	name := strings.ToLower(methodName)
	if name != "getmessage" && name != "getcode" && name != "getfile" && name != "getline" && name != "gettraceasstring" && name != "gettrace" && name != "getprevious" && name != "getseverity" {
		return value{}, false, nil
	}
	args, err := e.evalArguments(expressions)
	if err != nil {
		return value{}, true, err
	}
	if len(args) > 0 {
		// PHP's built-in accessors ignore extra positional arguments.
		if _, _, _, err := bindCallArguments(args, nil, line); err != nil {
			return value{}, true, err
		}
	}
	data := object.o.exception
	switch name {
	case "getmessage":
		return value{kind: stringKind, s: data.message}, true, nil
	case "getcode":
		if data.codeText != "" {
			return value{kind: stringKind, s: data.codeText}, true, nil
		}
		return integer(data.code), true, nil
	case "getseverity":
		return integer(data.severity), true, nil
	case "getfile":
		return value{kind: stringKind, s: data.file}, true, nil
	case "getline":
		return integer(int64(data.line)), true, nil
	case "getprevious":
		return data.previous, true, nil
	case "gettrace":
		return exceptionTraceArray(data.frames), true, nil
	default:
		return value{kind: stringKind, s: strings.TrimSuffix(formatStackTrace(data.frames), "\n")}, true, nil
	}
}

func exceptionTraceArray(frames []callFrame) value {
	result := arrayValue()
	for i := len(frames) - 1; i >= 0; i-- {
		frame := frames[i]
		entry := arrayValue()
		add := func(key string, v value) { entry.a = entry.a.with(arrayKey{isString: true, text: key}, v) }
		add("file", value{kind: stringKind, s: frame.filename})
		add("line", integer(int64(frame.line)))
		add("function", value{kind: stringKind, s: frame.method})
		if frame.class != "" {
			add("class", value{kind: stringKind, s: frame.class})
			separator := "->"
			if frame.static {
				separator = "::"
			}
			add("type", value{kind: stringKind, s: separator})
		}
		args := arrayValue()
		for n, arg := range frame.args {
			key := arrayKey{index: int64(n)}
			if n < len(frame.argKeys) {
				key = frame.argKeys[n]
			}
			args.a = args.a.with(key, arg)
		}
		add("args", args)
		result.a = result.a.with(arrayKey{index: int64(len(frames) - 1 - i)}, entry)
	}
	return result
}

// PHP prints the oldest cause first, with "Next" before each replacing
// exception. Keep object identities in the chain and stop on a cycle.
func exceptionDiagnostic(object value) string {
	var chain []value
	seen := make(map[*phpObject]bool)
	for object.kind == objectKind && object.o != nil && object.o.exception != nil && !seen[object.o] {
		seen[object.o] = true
		chain = append(chain, object)
		object = object.o.exception.previous
	}
	var result strings.Builder
	for i := len(chain) - 1; i >= 0; i-- {
		item := chain[i]
		data := item.o.exception
		if i != len(chain)-1 {
			result.WriteString("\n\nNext ")
		}
		trace := strings.TrimSuffix(formatStackTrace(data.frames), "\n")
		if data.message == "" {
			fmt.Fprintf(&result, "%s in %s:%d\nStack trace:\n%s", item.o.className(), data.file, data.line, trace)
		} else {
			fmt.Fprintf(&result, "%s: %s in %s:%d\nStack trace:\n%s", item.o.className(), data.message, data.file, data.line, trace)
		}
	}
	return result.String()
}

func chainException(next, prior value) {
	if next.kind != objectKind || prior.kind != objectKind || next.o == nil || prior.o == nil || next.o.exception == nil || prior.o.exception == nil || next.o == prior.o {
		return
	}
	seen := map[*phpObject]bool{next.o: true}
	for cursor := next.o; cursor != nil && cursor.exception != nil; {
		priorSeen := make(map[*phpObject]bool)
		for ancestor := prior.o; ancestor != nil && ancestor.exception != nil && !priorSeen[ancestor]; {
			priorSeen[ancestor] = true
			if seen[ancestor] {
				return
			}
			p := ancestor.exception.previous
			if p.kind != objectKind || p.o == nil {
				break
			}
			ancestor = p.o
		}
		previous := cursor.exception.previous
		if previous.kind != objectKind || previous.o == nil {
			cursor.exception.previous = prior
			return
		}
		if seen[previous.o] {
			return
		}
		seen[previous.o] = true
		cursor = previous.o
	}
}

func (e *evaluator) throwValue(v value, line int) error {
	if v.kind != objectKind || v.o == nil || !isThrowableClass(v.o.class) {
		return &runtimeError{line, "Error", "Can only throw objects"}
	}
	if v.o.exception == nil {
		v.o.exception = &exceptionData{file: e.filename, line: line}
	}
	if !v.o.exception.captured {
		e.captureFailure()
		v.o.exception.frames = append([]callFrame(nil), e.failureFrames...)
		v.o.exception.captured = true
	}
	return &thrown{object: v}
}

func (e *evaluator) catchable(err error) (*thrown, bool) {
	var raised *thrown
	if errors.As(err, &raised) {
		return raised, true
	}
	var failure *runtimeError
	if !errors.As(err, &failure) || failure.class == "" {
		return nil, false
	}
	e.initClasses()
	class := e.classes[strings.ToLower(failure.class)]
	if class == nil || !isThrowableClass(class) {
		return nil, false
	}
	filename := e.filename
	var sourced *sourceError
	if errors.As(err, &sourced) {
		filename = sourced.filename
	}
	e.captureFailure()
	object := value{kind: objectKind, o: &phpObject{class: class, properties: &phpArray{}, exception: &exceptionData{
		message: failure.message, file: filename, line: failure.line, frames: append([]callFrame(nil), e.failureFrames...), captured: true,
	}}}
	object.o.properties = object.o.properties.with(arrayKey{isString: true, text: "message"}, value{kind: stringKind, s: failure.message})
	object.o.properties = object.o.properties.with(arrayKey{isString: true, text: "code"}, integer(0))
	var pdo *pdoFailure
	if errors.As(err, &pdo) {
		object.o.exception.codeText = pdo.state
		object.o.properties = object.o.properties.with(arrayKey{isString: true, text: "code"}, value{kind: stringKind, s: pdo.state})
		info := arrayValue()
		info.a = info.a.with(arrayKey{index: 0}, value{kind: stringKind, s: pdo.state})
		info.a = info.a.with(arrayKey{index: 1}, integer(pdo.native))
		if !pdo.impl {
			info.a = info.a.with(arrayKey{index: 2}, value{kind: stringKind, s: pdo.detail})
		}
		object.o.properties = object.o.properties.with(arrayKey{isString: true, text: "errorInfo"}, info)
	} else if equalName(class.name, "PDOException") {
		object.o.properties = object.o.properties.with(arrayKey{isString: true, text: "errorInfo"}, value{kind: nullKind})
	}
	return &thrown{object: object}, true
}

func throwableMatches(class *classDecl, name string) bool {
	for current := class; current != nil; current = current.parent {
		if equalName(current.name, name) || classImplements(current, name) {
			return true
		}
	}
	return false
}

func (e *evaluator) executeTry(s statement) error {
	pending := e.execute(s.thenBranch)
	if pending != nil {
		if raised, ok := e.catchable(pending); ok {
			for _, clause := range s.catches {
				matches := false
				for _, name := range clause.types {
					if throwableMatches(raised.object.o.class, name) {
						matches = true
						break
					}
				}
				if !matches {
					continue
				}
				e.failureFrames = nil
				if clause.name != "" {
					if err := e.variableSlotAt(clause.name, s.line).set(raised.object); err != nil {
						return err
					}
				}
				pending = e.execute(clause.body)
				break
			}
		}
	}
	// A declaration failure aborts execution. It is not a Throwable, and a
	// return (or throw) in finally cannot replace that fatal diagnostic.
	var fatal *fatalDiagnostic
	if errors.As(pending, &fatal) {
		return pending
	}
	var stopped *exitSignal
	if errors.As(pending, &stopped) {
		return pending
	}
	if s.finally != nil {
		prior, hasPrior := e.catchable(pending)
		priorFrames := e.failureFrames
		e.failureFrames = nil
		if final := e.execute(s.finally); final != nil {
			if hasPrior {
				if replacement, ok := e.catchable(final); ok {
					chainException(replacement.object, prior.object)
					final = replacement
				}
			}
			return final
		}
		e.failureFrames = priorFrames
	}
	return pending
}
