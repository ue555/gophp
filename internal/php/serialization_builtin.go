package php

import (
	"fmt"
	"strconv"
	"strings"
)

// PHP's common N/b/i/d/s/a/O/E serialization tags. Reference tags and custom
// serialization hooks need a separate identity-aware implementation.
type serializedEnumFailure struct {
	message string
	offset  int
}

func (failure *serializedEnumFailure) Error() string { return failure.message }

func (e *evaluator) callSerialization(node functionCall) (value, error) {
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if strings.EqualFold(node.name, "serialize") {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}}, node.line)
		if err != nil {
			return value{}, err
		}
		if len(args) != 1 || !present[0] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "serialize() expects exactly 1 argument"}
		}
		var out strings.Builder
		if err := e.appendSerialized(&out, bound[0], map[*phpObject]bool{}, map[*phpArray]bool{}, node.line); err != nil {
			return value{}, err
		}
		return value{kind: stringKind, s: out.String()}, nil
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "data"}, {name: "options"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < 1 || len(args) > 2 || !present[0] {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "unserialize() expects 1 to 2 arguments"}
	}
	data, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, "unserialize", 1, "data", "string", bound[0])
	}
	if present[1] && bound[1].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "unserialize", 2, "options", "array", bound[1])
	}
	parser := serializedParser{input: data.s, e: e, line: node.line}
	result, err := parser.read(0)
	if err != nil {
		if failure, ok := err.(*serializedEnumFailure); ok {
			if err := e.warn(node.line, "unserialize(): "+failure.message); err != nil {
				return value{}, err
			}
			if err := e.warn(node.line, fmt.Sprintf("unserialize(): Error at offset %d of %d bytes", failure.offset, len(parser.input))); err != nil {
				return value{}, err
			}
			return boolean(false), nil
		}
		return value{}, err
	}
	if parser.at != len(parser.input) {
		if err := e.warn(node.line, "unserialize(): Extra data starting at offset "+strconv.Itoa(parser.at)+" of "+strconv.Itoa(len(parser.input))+" bytes"); err != nil {
			return value{}, err
		}
	}
	return result, nil
}

