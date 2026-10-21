package php

func (e *evaluator) callArraySum(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) != 1 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "array_sum() expects exactly 1 argument"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "array_sum", 1, "array", "array", bound[0])
	}
	total := integer(0)
	for _, entry := range bound[0].a.entries {
		v := dereference(entry.value)
		if v.kind == stringKind {
			numeric, _, ok := v.numeric()
			if ok {
				v = numeric
			} else {
				v = integer(0)
			}
		}
		if v.kind == boolKind {
			v = integer(v.i)
		}
		if v.kind == nullKind {
			continue
		}
		if v.kind != intKind && v.kind != floatKind {
			return value{}, &runtimeError{node.line, "TypeError", "array_sum(): Addition is not supported for this value"}
		}
		if total.kind == floatKind || v.kind == floatKind {
			left := float64(total.i)
			if total.kind == floatKind {
				left = total.f
			}
			right := float64(v.i)
			if v.kind == floatKind {
				right = v.f
			}
			total = floating(left + right)
		} else {
			total = integer(total.i + v.i)
		}
	}
	return total, nil
}
