package php

func (e *evaluator) callDebugBacktrace(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "options"}, {name: "limit"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) > 2 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "debug_backtrace() expects at most 2 arguments"}
	}
	options := int64(1)
	if present[0] {
		options = bound[0].i
	}
	limit := len(e.frames)
	if present[1] && bound[1].i > 0 && int(bound[1].i) < limit {
		limit = int(bound[1].i)
	}
	out := arrayValue()
	frames := e.snapshotCallFrames()
	for i := len(frames) - 1; i >= 0 && len(out.a.entries) < limit; i-- {
		frame := frames[i]
		item := arrayValue()
		add := func(key string, v value) { item.a = item.a.with(arrayKey{isString: true, text: key}, v) }
		add("file", value{kind: stringKind, s: frame.filename})
		add("line", integer(int64(frame.line)))
		add("function", value{kind: stringKind, s: frame.method})
		if frame.class != "" {
			add("class", value{kind: stringKind, s: frame.class})
			if options&1 != 0 && frame.object.kind == objectKind && frame.object.o != nil {
				add("object", frame.object)
			}
			typ := "->"
			if frame.static {
				typ = "::"
			}
			add("type", value{kind: stringKind, s: typ})
		}
		if options&2 == 0 {
			arguments := arrayValue()
			for j, arg := range frame.args {
				arguments.a = arguments.a.with(arrayKey{index: int64(j)}, arg)
			}
			add("args", arguments)
		}
		out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, item)
	}
	return out, nil
}
