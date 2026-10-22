package php

import (
	"fmt"
	"strings"
)

// Support the backed set { ... } hooks used by Symfony Request. Other hook
// forms must fail explicitly rather than being parsed as ordinary properties.
func (p *parser) propertySetHook(property *propertyDecl) error {
	line := p.current().line
	if property.static || property.readonly || property.owner.readonly {
		return syntax(line, "gophp: static/readonly property hooks are not supported")
	}
	p.advance()
	if !p.acceptKeyword("set") {
		return syntax(p.current().line, "gophp: only backed block set hooks are supported")
	}
	hook := &methodDecl{name: "$" + property.name + "::set", visibility: property.visibility,
		owner: property.owner, line: line, version: p.version, returnType: "void", propertyHook: true,
		source:     sourceContext{p.filename, p.namespace, p.strictTypes},
		parameters: []parameter{{name: "value", typeName: property.typeName, nullable: property.nullable, line: line}}}
	if p.current().kind != "{" {
		return syntax(p.current().line, "gophp: explicit hook parameters and arrow hooks are not supported")
	}
	start := p.pos
	previous := p.currentCallable
	p.currentCallable = hook.name
	var err error
	hook.body, err = p.block()
	p.currentCallable = previous
	if err != nil {
		return err
	}
	backed := false
	for i := start; i+2 < p.pos; i++ {
		if p.tokens[i].kind == "variable" && p.tokens[i].text == "this" && p.tokens[i+1].kind == "->" && p.tokens[i+2].text == property.name {
			backed = true
		}
	}
	if !backed {
		return syntax(line, "gophp: virtual property hooks are not supported")
	}
	if err := p.expect("}"); err != nil {
		return err
	}
	property.setHook = hook
	return nil
}

func (e *evaluator) reflectionRawProperty(property *propertyDecl, method string, args []callArgument, line int) (value, error) {
	write := strings.EqualFold(method, "setRawValue")
	method = "getRawValue"
	params := []parameter{{name: "object"}}
	if write {
		method = "setRawValue"
		params = append(params, parameter{name: "value"})
	}
	name := "ReflectionProperty::" + method
	bound, _, _, err := bindCallArguments(args, params, line)
	if err != nil {
		return value{}, err
	}
	if len(args) != len(params) {
		plural := "arguments"
		if len(params) == 1 {
			plural = "argument"
		}
		return value{}, &runtimeError{line, "ArgumentCountError", fmt.Sprintf("%s() expects exactly %d %s, %d given", name, len(params), plural, len(args))}
	}
	object := bound[0]
	if object.kind != objectKind {
		return value{}, builtinTypeError(line, name, 1, "object", "object", object)
	}
	compatible := false
	for class := object.o.class; class != nil; class = class.parent {
		if equalName(class.name, property.owner.name) {
			compatible = true
			break
		}
	}
	if !compatible {
		if !write {
			return value{}, &runtimeError{line, "ReflectionException", "Given object is not an instance of the class this property was declared in"}
		}
		// Reflection writes an unrelated receiver's named property as well.
		// In that case use its own declaration (or dynamic property), not the
		// reflected class's type and slot.
		key, err := e.propertyKey(object.o, property.name, line)
		if err != nil {
			return value{}, err
		}
		return value{}, e.storeObjectProperty(object.o, key, propertyForKey(object.o, key), bound[1], line)
	}
	if property.static {
		return value{}, &runtimeError{line, "Error", "gophp: raw access to static properties is not supported"}
	}
	key := propertyStorageKey(*property)
	if !write {
		if v, ok := object.o.properties.get(key); ok {
			return v, nil
		}
		if property.typeName != "" {
			return value{}, uninitializedPropertyError(*property, line)
		}
		return value{}, nil
	}
	// Raw reflection bypasses the hook, but retains type/readonly enforcement.
	copy := *property
	copy.setHook = nil
	outerScope, outerStrict := e.scope, e.strictTypes
	e.scope, e.strictTypes = property.owner, false
	defer func() { e.scope, e.strictTypes = outerScope, outerStrict }()
	return value{}, e.storeObjectProperty(object.o, key, &copy, bound[1], line)
}
