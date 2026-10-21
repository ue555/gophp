package php

import "strings"

func priorityHigher(a, b value) bool {
	a, b = dereference(a), dereference(b)
	if a.kind == intKind && b.kind == intKind {
		return a.i > b.i
	}
	if a.kind == floatKind || b.kind == floatKind {
		return a.asFloat() > b.asFloat()
	}
	return a.text() > b.text()
}

func priorityResult(entry priorityEntry, flags int64) value {
	switch flags {
	case 2:
		return entry.priority
	case 3:
		out := arrayValue()
		out.a = out.a.with(arrayKey{isString: true, text: "data"}, entry.data)
		out.a = out.a.with(arrayKey{isString: true, text: "priority"}, entry.priority)
		return out
	default:
		return entry.data
	}
}

func (e *evaluator) invokePriorityQueue(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method SplPriorityQueue::" + method.name + "() cannot be called statically"}
	}
	o := object.o
	if o.priorityFlags == 0 {
		o.priorityFlags = 1
	}
	switch strings.ToLower(method.name) {
	case "insert":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}, {name: "priority"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] || !present[1] {
			return value{}, &runtimeError{line, "ArgumentCountError", "SplPriorityQueue::insert() expects exactly 2 arguments"}
		}
		entry := priorityEntry{data: bound[0], priority: bound[1]}
		i := len(o.priorityItems)
		o.priorityItems = append(o.priorityItems, entry)
		for i > 0 && priorityHigher(entry.priority, o.priorityItems[i-1].priority) {
			o.priorityItems[i] = o.priorityItems[i-1]
			i--
		}
		o.priorityItems[i] = entry
		return boolean(true), nil
	case "count":
		return integer(int64(len(o.priorityItems))), nil
	case "isempty", "valid":
		return boolean((len(o.priorityItems) == 0) == (strings.EqualFold(method.name, "isEmpty"))), nil
	case "rewind":
		return value{kind: nullKind}, nil
	case "key":
		return integer(o.priorityIndex), nil
	case "getextractflags":
		return integer(o.priorityFlags), nil
	case "setextractflags":
		if len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "SplPriorityQueue::setExtractFlags() expects exactly 1 argument"}
		}
		v, _, ok := e.coerceType(args[0].value, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "SplPriorityQueue::setExtractFlags", 1, "flags", "int", args[0].value)
		}
		if v.i < 1 || v.i > 3 {
			return value{}, &runtimeError{line, "RuntimeException", "Extract mode must be one of SplPriorityQueue::EXTR_DATA, SplPriorityQueue::EXTR_PRIORITY, or SplPriorityQueue::EXTR_BOTH"}
		}
		o.priorityFlags = v.i
		return integer(v.i), nil
	case "top", "current", "extract", "next":
		if len(o.priorityItems) == 0 {
			return value{}, &runtimeError{line, "RuntimeException", "Can't extract from an empty heap"}
		}
		entry := o.priorityItems[0]
		if strings.EqualFold(method.name, "extract") || strings.EqualFold(method.name, "next") {
			o.priorityItems = o.priorityItems[1:]
			o.priorityIndex++
		}
		if strings.EqualFold(method.name, "next") {
			return value{kind: nullKind}, nil
		}
		return priorityResult(entry, o.priorityFlags), nil
	}
	return value{}, &runtimeError{line, "Error", "gophp: SplPriorityQueue method is not implemented"}
}