func (e *evaluator) appendSerialized(out *strings.Builder, v value, objects map[*phpObject]bool, arrays map[*phpArray]bool, line int) error {
	v = dereference(v)
	switch v.kind {
	case nullKind:
		out.WriteString("N;")
	case boolKind:
		if v.i != 0 {
			out.WriteString("b:1;")
		} else {
			out.WriteString("b:0;")
		}
	case intKind:
		fmt.Fprintf(out, "i:%d;", v.i)
	case floatKind:
		fmt.Fprintf(out, "d:%s;", strconv.FormatFloat(v.f, 'g', -1, 64))
	case stringKind:
		fmt.Fprintf(out, "s:%d:\"%s\";", len(v.s), v.s)
	case arrayKind:
		if arrays[v.a] {
			return &runtimeError{line, "Error", "gophp: cyclic array serialization is not implemented"}
		}
		arrays[v.a] = true
		defer delete(arrays, v.a)
		fmt.Fprintf(out, "a:%d:{", len(v.a.entries))
		for _, entry := range v.a.entries {
			key := value{kind: stringKind, s: entry.key.text}
			if !entry.key.isString {
				key = integer(entry.key.index)
			}
			if err := e.appendSerialized(out, key, objects, arrays, line); err != nil {
				return err
			}
			if err := e.appendSerialized(out, entry.value, objects, arrays, line); err != nil {
				return err
			}
		}
		out.WriteByte('}')
	case objectKind:
		if v.o == nil || v.o.class == nil {
			return &runtimeError{line, "Error", "gophp: unsupported object serialization"}
		}
		if v.o.enumCase != "" {
			caseName := v.o.class.name + ":" + v.o.enumCase
			fmt.Fprintf(out, "E:%d:\"%s\";", len(caseName), caseName)
			return nil
		}
		if objects[v.o] {
			return &runtimeError{line, "Error", "gophp: cyclic object serialization is not implemented"}
		}
		if e.lookupMethod(v.o.class, "__serialize") != nil || e.lookupMethod(v.o.class, "__sleep") != nil {
			return &runtimeError{line, "Error", "gophp: custom object serialization is not implemented"}
		}
		objects[v.o] = true
		defer delete(objects, v.o)
		arrayObject := false
		for class := v.o.class; class != nil; class = class.parent {
			if strings.EqualFold(class.name, "ArrayObject") || strings.EqualFold(class.name, "ArrayIterator") {
				arrayObject = true
				break
			}
		}
		count := 0
		if v.o.properties != nil {
			count = len(v.o.properties.entries)
		}
		if arrayObject {
			count += 4
		}
		fmt.Fprintf(out, "O:%d:\"%s\":%d:{", len(v.o.class.name), v.o.class.name, count)
		if arrayObject {
			for _, pair := range []struct {
				key  int64
				item value
			}{
				{0, integer(0)},
				{1, value{kind: arrayKind, a: copyArrayObjectData(v.o.iteratorArray)}},
				{2, arrayValue()},
				{3, value{kind: nullKind}},
			} {
				if err := e.appendSerialized(out, integer(pair.key), objects, arrays, line); err != nil {
					return err
				}
				if err := e.appendSerialized(out, pair.item, objects, arrays, line); err != nil {
					return err
				}
			}
		}
		if v.o.properties != nil {
			for _, entry := range v.o.properties.entries {
				key := value{kind: stringKind, s: entry.key.text}
				if !entry.key.isString {
					key = integer(entry.key.index)
				}
				if err := e.appendSerialized(out, key, objects, arrays, line); err != nil {
					return err
				}
				if err := e.appendSerialized(out, entry.value, objects, arrays, line); err != nil {
					return err
				}
			}
		}
		out.WriteByte('}')
	default:
		return &runtimeError{line, "Error", "gophp: serialization of this value is not implemented"}
	}
	return nil
}

type serializedParser struct {
	input string
	at    int
	e     *evaluator
	line  int
}

