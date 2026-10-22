package php

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
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
type priorityEntry struct{ data, priority value }

// Array containers copy on write. References stored in entries retain shared
// cells across container copies; nested arrays have value semantics, while
// objects retain shared identity.
type phpArray struct {
	entries []arrayEntry
	next    int64
	hasInt  bool
	full    bool
}
type phpObject struct {
	dom                *domNode
	pdo                *sqliteConnection
	pdoStatement       *sqliteStatement
	dateTime           *time.Time
	timezone           *time.Location
	weakMap            map[*phpObject]value
	objectStorage      map[*phpObject]value
	properties         *phpArray
	class              *classDecl
	enumCase           string
	closure            *phpClosure
	exception          *exceptionData
	generator          *phpGenerator
	reflected          *classDecl
	reflectedMethod    *methodDecl
	reflectedProperty  *propertyDecl
	reflectedClosure   *phpClosure
	reflectedAttribute *attributeDecl
	attributeScope     *classDecl
	reflectedParameter *parameter
	reflectedOwner     *classDecl
	reflectedTypeName  string
	reflectedNullable  bool
	iteratorArray      *phpArray
	iteratorIndex      int
	iteratorInner      value
	iteratorReady      bool
	iteratorOffset     int64
	iteratorLimit      int64
	iteratorPosition   int64
	appendIterators    []value
	appendIndex        int
	appendReady        bool
	cachingFlags       int64
	cachingCurrent     value
	cachingKey         value
	cachingValid       bool
	cachingArray       *phpArray
	cachingText        string
	filterInner        value
	filterValid        bool
	regexPattern       string
	regexCompiled      *regexp.Regexp
	regexExclude       *regexp.Regexp
	regexMode          int64
	regexCurrent       value
	filePath           string
	directoryPath      string
	directorySubPath   string
	directoryEntries   []string
	directoryIndex     int
	directoryFlags     int64
	recursiveItems     []recursiveItem
	recursiveIndex     int
	recursiveInner     value
	priorityItems      []priorityEntry
	priorityFlags      int64
	priorityIndex      int64
}

func (o *phpObject) className() string {
	if o.class == nil {
		return "stdClass"
	}
	return o.class.name
}

