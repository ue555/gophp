package php

func (e *evaluator) callIteratorCount(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "iterator"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) != 1 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "iterator_count() expects exactly 1 argument"}
	}
	iterator := bound[0]
	if iterator.kind == arrayKind {
		return integer(int64(len(iterator.a.entries))), nil
	}
	if iterator.kind != objectKind || iterator.o == nil || !classImplements(iterator.o.class, "Traversable") {
		return value{}, builtinTypeError(node.line, "iterator_count", 1, "iterator", "Traversable|array", iterator)
	}
	if classImplements(iterator.o.class, "IteratorAggregate") {
		iterator, err = e.iteratorMethod(iterator, "getIterator", node.line)
		if err != nil {
			return value{}, err
		}
		if iterator.kind != objectKind || iterator.o == nil || !classImplements(iterator.o.class, "Traversable") {
			return value{}, &runtimeError{node.line, "TypeError", "IteratorAggregate::getIterator() must return Traversable"}
		}
	}
	if _, err := e.iteratorMethod(iterator, "rewind", node.line); err != nil {
		return value{}, err
	}
	var count int64
	for {
		valid, err := e.iteratorMethod(iterator, "valid", node.line)
		if err != nil {
			return value{}, err
		}
		if !valid.truthy() {
			return integer(count), nil
		}
		if e.bounded {
			if e.steps == 0 {
				return value{}, &runtimeError{node.line, "Error", "gophp: execution budget exhausted"}
			}
			e.steps--
		}
		count++
		if _, err := e.iteratorMethod(iterator, "next", node.line); err != nil {
			return value{}, err
		}
	}
}

func (e *evaluator) callIteratorToArray(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "iterator"}, {name: "preserve_keys"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "iterator_to_array() expects 1 to 2 arguments"}
	}
	iterator := bound[0]
	preserve := true
	if present[1] {
		preserve = bound[1].truthy()
	}
	return e.iteratorToArrayValue(iterator, preserve, node.line)
}

func (e *evaluator) iteratorToArrayValue(iterator value, preserve bool, line int) (value, error) {
	if iterator.kind == arrayKind {
		if preserve {
			return iterator, nil
		}
		out := arrayValue()
		for _, entry := range iterator.a.entries {
			k, keyErr := out.a.appendKey(line)
			if keyErr != nil {
				return value{}, keyErr
			}
			out.a = out.a.with(k, dereference(entry.value))
		}
		return out, nil
	}
	if iterator.kind != objectKind || !classImplements(iterator.o.class, "Traversable") {
		return value{}, builtinTypeError(line, "iterator_to_array", 1, "iterator", "Traversable|array", iterator)
	}
	if classImplements(iterator.o.class, "IteratorAggregate") {
		var err error
		iterator, err = e.iteratorMethod(iterator, "getIterator", line)
		if err != nil {
			return value{}, err
		}
		if iterator.kind != objectKind || !classImplements(iterator.o.class, "Traversable") {
			return value{}, &runtimeError{line, "TypeError", "IteratorAggregate::getIterator() must return Traversable"}
		}
	}
	if _, err := e.iteratorMethod(iterator, "rewind", line); err != nil {
		return value{}, err
	}
	out := arrayValue()
	for {
		valid, err := e.iteratorMethod(iterator, "valid", line)
		if err != nil {
			return value{}, err
		}
		if !valid.truthy() {
			return out, nil
		}
		key, err := e.iteratorMethod(iterator, "key", line)
		if err != nil {
			return value{}, err
		}
		item, err := e.iteratorMethod(iterator, "current", line)
		if err != nil {
			return value{}, err
		}
		var k arrayKey
		if preserve {
			k, err = e.arrayKey(key, line)
		} else {
			k, err = out.a.appendKey(line)
		}
		if err != nil {
			return value{}, err
		}
		out.a = out.a.with(k, item)
		if _, err = e.iteratorMethod(iterator, "next", line); err != nil {
			return value{}, err
		}
	}
}
