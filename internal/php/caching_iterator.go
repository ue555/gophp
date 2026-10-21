package php

import (
	"fmt"
	"strings"
)

const cachingFullCache = int64(256)
// SPL reserves the upper bits for iterator state. Lower bits, including
// unrecognized ones, remain visible through getFlags().
const cachingPublicFlags = int64(0xffff)

func (e *evaluator) invokeCachingIterator(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil {
		return value{}, &runtimeError{line, "Error", "Non-static method CachingIterator::" + method.name + "() cannot be called statically"}
	}
	o := object.o
	name := strings.ToLower(method.name)
	if name == "__construct" {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "iterator"}, {name: "flags"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] || len(args) > 2 {
			return value{}, &runtimeError{line, "ArgumentCountError", "CachingIterator::__construct() expects 1 or 2 arguments"}
		}
		inner := bound[0]
		if inner.kind != objectKind || inner.o == nil || !classImplements(inner.o.class, "Iterator") {
			return value{}, builtinTypeError(line, "CachingIterator::__construct", 1, "iterator", "Iterator", inner)
		}
		flags := int64(1)
		if present[1] {
			v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(line, "CachingIterator::__construct", 2, "flags", "int", bound[1])
			}
			flags = v.i
		}
		if err := checkCachingFlags(flags, line, 2); err != nil {
			return value{}, err
		}
		o.iteratorInner = inner
		o.cachingFlags = flags & cachingPublicFlags
		o.cachingArray = &phpArray{}
		o.cachingValid = false
		return value{kind: nullKind}, nil
	}
	if o.cachingArray == nil {
		return value{}, &runtimeError{line, "Error", "Object is not initialized"}
	}
	if name == "setflags" {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "flags"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] || len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "CachingIterator::setFlags() expects exactly 1 argument"}
		}
		v, _, ok := e.coerceType(bound[0], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "CachingIterator::setFlags", 1, "flags", "int", bound[0])
		}
		if err := checkCachingFlags(v.i, line, 1); err != nil {
			return value{}, err
		}
		if o.cachingFlags&1 != 0 && v.i&1 == 0 {
			return value{}, &runtimeError{line, "InvalidArgumentException", "Unsetting flag CALL_TO_STRING is not possible"}
		}
		if o.cachingFlags&8 != 0 && v.i&8 == 0 {
			return value{}, &runtimeError{line, "InvalidArgumentException", "Unsetting flag TOSTRING_USE_INNER is not possible"}
		}
		if v.i&cachingFullCache != 0 && o.cachingFlags&cachingFullCache == 0 {
			o.cachingArray = &phpArray{}
		}
		o.cachingFlags = v.i & cachingPublicFlags
		return value{kind: nullKind}, nil
	}
	if strings.HasPrefix(name, "offset") {
		return e.cachingOffset(object, name, args, line)
	}
	if len(args) != 0 {
		return value{}, &runtimeError{line, "ArgumentCountError", "CachingIterator::" + method.name + "() expects exactly 0 arguments"}
	}
	switch name {
	case "getflags":
		flags := o.cachingFlags
		if o.cachingValid {
			flags |= 0x10000 // CIT_VALID is visible through getFlags().
		}
		return integer(flags), nil
	case "getcache", "count":
		if o.cachingFlags&cachingFullCache == 0 {
			return value{}, &runtimeError{line, "BadMethodCallException", "CachingIterator does not use a full cache (see CachingIterator::__construct)"}
		}
		if name == "count" {
			return integer(int64(len(o.cachingArray.entries))), nil
		}
		return value{kind: arrayKind, a: o.cachingArray}, nil
	case "__tostring":
		if o.cachingFlags&15 == 0 {
			return value{}, &runtimeError{line, "BadMethodCallException", "CachingIterator does not fetch string value (see CachingIterator::__construct)"}
		}
		if !o.cachingValid {
			return value{kind: stringKind}, nil
		}
		if o.cachingFlags&9 != 0 {
			return value{kind: stringKind, s: o.cachingText}, nil
		}
		v := o.cachingCurrent
		if o.cachingFlags&2 != 0 {
			v = o.cachingKey
		}
		s, err := e.toText(v, line)
		return value{kind: stringKind, s: s}, err
	}
	return e.cachingIteratorMethod(object, name, line)
}