func (o *phpObject) dynamicProperty(name string, scope *classDecl) bool {
	if o.class == nil || strings.EqualFold(o.class.name, "stdClass") {
		return false
	}
	for _, property := range o.class.properties {
		if !property.static && property.name == name && (property.visibility != "private" || property.owner == scope) {
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
	v, ok := a.getRaw(k)
	return dereference(v), ok
}
func (a *phpArray) getRaw(k arrayKey) (value, bool) {
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
			if err := e.diagnostic("Deprecated", line, "Using null as an array offset is deprecated, use an empty string instead"); err != nil {
				return arrayKey{}, err
			}
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
			if err := e.warn(line, fmt.Sprintf("The float %s is not representable as an int, cast occurred", floatText(v.f, -1))); err != nil {
				return arrayKey{}, err
			}
		}
		// Zend's range macro regards NaN as in range (both ordered tests are
		// false), so 8.5 emits the precision deprecation after its warning.
		if float64(n) != v.f && (e.version != PHP85 || !outOfRange || math.IsNaN(v.f)) {
			if err := e.diagnostic("Deprecated", line, fmt.Sprintf("Implicit conversion from float %s to int loses precision", floatText(v.f, -1))); err != nil {
				return arrayKey{}, err
			}
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
		if spread, ok := element.value.(unpackArgument); ok {
			iterable, err := e.eval(spread.operand)
			if err != nil {
				return value{}, err
			}
			if iterable.kind == objectKind && iterable.o != nil && classImplements(iterable.o.class, "Traversable") {
				iterable, err = e.iteratorToArrayValue(iterable, true, spread.line)
				if err != nil {
					return value{}, err
				}
			}
			if iterable.kind != arrayKind {
				return value{}, &runtimeError{spread.line, "Error", "Only arrays and Traversables can be unpacked"}
			}
			for _, entry := range iterable.a.entries {
				key := entry.key
				if !key.isString {
					key, err = result.a.appendKey(spread.line)
					if err != nil {
						return value{}, err
					}
				}
				result.a = result.a.with(key, dereference(entry.value))
			}
			continue
		}
		line := expressionLine(element.value)
		if element.byRef {
			if !referenceTargetSyntax(element.value) {
				return value{}, &runtimeError{line, "Error", "Cannot create references to/from string offsets"}
			}
			var key arrayKey
			if element.key == nil {
				var err error
				key, err = result.a.appendKey(line)
				if err != nil {
					return value{}, err
				}
			} else {
				raw, err := e.eval(element.key)
				if err != nil {
					return value{}, err
				}
				key, err = e.arrayKey(raw, line)
				if err != nil {
					return value{}, err
				}
			}
			place, err := e.referencePlace(element.value, true)
			if err != nil {
				return value{}, err
			}
			stored := place.read()
			if stored.kind != refKind {
				stored = value{kind: refKind, ref: &referenceCell{value: dereference(stored)}}
				if err := place.replace(stored); err != nil {
					return value{}, err
				}
			}
			result.a = result.a.with(key, stored)
			continue
		}
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
		if classImplements(base.o.class, "ArrayAccess") {
			method := e.lookupMethod(base.o.class, "offsetGet")
			return e.invoke(base, method, []callArgument{{value: raw}}, node.line)
		}
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
		if err := e.warn(node.line, "Trying to access array offset on "+name); err != nil {
			return value{}, err
		}
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
	if err := e.warn(node.line, "Undefined array key "+label); err != nil {
		return value{}, err
	}
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
	name, err := e.instancePropertyName(node.name, node.nameExpr, node.line)
	if err != nil {
		return value{}, err
	}
	node.name = name
	if base.kind != objectKind {
		if err := e.warn(node.line, fmt.Sprintf("Attempt to read property \"%s\" on %s", node.name, base.typeName())); err != nil {
			return value{}, err
		}
		return value{}, nil
	}
	if base.o.enumCase != "" && node.name != "name" && node.name != "value" {
		return value{}, &runtimeError{node.line, "Error", "Undefined property: " + base.o.className() + "::$" + node.name}
	}
	if base.o.reflectedMethod != nil && node.name == "class" && base.o.reflectedMethod.owner != nil {
		return value{kind: stringKind, s: base.o.reflectedMethod.owner.name}, nil
	}
	if base.o.reflected != nil && node.name == "name" {
		return value{kind: stringKind, s: base.o.reflected.name}, nil
	}
	if base.o.reflectedProperty != nil && node.name == "name" {
		return value{kind: stringKind, s: base.o.reflectedProperty.name}, nil
	}
	if base.o.reflectedMethod != nil && node.name == "name" {
		return value{kind: stringKind, s: base.o.reflectedMethod.name}, nil
	}
	k, err := e.propertyKey(base.o, node.name, node.line)
	if err != nil {
		return value{}, err
	}
	if v, ok := base.o.properties.get(k); ok {
		return v, nil
	}
	if property := propertyForKey(base.o, k); property != nil && property.typeName != "" {
		return value{}, uninitializedPropertyError(*property, node.line)
	}
	if method := e.lookupMethod(base.o.class, "__get"); method != nil {
		return e.invoke(base, method, []callArgument{{value: value{kind: stringKind, s: node.name}}}, node.line)
	}
	if err := e.warn(node.line, "Undefined property: "+base.o.className()+"::$"+node.name); err != nil {
		return value{}, err
	}
	return value{}, nil
}

// A location resolves the delayed write path after RHS subexpressions.
type location struct {
	get func() value
	set func(value) error
	// typed marks a typed property, whose compound result is type-checked
	// before it is stored.
	typed bool
}

type locationMode uint8

const (
	locationWrite locationMode = iota
	locationReadWrite
	locationIncrement
	locationModify
	locationIndirectWrite
)

func (e *evaluator) locate(target expr, modes ...locationMode) (location, error) {
	mode := locationWrite
	if len(modes) > 0 {
		mode = modes[0]
	}
	rw := mode != locationWrite && mode != locationIndirectWrite
	switch node := target.(type) {
	case literal:
		if node.value.kind == objectKind && node.value.o != nil && classImplements(node.value.o.class, "ArrayAccess") {
			return location{get: func() value { return node.value }, set: func(value) error {
				return &runtimeError{node.line, "Error", "Cannot assign to temporary object"}
			}}, nil
		}
		return location{}, &runtimeError{node.line, "Error", "Cannot use temporary expression in write context"}
	case variable:
		if e.variables == nil {
			e.variables = make(map[string]value)
		}
		if _, exists := e.readVariable(node.token.text); rw && !exists {
			if err := e.warn(node.token.line, "Undefined variable $"+node.token.text); err != nil {
				return location{}, err
			}
		}
		return e.variableSlotAt(node.token.text, node.token.line), nil
	case staticPropertyExpr:
		loc, err := e.staticPropertyLocation(node)
		if err != nil {
			return location{}, err
		}
		if rw {
			class, err := e.resolveClassName(node.className, node.line)
			if err != nil {
				return location{}, err
			}
			if slot := class.staticProperties[node.name]; slot != nil && !slot.initialized {
				return location{}, uninitializedPropertyError(slot.decl, node.line)
			}
		}
		return loc, nil
	case propertyExpr:
		var base value
		var err error
		if v, ok := node.base.(variable); ok && v.token.text != "this" {
			base, _ = e.readVariable(v.token.text)
			if _, exists := e.readVariable(v.token.text); rw && !exists {
				if err := e.warn(v.token.line, "Undefined variable $"+v.token.text); err != nil {
					return location{}, err
				}
			}
		} else {
			base, err = e.eval(node.base)
		}
		if err != nil {
			return location{}, err
		}
		name, err := e.instancePropertyName(node.name, node.nameExpr, node.line)
		if err != nil {
			return location{}, err
		}
		node.name, node.nameExpr = name, nil
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
		if base.o.enumCase != "" {
			return location{}, &runtimeError{node.line, "Error", "Cannot modify readonly property " + base.o.className() + "::$" + node.name}
		}
		k, err := e.propertyKey(base.o, node.name, node.line)
		if err != nil {
			return location{}, err
		}
		property := propertyForKey(base.o, k)
		if (mode == locationModify || mode == locationIndirectWrite) && property != nil && property.setHook != nil && !(e.currentMethod == property.setHook && e.variables["this"].o == base.o) {
			return location{}, &runtimeError{node.line, "Error", "Indirect modification of " + property.owner.name + "::$" + property.name + " is not allowed"}
		}
		if property == nil && base.o.properties.find(k) < 0 {
			if method := e.lookupMethod(base.o.class, "__set"); method != nil && !rw {
				return location{get: func() value { return value{} }, set: func(v value) error {
					_, callErr := e.invoke(base, method, []callArgument{{value: value{kind: stringKind, s: node.name}}, {value: v}}, node.line)
					return callErr
				}}, nil
			}
		}
		if property == nil && base.o.class != nil && base.o.class.readonly {
			return location{}, &runtimeError{node.line, "Error", "Cannot create dynamic property " + base.o.class.name + "::$" + node.name}
		}
		if rw && property != nil && property.typeName != "" && base.o.properties.find(k) < 0 {
			return location{}, uninitializedPropertyError(*property, node.line)
		}
		if rw && base.o.properties.find(k) < 0 {
			if base.o.dynamicProperty(node.name, e.scope) {
				if err := e.diagnostic("Deprecated", node.line, "Creation of dynamic property "+base.o.className()+"::$"+node.name+" is deprecated"); err != nil {
					return location{}, err
				}
			}
			if err := e.warn(node.line, "Undefined property: "+base.o.className()+"::$"+node.name); err != nil {
				return location{}, err
			}
			base.o.properties = base.o.properties.with(k, value{})
		}
		return location{func() value { v, _ := base.o.properties.get(k); return v }, func(v value) error {
			creating := base.o.dynamicProperty(node.name, e.scope) && base.o.properties.find(k) < 0
			if err := e.storeObjectProperty(base.o, k, property, v, node.line); err != nil {
				return err
			}
			// Zend creates the property first, so a throwing handler keeps it.
			if creating {
				return e.diagnostic("Deprecated", node.line, "Creation of dynamic property "+base.o.className()+"::$"+node.name+" is deprecated")
			}
			return nil
		}, property != nil && property.typeName != ""}, nil
	case offsetExpr:
		if root, ok := node.base.(variable); ok && root.token.text == "GLOBALS" {
			if node.key == nil {
				return location{}, &runtimeError{node.line, "Error", "Cannot append to $GLOBALS"}
			}
			raw, err := e.eval(node.key)
			if err != nil {
				return location{}, err
			}
			key, err := e.arrayKey(raw, node.line)
			if err != nil {
				return location{}, err
			}
			name := key.text
			if !key.isString {
				name = fmt.Sprint(key.index)
			}
			if name == "GLOBALS" {
				return location{}, &runtimeError{node.line, "Error", "Cannot modify $GLOBALS array directly"}
			}
			if rw {
				if _, exists := e.rootScope[name]; !exists {
					if err := e.warn(node.line, "Undefined global variable $"+name); err != nil {
						return location{}, err
					}
				}
			}
			return e.variableScopeSlot(e.rootScope, name, node.line), nil
		}
		parentMode := locationIndirectWrite
		if rw {
			parentMode = locationModify
		}
		parent, err := e.locate(node.base, parentMode)
		if err != nil {
			return location{}, err
		}
		base := parent.get()
		if base.kind == objectKind && classImplements(base.o.class, "ArrayAccess") {
			key := value{kind: nullKind}
			if node.key != nil {
				key, err = e.eval(node.key)
				if err != nil {
					return location{}, err
				}
			}
			return location{get: func() value {
				method := e.lookupMethod(base.o.class, "offsetGet")
				v, callErr := e.invoke(base, method, []callArgument{{value: key}}, node.line)
				if callErr != nil {
					e.pendingDiagnostic = callErr
				}
				return v
			}, set: func(v value) error {
				method := e.lookupMethod(base.o.class, "offsetSet")
				_, callErr := e.invoke(base, method, []callArgument{{value: key}, {value: v}}, node.line)
				return callErr
			}}, nil
		}
		if base.kind == nullKind || base.kind == boolKind && base.i == 0 {
			if base.kind == boolKind {
				if err := e.diagnostic("Deprecated", node.line, "Automatic conversion of false to array is deprecated"); err != nil {
					return location{}, err
				}
			}
			base = arrayValue()
			if err := parent.set(base); err != nil {
				return location{}, err
			}
		}
		if base.kind == stringKind && node.key != nil {
			raw, keyErr := e.eval(node.key)
			if keyErr != nil {
				return location{}, keyErr
			}
			if raw.kind == stringKind {
				converted, trailing, ok := raw.numeric()
				if !ok || converted.kind != intKind {
					return location{}, &runtimeError{node.line, "TypeError", "Cannot access offset of type string on string"}
				}
				if trailing {
					if err := e.warn(node.line, "Illegal string offset \""+raw.s+"\""); err != nil {
						return location{}, err
					}
				}
				raw = converted
			} else if raw.kind != intKind {
				if raw.kind == arrayKind || raw.kind == objectKind {
					return location{}, &runtimeError{node.line, "TypeError", "Cannot access offset of type " + raw.typeName() + " on string"}
				}
				if err := e.warn(node.line, "String offset cast occurred"); err != nil {
					return location{}, err
				}
				if raw.kind == floatKind {
					raw = integer(floatToInt(raw.f))
				} else {
					raw = integer(raw.i)
				}
			}
			index := raw.i
			if index < 0 {
				index += int64(len(base.s))
			}
			return location{get: func() value { v, _ := e.stringOffset(parent.get(), integer(index), node.line); return v }, set: func(v value) error {
				if index < 0 {
					return &runtimeError{node.line, "Error", "Illegal string offset "}
				}
				converted, err := e.toText(v, node.line)
				if err != nil {
					return err
				}
				if converted == "" {
					return &runtimeError{node.line, "Error", "Cannot assign an empty string to a string offset"}
				}
				current := parent.get()
				if current.kind != stringKind {
					return &runtimeError{node.line, "TypeError", "Cannot use a scalar value as an array"}
				}
				bytes := []byte(current.s)
				for int64(len(bytes)) <= index {
					bytes = append(bytes, ' ')
				}
				bytes[index] = converted[0]
				return parent.set(value{kind: stringKind, s: string(bytes)})
			}}, nil
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
				err = e.warn(node.line, fmt.Sprintf("Undefined array key \"%s\"", k.text))
			} else {
				err = e.warn(node.line, fmt.Sprintf("Undefined array key %d", k.index))
			}
			if err != nil {
				return location{}, err
			}
		}
		// The RHS has already been evaluated at this point.
		if node.key == nil {
			base.a = e.arrayWith(base.a, k, value{})
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
			latest.a, err = e.arrayWithAssigned(latest.a, k, v, node.line)
			if err != nil {
				return err
			}
			return parent.set(latest)
		}, false}, nil
	}
	return location{}, fmt.Errorf("gophp: unsupported assignment target %T at %s:%d", target, e.filename, expressionLine(target))
}

