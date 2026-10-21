package php

import "strings"

func (e *evaluator) reflectionAttributes(attributes []attributeDecl, scope *classDecl, args []callArgument, line int) (value, error) {
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "name"}, {name: "flags"}}, line)
	if err != nil {
		return value{}, err
	}
	if len(args) > 2 {
		return value{}, &runtimeError{line, "ArgumentCountError", "ReflectionClass::getAttributes() expects at most 2 arguments"}
	}
	name := ""
	if present[0] && bound[0].kind != nullKind {
		v, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "ReflectionClass::getAttributes", 1, "name", "?string", bound[0])
		}
		name = strings.TrimPrefix(v.s, "\\")
	}
	flags := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "ReflectionClass::getAttributes", 2, "flags", "int", bound[1])
		}
		flags = v.i
	}
	if flags != 0 && flags != 2 {
		return value{}, &runtimeError{line, "ValueError", "ReflectionClass::getAttributes(): Argument #2 ($flags) must be a valid attribute filter flag"}
	}
	out := arrayValue()
	for i := range attributes {
		attr := &attributes[i]
		if name != "" {
			match := strings.EqualFold(attr.name, name)
			if !match && flags == 2 {
				class, resolveErr := e.resolveClassName(attr.name, line)
				if resolveErr == nil {
					for c := class; c != nil; c = c.parent {
						if strings.EqualFold(c.name, name) {
							match = true
							break
						}
					}
					if !match {
						match = classImplements(class, name)
					}
				}
			}
			if !match {
				continue
			}
		}
		e.initClasses()
		wrapped := value{kind: objectKind, o: &phpObject{class: e.classes["reflectionattribute"], reflectedAttribute: attr, attributeScope: scope, properties: &phpArray{}}}
		out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, wrapped)
	}
	return out, nil
}

func (e *evaluator) reflectionAttributeArguments(attr attributeDecl, scope *classDecl) (value, error) {
	previous := e.scope
	e.scope = scope
	defer func() { e.scope = previous }()
	args, err := e.evalArguments(attr.args)
	if err != nil {
		return value{}, err
	}
	out := arrayValue()
	for _, arg := range args {
		key := arrayKey{index: int64(len(out.a.entries))}
		if arg.named {
			key = arrayKey{isString: true, text: arg.name}
		}
		out.a = out.a.with(key, arg.value)
	}
	return out, nil
}

func (e *evaluator) reflectionAttributeInstance(attr attributeDecl, scope *classDecl) (value, error) {
	previous := e.scope
	e.scope = scope
	defer func() { e.scope = previous }()
	return e.newInstance(newObject{name: attr.name, args: attr.args, arguments: true, line: attr.line})
}
