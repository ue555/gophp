package php

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type arrayKey struct {
	text     string
	index    int64
	isString bool
}
type arrayEntry struct {
	key   arrayKey
	value value
}

// Arrays are immutable after publication. Each write copies the affected path;
// nested arrays have value semantics, while objects retain shared identity.
type phpArray struct {
	entries []arrayEntry
	next    int64
	hasInt  bool
	full    bool
}
type phpObject struct {
	properties *phpArray
	class      *classDecl
	closure    *phpClosure
	exception  *exceptionData
}

func (o *phpObject) className() string {
	if o.class == nil {
		return "stdClass"
	}
	return o.class.name
}

func (o *phpObject) dynamicProperty(name string) bool {
	if o.class == nil || strings.EqualFold(o.class.name, "stdClass") {
		return false
	}
	for _, property := range o.class.properties {
		if property.name == name {
			return false
		}
	}
	return true
}

func arrayValue() value { return value{kind: arrayKind, a: &phpArray{}} }
func (a *phpArray) find(k arrayKey) int {
	for i, entry := range a.entries {
		if entry.key == k {
			return i
		}
	}
	return -1
}
func (a *phpArray) get(k arrayKey) (value, bool) {
	if i := a.find(k); i >= 0 {
		return a.entries[i].value, true
	}
	return value{}, false
}
func (a *phpArray) with(k arrayKey, v value) *phpArray {
	b := *a
	b.entries = append([]arrayEntry(nil), a.entries...)
	if i := b.find(k); i >= 0 {
		b.entries[i].value = v
	} else {
		b.entries = append(b.entries, arrayEntry{k, v})
	}
	if !k.isString && (!b.hasInt || k.index >= b.next) {
		b.hasInt = true
		if k.index == math.MaxInt64 {
			b.next, b.full = k.index, true
		} else {
			b.next = k.index + 1
		}
	}
	return &b
}
func (a *phpArray) appendKey(line int) (arrayKey, error) {
	if a.full && a.find(arrayKey{index: a.next}) >= 0 {
		return arrayKey{}, &runtimeError{line, "Error", "Cannot add element to the array as the next element is already occupied"}
	}
	return arrayKey{index: a.next}, nil
}

func (e *evaluator) arrayKey(v value, line int) (arrayKey, error) {
	switch v.kind {
	case intKind, boolKind:
		return arrayKey{index: v.i}, nil
	case nullKind:
		if e.version == PHP85 {
			e.diagnostic("Deprecated", line, "Using null as an array offset is deprecated, use an empty string instead")
		}
		return arrayKey{isString: true}, nil
	case stringKind:
		// No plus sign, leading zero, -0, or surrounding whitespace converts.
		if n, err := strconv.ParseInt(v.s, 10, 64); err == nil && strconv.FormatInt(n, 10) == v.s {
			return arrayKey{index: n}, nil
		}
		return arrayKey{text: v.s, isString: true}, nil
	case floatKind:
		n := floatToInt(v.f)
		outOfRange := math.IsNaN(v.f) || math.IsInf(v.f, 0) || v.f >= math.Ldexp(1, 63) || v.f < -math.Ldexp(1, 63)
		if e.version == PHP85 && outOfRange {
			e.warn(line, fmt.Sprintf("The float %s is not representable as an int, cast occurred", floatText(v.f, -1)))
		}
		// Zend's range macro regards NaN as in range (both ordered tests are
		// false), so 8.5 emits the precision deprecation after its warning.
		if float64(n) != v.f && (e.version != PHP85 || !outOfRange || math.IsNaN(v.f)) {
			e.diagnostic("Deprecated", line, fmt.Sprintf("Implicit conversion from float %s to int loses precision", floatText(v.f, -1)))
		}
		return arrayKey{index: n}, nil
	default:
		return arrayKey{}, &runtimeError{line, "TypeError", "Cannot access offset of type " + v.typeName() + " on array"}
	}
}

// PHP's finite float-to-int cast wraps modulo 2^64; NaN and infinity become 0.
func floatToInt(f float64) int64 {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	f = math.Mod(math.Trunc(f), math.Ldexp(1, 64))
	if f >= math.Ldexp(1, 63) {
		f -= math.Ldexp(1, 64)
	}
	if f < -math.Ldexp(1, 63) {
		f += math.Ldexp(1, 64)
	}
	return int64(f)
}

func (e *evaluator) evalArray(node arrayLiteral) (value, error) {
	result := arrayValue()
	for _, element := range node.elements {
		line := expressionLine(element.value)
		values, err := e.evalOperands(element.key, element.value, line, true)
		if err != nil {
			return value{}, err
		}
		var k arrayKey
		if element.key == nil {
			k, err = result.a.appendKey(line)
		} else {
			k, err = e.arrayKey(values[0], line)
		}
		if err != nil {
			return value{}, err
		}
		result.a = result.a.with(k, values[1])
	}
	return result, nil
}