// Target key subexpressions execute before the RHS, but the writes and bare
// variable fetches are delayed, as in Zend's delayed compilation of assignments.
func (e *evaluator) prepareTarget(target expr) (expr, error) {
	switch node := target.(type) {
	case variable:
		return node, nil
	case dynamicVariable:
		name, err := e.dynamicVariableName(node)
		if err != nil {
			return nil, err
		}
		return variable{token{kind: "variable", text: name, line: node.line}}, nil
	case staticPropertyExpr:
		return e.prepareStaticProperty(node)
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
		if err != nil {
			return nil, err
		}
		name, err := e.instancePropertyName(node.name, node.nameExpr, node.line)
		if err != nil {
			return nil, err
		}
		node.name, node.nameExpr = name, nil
		return node, nil
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
	if pattern, ok := node.target.(arrayLiteral); ok {
		v, err := e.eval(node.source)
		if err != nil {
			return value{}, err
		}
		if err := e.assignDestructured(pattern, v); err != nil {
			return value{}, err
		}
		return v, nil
	}
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
	return target.get(), nil
}

func (e *evaluator) assignDestructured(pattern arrayLiteral, source value) error {
	next := int64(0)
	for _, entry := range pattern.elements {
		if entry.value == nil {
			next++
			continue
		}
		key := value{kind: intKind, i: next}
		if entry.key != nil {
			v, err := e.eval(entry.key)
			if err != nil {
				return err
			}
			key = v
		} else {
			next++
		}
		line := expressionLine(entry.value)
		item := value{kind: nullKind}
		if source.kind != nullKind {
			var err error
			item, err = e.readOffset(offsetExpr{base: literal{source, line}, key: literal{key, line}, line: line})
			if err != nil {
				return err
			}
		}
		if nested, ok := entry.value.(arrayLiteral); ok {
			if err := e.assignDestructured(nested, item); err != nil {
				return err
			}
			continue
		}
		place, err := e.prepareTarget(entry.value)
		if err != nil {
			return err
		}
		target, err := e.locate(place)
		if err != nil {
			return err
		}
		if err := target.set(item); err != nil {
			return err
		}
	}
	return nil
}