func (p *serializedParser) take(s string) bool {
	if strings.HasPrefix(p.input[p.at:], s) {
		p.at += len(s)
		return true
	}
	return false
}
func (p *serializedParser) number() (int64, error) {
	start := p.at
	for p.at < len(p.input) && p.input[p.at] >= '0' && p.input[p.at] <= '9' || p.at == start && p.at < len(p.input) && p.input[p.at] == '-' {
		p.at++
	}
	n, err := strconv.ParseInt(p.input[start:p.at], 10, 64)
	if err != nil {
		return 0, p.invalid()
	}
	return n, nil
}
func (p *serializedParser) invalid() error {
	return &runtimeError{p.line, "Error", "gophp: malformed serialized value"}
}
func (p *serializedParser) read(depth int) (value, error) {
	if depth > 256 || p.at >= len(p.input) {
		return value{}, p.invalid()
	}
	tag := p.input[p.at]
	p.at++
	if tag == 'N' {
		if !p.take(";") {
			return value{}, p.invalid()
		}
		return value{kind: nullKind}, nil
	}
	if !p.take(":") {
		return value{}, p.invalid()
	}
	switch tag {
	case 'b', 'i':
		n, err := p.number()
		if err != nil || !p.take(";") {
			return value{}, p.invalid()
		}
		if tag == 'b' {
			return boolean(n != 0), nil
		}
		return integer(n), nil
	case 'd':
		start := p.at
		for p.at < len(p.input) && p.input[p.at] != ';' {
			p.at++
		}
		if !p.take(";") {
			return value{}, p.invalid()
		}
		n, err := strconv.ParseFloat(p.input[start:p.at-1], 64)
		if err != nil {
			return value{}, p.invalid()
		}
		return floating(n), nil
	case 's':
		return p.readString()
	case 'E':
		serializedName, err := p.readString()
		if err != nil || serializedName.kind != stringKind {
			return value{}, p.invalid()
		}
		className, caseName, ok := strings.Cut(serializedName.s, ":")
		if !ok || className == "" || caseName == "" {
			return value{}, &serializedEnumFailure{message: fmt.Sprintf("Invalid enum name '%s' (missing colon)", serializedName.s), offset: p.at}
		}
		p.e.initClasses()
		class := p.e.classes[strings.ToLower(className)]
		if class == nil {
			if err := p.e.tryAutoload(className, p.line); err != nil {
				return value{}, err
			}
			class = p.e.classes[strings.ToLower(className)]
		}
		if class == nil || class.kind != "enum" {
			message := fmt.Sprintf("Class '%s' not found", className)
			if class != nil {
				message = fmt.Sprintf("Class '%s' is not an enum", className)
			}
			return value{}, &serializedEnumFailure{message: message, offset: p.at}
		}
		candidate, found := class.constants[caseName]
		if !found {
			return value{}, &serializedEnumFailure{message: fmt.Sprintf("Undefined constant %s::%s", className, caseName), offset: p.at}
		}
		if candidate.value.kind != objectKind || candidate.value.o == nil || candidate.value.o.enumCase != caseName {
			return value{}, &serializedEnumFailure{message: fmt.Sprintf("%s::%s is not an enum case", className, caseName), offset: p.at}
		}
		return candidate.value, nil
	case 'a', 'O':
		var class *classDecl
		if tag == 'O' {
			length, err := p.number()
			if err != nil || length < 0 || !p.take(":\"") || length > int64(len(p.input)-p.at) {
				return value{}, p.invalid()
			}
			name := value{kind: stringKind, s: p.input[p.at : p.at+int(length)]}
			p.at += int(length)
			if !p.take("\":") {
				return value{}, p.invalid()
			}
			p.e.initClasses()
			class = p.e.classes[strings.ToLower(name.s)]
			if class == nil {
				if err := p.e.tryAutoload(name.s, p.line); err != nil {
					return value{}, err
				}
				class = p.e.classes[strings.ToLower(name.s)]
			}
			if class == nil {
				return value{}, p.invalid()
			}
		}
		count, err := p.number()
		if err != nil || count < 0 || count > 100000 || !p.take(":{") {
			return value{}, p.invalid()
		}
		data := arrayValue()
		for i := int64(0); i < count; i++ {
			key, err := p.read(depth + 1)
			if err != nil {
				return value{}, err
			}
			item, err := p.read(depth + 1)
			if err != nil {
				return value{}, err
			}
			arrayKey, err := p.e.arrayKey(key, p.line)
			if err != nil {
				return value{}, err
			}
			data.a = data.a.with(arrayKey, item)
		}
		if !p.take("}") {
			return value{}, p.invalid()
		}
		if tag == 'a' {
			return data, nil
		}
		object := &phpObject{class: class, properties: data.a}
		for ancestor := class; ancestor != nil; ancestor = ancestor.parent {
			if !strings.EqualFold(ancestor.name, "ArrayObject") && !strings.EqualFold(ancestor.name, "ArrayIterator") {
				continue
			}
			if storage, ok := data.a.getRaw(arrayKey{index: 1}); ok && storage.kind == arrayKind {
				object.iteratorArray = copyArrayObjectData(storage.a)
			}
			properties := arrayValue()
			for _, entry := range data.a.entries {
				if entry.key.isString {
					properties.a = properties.a.with(entry.key, entry.value)
				}
			}
			object.properties = properties.a
			break
		}
		return value{kind: objectKind, o: object}, nil
	}
	return value{}, p.invalid()
}
func (p *serializedParser) readString() (value, error) {
	length, err := p.number()
	if err != nil || length < 0 || !p.take(":\"") || length > int64(len(p.input)-p.at) {
		return value{}, p.invalid()
	}
	s := p.input[p.at : p.at+int(length)]
	p.at += int(length)
	if !p.take("\";") {
		return value{}, p.invalid()
	}
	return value{kind: stringKind, s: s}, nil
}
