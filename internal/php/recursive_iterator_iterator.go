package php

import "strings"

type recursiveItem struct {
	key   value
	item  value
	inner value
	depth int
}

func (e *evaluator) invokeRecursiveIteratorIterator(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method RecursiveIteratorIterator::" + method.name + "() cannot be called statically"}
	}
	o := object.o
	switch strings.ToLower(method.name) {
	case "__construct":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "iterator"}, {name: "mode"}, {name: "flags"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", "RecursiveIteratorIterator::__construct() expects at least 1 argument"}
		}
		if bound[0].kind != objectKind || bound[0].o == nil || !classImplements(bound[0].o.class, "RecursiveIterator") {
			return value{}, builtinTypeError(line, "RecursiveIteratorIterator::__construct", 1, "iterator", "RecursiveIterator", bound[0])
		}
		o.recursiveInner = bound[0]
		if present[1] {
			mode, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(line, "RecursiveIteratorIterator::__construct", 2, "mode", "int", bound[1])
			}
			o.directoryFlags = mode.i
		}
		return value{kind: nullKind}, nil
	case "rewind":
		o.recursiveItems = nil
		o.recursiveIndex = 0
		if err := e.collectRecursiveItems(o.recursiveInner, int(o.directoryFlags), 0, o.recursiveMaxDepth-1, &o.recursiveItems, line); err != nil {
			return value{}, err
		}
		return value{kind: nullKind}, nil
	case "setmaxdepth":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "maxDepth"}}, line)
		if err != nil {
			return value{}, err
		}
		if len(args) > 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "RecursiveIteratorIterator::setMaxDepth() expects at most 1 argument"}
		}
		maxDepth := int64(-1)
		if present[0] {
			converted, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(line, "RecursiveIteratorIterator::setMaxDepth", 1, "maxDepth", "int", bound[0])
			}
			maxDepth = converted.i
		}
		if maxDepth < -1 {
			return value{}, &runtimeError{line, "ValueError", "RecursiveIteratorIterator::setMaxDepth(): Argument #1 ($maxDepth) must be greater than or equal to -1"}
		}
		o.recursiveMaxDepth = int(maxDepth) + 1
		return value{kind: nullKind}, nil
	case "getmaxdepth":
		if o.recursiveMaxDepth == 0 {
			return boolean(false), nil
		}
		return integer(int64(o.recursiveMaxDepth - 1)), nil
	case "valid":
		return boolean(o.recursiveIndex >= 0 && o.recursiveIndex < len(o.recursiveItems)), nil
	case "next":
		o.recursiveIndex++
		return value{kind: nullKind}, nil
	case "key", "current", "getdepth", "getsubiterator":
		if o.recursiveIndex < 0 || o.recursiveIndex >= len(o.recursiveItems) {
			return value{kind: nullKind}, nil
		}
		item := o.recursiveItems[o.recursiveIndex]
		switch strings.ToLower(method.name) {
		case "key":
			return item.key, nil
		case "current":
			return item.item, nil
		case "getdepth":
			return integer(int64(item.depth)), nil
		default:
			return item.inner, nil
		}
	case "getinneriterator":
		return o.recursiveInner, nil
	}
	return value{}, &runtimeError{line, "Error", "gophp: RecursiveIteratorIterator method is not implemented"}
}

// collectRecursiveItems flattens the recursive iterator eagerly. maxDepth < 0
// means unlimited; at the maximum depth an element that has children is not
// descended into, yet it still counts as having children (LEAVES_ONLY skips
// it, SELF_FIRST and CHILD_FIRST yield it), matching ext/spl.
func (e *evaluator) collectRecursiveItems(inner value, mode, depth, maxDepth int, out *[]recursiveItem, line int) error {
	if depth > 64 {
		return &runtimeError{line, "Error", "gophp: recursive iterator nesting limit exceeded"}
	}
	if _, err := e.iteratorMethod(inner, "rewind", line); err != nil {
		return err
	}
	for {
		valid, err := e.iteratorMethod(inner, "valid", line)
		if err != nil {
			return err
		}
		if !valid.truthy() {
			return nil
		}
		key, err := e.iteratorMethod(inner, "key", line)
		if err != nil {
			return err
		}
		current, err := e.iteratorMethod(inner, "current", line)
		if err != nil {
			return err
		}
		hasChildren, err := e.iteratorMethod(inner, "hasChildren", line)
		if err != nil {
			return err
		}
		entry := recursiveItem{key: key, item: current, inner: inner, depth: depth}
		if !hasChildren.truthy() || mode == 1 {
			*out = append(*out, entry)
		}
		if hasChildren.truthy() && (maxDepth < 0 || depth < maxDepth) {
			child, err := e.iteratorMethod(inner, "getChildren", line)
			if err != nil {
				return err
			}
			if child.kind != objectKind || child.o == nil {
				return &runtimeError{line, "TypeError", "RecursiveIterator::getChildren() must return a RecursiveIterator"}
			}
			if err := e.collectRecursiveItems(child, mode, depth+1, maxDepth, out, line); err != nil {
				return err
			}
		}
		if hasChildren.truthy() && mode == 2 {
			*out = append(*out, entry)
		}
		if _, err := e.iteratorMethod(inner, "next", line); err != nil {
			return err
		}
	}
}