func (e *evaluator) readOffset(node offsetExpr) (value, error) {
	if node.key == nil {
		return value{}, &runtimeError{node.line, "Error", "Cannot use [] for reading"}
	}
	values, err := e.evalOperands(node.base, node.key, node.line, false)
	if err != nil {
		return value{}, err
	}
	base, raw := values[0], values[1]
	if base.kind == objectKind {
		return value{}, &runtimeError{node.line, "Error", "Cannot use object of type " + base.o.className() + " as array"}
	}
	if base.kind == stringKind {
		return e.stringOffset(base, raw, node.line)
	}
	if base.kind != arrayKind {
		name := base.typeName()
		if base.kind == boolKind {
			name = "false"
			if base.i != 0 {
				name = "true"
			}
		}
		e.warn(node.line, "Trying to access array offset on "+name)
		return value{}, nil
	}
	k, err := e.arrayKey(raw, node.line)
	if err != nil {
		return value{}, err
	}
	if v, ok := base.a.get(k); ok {
		return v, nil
	}
	label := strconv.FormatInt(k.index, 10)
	if k.isString {
		label = "\"" + k.text + "\""
	}
	e.warn(node.line, "Undefined array key "+label)
	return value{}, nil
}

func (e *evaluator) stringOffset(base, key value, line int) (value, error) {
	if key.kind == stringKind {
		n, trailing, ok := key.numeric()
		if !ok || n.kind != intKind {
			return value{}, &runtimeError{line, "TypeError", "Cannot access offset of type string on string"}
		}
		if trailing {
			e.warn(line, "Illegal string offset \""+key.s+"\"")
		}
		key = n
	} else if key.kind != intKind {
		if key.kind == arrayKind || key.kind == objectKind {
			return value{}, &runtimeError{line, "TypeError", "Cannot access offset of type " + key.typeName() + " on string"}
		}
		e.warn(line, "String offset cast occurred")
		if key.kind == floatKind {
			key = integer(floatToInt(key.f))
		} else {
			key = integer(key.i)
		}
	}
	i := key.i
	if i < 0 {
		i += int64(len(base.s))
	}
	if i < 0 || i >= int64(len(base.s)) {
		e.warn(line, fmt.Sprintf("Uninitialized string offset %d", key.i))
		return value{kind: stringKind}, nil
	}
	return value{kind: stringKind, s: base.s[i : i+1]}, nil
}

func (e *evaluator) readProperty(node propertyExpr) (value, error) {
	base, err := e.eval(node.base)
	if err != nil {
		return value{}, err
	}
	if base.kind != objectKind {
		e.warn(node.line, fmt.Sprintf("Attempt to read property \"%s\" on %s", node.name, base.typeName()))
		return value{}, nil
	}
	k, err := e.propertyKey(base.o, node.name, node.line)
	if err != nil {
		return value{}, err
	}
	if v, ok := base.o.properties.get(k); ok {
		return v, nil
	}
	e.warn(node.line, "Undefined property: "+base.o.className()+"::$"+node.name)
	return value{}, nil
}

// A location resolves the delayed write path after RHS subexpressions.
type location struct {
	get func() value
	set func(value) error
}

type locationMode uint8

const (
	locationWrite locationMode = iota
	locationReadWrite
	locationIncrement
	locationModify
)