func checkCachingFlags(flags int64, line, argument int) error {
	toString := flags & 15
	if toString != 0 && toString&(toString-1) != 0 {
		name := "CachingIterator::__construct"
		if argument == 1 {
			name = "CachingIterator::setFlags"
		}
		return &runtimeError{line, "ValueError", fmt.Sprintf("%s(): Argument #%d ($flags) must contain only one of CachingIterator::CALL_TOSTRING, CachingIterator::TOSTRING_USE_KEY, CachingIterator::TOSTRING_USE_CURRENT, or CachingIterator::TOSTRING_USE_INNER", name, argument)}
	}
	return nil
}

// PHP's caching iterator keeps the current item while the inner iterator is
// already positioned on the next one. hasNext() therefore observes the inner
// cursor without advancing it again.
func (e *evaluator) cachingAdvance(o *phpObject, line int) error {
	valid, err := e.iteratorMethod(o.iteratorInner, "valid", line)
	if err != nil {
		return err
	}
	if !valid.truthy() {
		o.cachingValid = false
		return nil
	}
	// spl_dual_it_fetch() reads the value before the key. User Iterators can
	// observe this order (and either method may have side effects).
	current, err := e.iteratorMethod(o.iteratorInner, "current", line)
	if err != nil {
		return err
	}
	key, err := e.iteratorMethod(o.iteratorInner, "key", line)
	if err != nil {
		return err
	}
	o.cachingKey, o.cachingCurrent, o.cachingValid = key, current, true
	if o.cachingFlags&cachingFullCache != 0 {
		arrayKey, keyErr := e.arrayKey(key, line)
		if keyErr != nil {
			return keyErr
		}
		o.cachingArray = o.cachingArray.with(arrayKey, current)
	}
	if o.cachingFlags&9 != 0 {
		textValue := current
		if o.cachingFlags&8 != 0 {
			textValue = o.iteratorInner
		}
		text, err := e.toText(textValue, line)
		if err != nil {
			return err
		}
		o.cachingText = text
	}
	_, err = e.iteratorMethod(o.iteratorInner, "next", line)
	return err
}

func (e *evaluator) cachingIteratorMethod(object value, name string, line int) (value, error) {
	o := object.o
	switch strings.ToLower(name) {
	case "getinneriterator":
		return o.iteratorInner, nil
	case "rewind":
		if _, err := e.iteratorMethod(o.iteratorInner, "rewind", line); err != nil {
			return value{}, err
		}
		o.cachingArray = &phpArray{}
		o.cachingValid = false
		return value{kind: nullKind}, e.cachingAdvance(o, line)
	case "next":
		return value{kind: nullKind}, e.cachingAdvance(o, line)
	case "valid":
		return boolean(o.cachingValid), nil
	case "hasnext":
		return e.iteratorMethod(o.iteratorInner, "valid", line)
	case "key":
		if o.cachingValid {
			return o.cachingKey, nil
		}
		return value{kind: nullKind}, nil
	case "current":
		if o.cachingValid {
			return o.cachingCurrent, nil
		}
		return value{kind: nullKind}, nil
	}
	return value{}, &runtimeError{line, "Error", "gophp: CachingIterator method is not implemented"}
}

func (e *evaluator) cachingOffset(object value, name string, args []callArgument, line int) (value, error) {
	o := object.o
	parameters := []parameter{{name: "key"}}
	if name == "offsetset" {
		parameters = append(parameters, parameter{name: "value"})
	}
	bound, present, _, err := bindCallArguments(args, parameters, line)
	if err != nil {
		return value{}, err
	}
	if len(args) != len(parameters) || !present[0] {
		return value{}, &runtimeError{line, "ArgumentCountError", "CachingIterator::" + name + "() expects the required arguments"}
	}
	keyText, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(line, "CachingIterator::"+name, 1, "key", "string", bound[0])
	}
	if o.cachingFlags&cachingFullCache == 0 {
		return value{}, &runtimeError{line, "BadMethodCallException", "CachingIterator does not use a full cache (see CachingIterator::__construct)"}
	}
	key, err := e.arrayKey(keyText, line)
	if err != nil {
		return value{}, err
	}
	switch name {
	case "offsetexists":
		_, ok := o.cachingArray.get(key)
		return boolean(ok), nil
	case "offsetget":
		v, ok := o.cachingArray.get(key)
		if !ok {
			return value{kind: nullKind}, e.warn(line, "Undefined array key \""+keyText.s+"\"")
		}
		return v, nil
	case "offsetset":
		o.cachingArray = o.cachingArray.with(key, bound[1])
	case "offsetunset":
		o.cachingArray = o.cachingArray.without(key)
	}
	return value{kind: nullKind}, nil
}
