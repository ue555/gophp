package php

import (
	"fmt"
	"strings"
)

// Enum cases are allocated once per declaration. The class constant table
// retains their identity for the lifetime of this evaluator only.
func (e *evaluator) initEnumCases(decl *classDecl) error {
	seenValues := make(map[string]string)
	for _, c := range decl.enumCases {
		if _, exists := decl.constants[c.name]; exists {
			return &fatalDiagnostic{c.line, "Cannot redefine class constant " + decl.name + "::" + c.name}
		}
		if decl.enumBacking == "" && c.initial != nil {
			return &fatalDiagnostic{c.line, "A unit enum cannot have a case value"}
		}
		if decl.enumBacking != "" && c.initial == nil {
			return &fatalDiagnostic{c.line, "Case " + decl.name + "::" + c.name + " of backed enum must have a value"}
		}
		var backing value
		if c.initial != nil {
			if !constantInitializer(c.initial) {
				return &fatalDiagnostic{c.line, "Constant expression contains invalid operations"}
			}
			v, err := e.eval(c.initial)
			if err != nil {
				return err
			}
			if decl.enumBacking == "int" && v.kind != intKind || decl.enumBacking == "string" && v.kind != stringKind {
				return &fatalDiagnostic{c.line, "Enum case type " + v.typeName() + " does not match enum backing type " + decl.enumBacking}
			}
			backing = v
			key := fmt.Sprintf("%d:%s:%d", v.kind, v.s, v.i)
			if previous := seenValues[key]; previous != "" {
				decl.enumDuplicate = "Duplicate value in enum " + decl.name + " for cases " + previous + " and " + c.name
			}
			seenValues[key] = c.name
		}
		properties := &phpArray{}
		properties = properties.with(arrayKey{text: "name", isString: true}, value{kind: stringKind, s: c.name})
		if decl.enumBacking != "" {
			properties = properties.with(arrayKey{text: "value", isString: true}, backing)
		}
		caseValue := value{kind: objectKind, o: &phpObject{class: decl, properties: properties, enumCase: c.name}}
		decl.constants[c.name] = classConstantValue{value: caseValue, decl: classConstantDecl{name: c.name, visibility: "public", final: true, owner: decl, line: c.line}}
	}
	return nil
}

func (e *evaluator) callEnumStatic(class *classDecl, name string, args []expr, line int) (value, bool, error) {
	if class.kind != "enum" {
		return value{}, false, nil
	}
	switch strings.ToLower(name) {
	case "cases":
		arguments, err := e.evalArguments(args)
		if err != nil {
			return value{}, true, err
		}
		if len(arguments) != 0 {
			return value{}, true, &runtimeError{line, "ArgumentCountError", fmt.Sprintf("%s::cases() expects exactly 0 arguments, %d given", class.name, len(arguments))}
		}
		if class.enumDuplicate != "" {
			return value{}, true, &runtimeError{line, "Error", class.enumDuplicate}
		}
		result := arrayValue()
		for i, c := range class.enumCases {
			result.a = result.a.with(arrayKey{index: int64(i)}, class.constants[c.name].value)
		}
		return result, true, nil
	case "from", "tryfrom":
		if class.enumBacking == "" {
			return value{}, false, nil
		}
		arguments, err := e.evalArguments(args)
		if err != nil {
			return value{}, true, err
		}
		if len(arguments) != 1 {
			return value{}, true, &runtimeError{line, "ArgumentCountError", fmt.Sprintf("%s::%s() expects exactly 1 argument, %d given", class.name, name, len(arguments))}
		}
		if arguments[0].named && arguments[0].name != "value" {
			return value{}, true, &runtimeError{line, "Error", "Unknown named parameter $" + arguments[0].name}
		}
		v := arguments[0].value
		converted, _, ok := e.coerceType(v, class.enumBacking, false, e.strictTypes)
		if !ok {
			return value{}, true, &runtimeError{line, "TypeError", fmt.Sprintf("%s::%s(): Argument #1 ($value) must be of type %s, %s given", class.name, name, class.enumBacking, v.typeName())}
		}
		if class.enumDuplicate != "" {
			return value{}, true, &runtimeError{line, "Error", class.enumDuplicate}
		}
		for _, c := range class.enumCases {
			candidate := class.constants[c.name].value
			backing, _ := candidate.o.properties.get(arrayKey{text: "value", isString: true})
			if backing.identical(converted) {
				return candidate, true, nil
			}
		}
		if strings.EqualFold(name, "tryFrom") {
			return value{}, true, nil
		}
		if converted.kind == stringKind {
			return value{}, true, &runtimeError{line, "ValueError", fmt.Sprintf("%q is not a valid backing value for enum %s", converted.s, class.name)}
		}
		return value{}, true, &runtimeError{line, "ValueError", fmt.Sprintf("%d is not a valid backing value for enum %s", converted.i, class.name)}
	}
	return value{}, false, nil
}