func (e *evaluator) locate(target expr, modes ...locationMode) (location, error) {
	mode := locationWrite
	if len(modes) > 0 {
		mode = modes[0]
	}
	rw := mode != locationWrite
	switch node := target.(type) {
	case variable:
		if e.variables == nil {
			e.variables = make(map[string]value)
		}
		if _, exists := e.readVariable(node.token.text); rw && !exists {
			e.warn(node.token.line, "Undefined variable $"+node.token.text)
		}
		return e.variableSlot(node.token.text), nil
	case staticPropertyExpr:
		return e.staticPropertyLocation(node)
	case propertyExpr:
		var base value
		var err error
		if v, ok := node.base.(variable); ok && v.token.text != "this" {
			base, _ = e.readVariable(v.token.text)
			if _, exists := e.readVariable(v.token.text); rw && !exists {
				e.warn(v.token.line, "Undefined variable $"+v.token.text)
			}
		} else {
			base, err = e.eval(node.base)
		}
		if err != nil {
			return location{}, err
		}
		if base.kind != objectKind {
			action := "assign"
			if mode == locationIncrement {
				action = "increment/decrement"
			}
			if mode == locationModify {
				action = "modify"
			}
			name := base.typeName()
			if base.kind == boolKind {
				if base.i == 0 {
					name = "false"
				} else {
					name = "true"
				}
			}
			return location{}, &runtimeError{node.line, "Error", fmt.Sprintf("Attempt to %s property \"%s\" on %s", action, node.name, name)}
		}
		k, err := e.propertyKey(base.o, node.name, node.line)
		if err != nil {
			return location{}, err
		}
		if rw && base.o.properties.find(k) < 0 {
			if base.o.dynamicProperty(node.name) {
				e.diagnostic("Deprecated", node.line, "Creation of dynamic property "+base.o.className()+"::$"+node.name+" is deprecated")
			}
			e.warn(node.line, "Undefined property: "+base.o.className()+"::$"+node.name)
			base.o.properties = base.o.properties.with(k, value{})
		}
		return location{func() value { v, _ := base.o.properties.get(k); return v }, func(v value) error {
			if base.o.dynamicProperty(node.name) && base.o.properties.find(k) < 0 {
				e.diagnostic("Deprecated", node.line, "Creation of dynamic property "+base.o.className()+"::$"+node.name+" is deprecated")
			}
			base.o.properties = base.o.properties.with(k, v)
			return nil
		}}, nil
	case offsetExpr:
		parentMode := locationWrite
		if rw {
			parentMode = locationModify
		}
		parent, err := e.locate(node.base, parentMode)
		if err != nil {
			return location{}, err
		}
		base := parent.get()
		if base.kind == nullKind || base.kind == boolKind && base.i == 0 {
			if base.kind == boolKind {
				e.diagnostic("Deprecated", node.line, "Automatic conversion of false to array is deprecated")
			}
			base = arrayValue()
			if err := parent.set(base); err != nil {
				return location{}, err
			}
		}
		if base.kind != arrayKind {
			message := "Cannot use a scalar value as an array"
			if base.kind == objectKind {
				message = "Cannot use object of type " + base.o.className() + " as array"
			}
			if base.kind == stringKind {
				message = "unsupported syntax: string offset assignment"
			}
			return location{}, &runtimeError{node.line, "Error", message}
		}
		var k arrayKey
		if node.key == nil {
			k, err = base.a.appendKey(node.line)
		} else {
			var raw value
			raw, err = e.eval(node.key)
			if err == nil {
				k, err = e.arrayKey(raw, node.line)
			}
		}
		if err != nil {
			return location{}, err
		}
		// Reserve intermediate append slots for nested writes (e.g. $a[][]).
		if rw && node.key != nil && base.a.find(k) < 0 {
			if k.isString {
				e.warn(node.line, fmt.Sprintf("Undefined array key \"%s\"", k.text))
			} else {
				e.warn(node.line, fmt.Sprintf("Undefined array key %d", k.index))
			}
		}
		// The RHS has already been evaluated at this point.
		if node.key == nil {
			base.a = base.a.with(k, value{})
			if err := parent.set(base); err != nil {
				return location{}, err
			}
		}
		return location{func() value {
			base := parent.get()
			if base.kind != arrayKind {
				return value{}
			}
			v, _ := base.a.get(k)
			return v
		}, func(v value) error {
			latest := parent.get()
			if latest.kind != arrayKind {
				return &runtimeError{node.line, "Error", "Cannot use a scalar value as an array"}
			}
			latest.a = latest.a.with(k, v)
			return parent.set(latest)
		}}, nil
	}
	return location{}, fmt.Errorf("unsupported assignment target")
}

// Target key subexpressions execute before the RHS, but the writes and bare
// variable fetches are delayed, as in Zend's delayed compilation of assignments.
func (e *evaluator) prepareTarget(target expr) (expr, error) {
	switch node := target.(type) {
	case variable:
		return node, nil
	case staticPropertyExpr:
		return node, nil
	case offsetExpr:
		base, err := e.prepareTarget(node.base)
		if err != nil {
			return nil, err
		}
		node.base = base
		if node.key != nil {
			if _, bare := node.key.(variable); !bare {
				v, err := e.eval(node.key)
				if err != nil {
					return nil, err
				}
				node.key = literal{v, expressionLine(node.key)}
			}
		}
		return node, nil
	case propertyExpr:
		base, err := e.prepareTarget(node.base)
		node.base = base
		return node, err
	default:
		v, err := e.eval(target)
		return literal{v, expressionLine(target)}, err
	}
}

func rootVariable(target expr) string {
	switch node := target.(type) {
	case variable:
		return node.token.text
	case offsetExpr:
		return rootVariable(node.base)
	case propertyExpr:
		return rootVariable(node.base)
	}
	return ""
}

func (e *evaluator) evalAssignment(node assignment) (value, error) {
	prepared, err := e.prepareTarget(node.target)
	if err != nil {
		return value{}, err
	}
	var v value
	rhs, bare := node.source.(variable)
	_, dimension := node.target.(offsetExpr)
	deferRHS := bare && dimension && rhs.token.text != rootVariable(node.target)
	if !deferRHS {
		v, err = e.eval(node.source)
		if err != nil {
			return value{}, err
		}
	}
	target, err := e.locate(prepared)
	if err != nil {
		return value{}, err
	}
	if deferRHS {
		v, err = e.eval(node.source)
		if err != nil {
			return value{}, err
		}
	}
	if err := target.set(v); err != nil {
		return value{}, err
	}
	return v, nil
}
