package php

func (e *evaluator) callRange(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "start"}, {name: "end"}, {name: "step"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 2 || len(args) > 3 || !present[0] || !present[1] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "range() expects 2 to 3 arguments"}
	}
	start, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "range", 1, "start", "int|float|string", bound[0])
	}
	end, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "range", 2, "end", "int|float|string", bound[1])
	}
	step := int64(1)
	if present[2] {
		converted, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "range", 3, "step", "int|float", bound[2])
		}
		step = converted.i
	}
	if step < 1 {
		return value{}, &runtimeError{node.line, "ValueError", "range(): Argument #3 ($step) must be greater than 0"}
	}
	out := arrayValue()
	for i := start.i; ; {
		out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, integer(i))
		if len(out.a.entries) > 100000 {
			return value{}, &runtimeError{node.line, "Error", "gophp: range() size limit exceeded"}
		}
		if start.i <= end.i {
			if i >= end.i || end.i-i < step {
				break
			}
			i += step
		} else {
			if i <= end.i || i-end.i < step {
				break
			}
			i -= step
		}
	}
	return out, nil
}
