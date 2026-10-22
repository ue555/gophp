package php

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
)

type propertyDecl struct {
	attributes       []attributeDecl
	name, visibility string
	typeName         string
	nullable         bool
	initial          expr
	static           bool
	readonly         bool
	fromTrait        bool
	traitSource      string
	line             int
	owner            *classDecl
	setHook          *methodDecl
}
type parameter struct {
	attributes       []attributeDecl
	name             string
	initial          expr
	line             int
	typeName         string
	nullable         bool
	implicitNullable bool
	variadic         bool
	byRef            bool
	promoted         bool
	visibility       string
	readonly         bool
}
type methodDecl struct {
	attributes              []attributeDecl
	returnsByReference      bool
	name, visibility        string
	final, abstract, static bool
	traitSource             string
	parameters              []parameter
	returnType              string
	returnNullable          bool
	arrow                   bool
	propertyHook            bool
	version                 Version
	body                    []statement
	line                    int
	source                  sourceContext
	owner                   *classDecl
}
type classConstantDecl struct {
	name, visibility string
	typeName         string
	nullable         bool
	initial          expr
	final            bool
	traitSource      string
	line             int
	owner            *classDecl
}
type classConstantValue struct {
	value     value
	decl      classConstantDecl
	pending   bool
	resolving bool
}
type staticPropertySlot struct {
	value       value
	decl        propertyDecl
	initialized bool
}
type classDecl struct {
	attributes       []attributeDecl
	name             string
	kind             string // class, interface, trait, or enum
	enumBacking      string
	enumCases        []enumCaseDecl
	enumDuplicate    string
	final, abstract  bool
	readonly         bool
	parentName       string
	parent           *classDecl
	interfaceNames   []string
	interfaces       []*classDecl
	traitUses        []traitUse
	properties       []propertyDecl
	methods          map[string]*methodDecl
	methodOrder      []*methodDecl
	constantDecls    []classConstantDecl
	constants        map[string]classConstantValue
	staticProperties map[string]*staticPropertySlot
	line             int
	source           sourceContext
}

type enumCaseDecl struct {
	name    string
	initial expr
	line    int
}

type traitUse struct {
	names       []string
	adaptations []traitAdaptation
}

type traitAdaptation struct {
	trait, method string
	insteadOf     []string
	alias         string
	visibility    string
	final         bool
	line          int
}

func validateReadonlyClassDeclaration(decl *classDecl) error {
	for _, property := range decl.properties {
		if property.static {
			return &fatalDiagnostic{property.line, "Readonly class " + decl.name + " cannot declare static properties"}
		}
		if property.typeName == "" {
			return &fatalDiagnostic{property.line, "Readonly property " + decl.name + "::$" + property.name + " must have type"}
		}
		if property.initial != nil {
			return &fatalDiagnostic{property.line, "Readonly property " + decl.name + "::$" + property.name + " cannot have default value"}
		}
	}
	return nil
}

func (p *parser) classDeclaration(kind string) (*classDecl, error) {
	name := p.current()
	if err := p.expect("name"); err != nil {
		return nil, err
	}
	if err := p.registerDeclaration(name, "class"); err != nil {
		return nil, err
	}
	return p.classDeclarationNamed(kind, name)
}

func (p *parser) classDeclarationNamed(kind string, name token) (*classDecl, error) {
	decl := &classDecl{name: qualify(p.namespace, name.text), kind: kind, line: name.line, methods: make(map[string]*methodDecl), source: sourceContext{p.filename, p.namespace, p.strictTypes}}
	if kind == "enum" {
		decl.final = true
	}
	if kind == "enum" && p.accept(":") {
		backing := p.current()
		if err := p.expect("name"); err != nil {
			return nil, err
		}
		decl.enumBacking = strings.ToLower(backing.text)
		if decl.enumBacking != "int" && decl.enumBacking != "string" {
			return nil, &fatalDiagnostic{backing.line, "Enum backing type must be int or string, " + backing.text + " given"}
		}
	}
	if p.acceptKeyword("extends") {
		for {
			parent, err := p.readName()
			if err != nil {
				return nil, err
			}
			resolved, _ := p.resolveName(parent.text, "class")
			if kind == "interface" {
				decl.interfaceNames = append(decl.interfaceNames, resolved)
			} else {
				decl.parentName = resolved
				break
			}
			if !p.accept(",") {
				break
			}
		}
	}
	if (kind == "class" || kind == "enum") && p.acceptKeyword("implements") {
		for {
			iface, err := p.readName()
			if err != nil {
				return nil, err
			}
			resolved, _ := p.resolveName(iface.text, "class")
			decl.interfaceNames = append(decl.interfaceNames, resolved)
			if !p.accept(",") {
				break
			}
		}
	}
	if err := p.expect("{"); err != nil {
		return nil, err
	}
	p.classDepth++
	outerClass := p.currentClass
	p.currentClass = decl
	defer func() { p.classDepth--; p.currentClass = outerClass }()
	for p.current().kind != "}" {
		attributes, err := p.attributes()
		if err != nil {
			return nil, err
		}
		if kind == "enum" && p.acceptKeyword("case") {
			caseName := p.current()
			if err := p.expect("name"); err != nil {
				return nil, err
			}
			var initial expr
			if p.accept("=") {
				var err error
				initial, err = p.expression(0, ";")
				if err != nil {
					return nil, err
				}
			}
			if err := p.expect(";"); err != nil {
				return nil, err
			}
			decl.enumCases = append(decl.enumCases, enumCaseDecl{name: caseName.text, initial: initial, line: caseName.line})
			continue
		}
		if kind != "interface" && p.acceptKeyword("use") {
			use, err := p.traitUseDeclaration()
			if err != nil {
				return nil, err
			}
			decl.traitUses = append(decl.traitUses, use)
			continue
		}
		visibility := "public"
		explicit := false
		methodFinal := false
		methodAbstract := false
		methodStatic := false
		memberReadonly := false
		for {
			if p.acceptKeyword("readonly") {
				if memberReadonly {
					return nil, unexpected(p.tokens[p.pos-1])
				}
				memberReadonly = true
				continue
			}
			if p.acceptKeyword("abstract") {
				if methodAbstract {
					return nil, unexpected(p.tokens[p.pos-1])
				}
				methodAbstract = true
				continue
			}
			if p.acceptKeyword("final") {
				if methodFinal {
					return nil, unexpected(p.tokens[p.pos-1])
				}
				methodFinal = true
				continue
			}
			if p.acceptKeyword("static") {
				if methodStatic {
					return nil, unexpected(p.tokens[p.pos-1])
				}
				methodStatic = true
				continue
			}
			if p.keyword("public") || p.keyword("protected") || p.keyword("private") || p.keyword("var") {
				modifier := strings.ToLower(p.advance().text)
				if explicit {
					return nil, unexpected(p.tokens[p.pos-1])
				}
				explicit = true
				if modifier != "var" {
					visibility = modifier
				}
				continue
			}
			break
		}
		if p.acceptKeyword("const") {
			if methodAbstract || methodStatic || memberReadonly {
				return nil, syntax(p.current().line, "unsupported class constant modifier")
			}
			constantType, constantNullable := "", false
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].kind != "=" {
				var err error
				constantType, constantNullable, err = p.typeDeclaration(true)
				if err != nil {
					return nil, err
				}
				if err = p.validMemberType(constantType, "constant"); err != nil {
					return nil, err
				}
			}
			for {
				name := p.current()
				if err := p.expect("name"); err != nil {
					return nil, err
				}
				for _, previous := range decl.constantDecls {
					if previous.name == name.text {
						return nil, &fatalDiagnostic{name.line, "Cannot redefine class constant " + decl.name + "::" + name.text}
					}
				}
				if err := p.expect("="); err != nil {
					return nil, err
				}
				initial, err := p.expression(0, ",", ";")
				if err != nil {
					return nil, err
				}
				decl.constantDecls = append(decl.constantDecls, classConstantDecl{name: name.text, visibility: visibility, typeName: constantType, nullable: constantNullable, initial: initial, final: methodFinal, line: name.line, owner: decl})
				if !p.accept(",") {
					break
				}
			}
			if err := p.expect(";"); err != nil {
				return nil, err
			}
			continue
		}
		if p.acceptKeyword("function") {
			byRefReturn := p.accept("&")
			if memberReadonly {
				return nil, &fatalDiagnostic{p.current().line, "Cannot use the readonly modifier on a method"}
			}
			name := p.current()
			if err := p.expect("name"); err != nil {
				return nil, err
			}
			lower := strings.ToLower(name.text)
			if strings.HasPrefix(lower, "__") && lower != "__construct" && lower != "__destruct" && lower != "__call" && lower != "__callstatic" && lower != "__tostring" && lower != "__get" && lower != "__set" && lower != "__isset" && lower != "__unset" && lower != "__clone" && lower != "__debuginfo" && lower != "__serialize" && lower != "__unserialize" && lower != "__sleep" && lower != "__wakeup" && lower != "__invoke" && lower != "__set_state" {
				return nil, syntax(name.line, "unsupported magic method %q", name.text)
			}
			if _, exists := decl.methods[lower]; exists {
				return nil, syntax(name.line, "duplicate method %s::%s", decl.name, name.text)
			}
			if methodAbstract && (methodFinal || visibility == "private") {
				return nil, syntax(name.line, "abstract method cannot be final or private")
			}
			if kind == "interface" && (visibility != "public" || methodFinal) {
				return nil, syntax(name.line, "unsupported interface method modifier")
			}
			method := &methodDecl{name: name.text, visibility: visibility, final: methodFinal, abstract: methodAbstract || kind == "interface", static: methodStatic, line: name.line, owner: decl, version: p.version, attributes: attributes, returnsByReference: byRefReturn}
			previousCallable := p.currentCallable
			p.currentCallable = name.text
			var err error
			method.source = sourceContext{p.filename, p.namespace, p.strictTypes}
			if method.abstract {
				method.parameters, err = p.parameterList(false)
				if err == nil {
					method.returnType, method.returnNullable, err = p.returnDeclaration()
				}
				if err == nil {
					err = p.expect(";")
				}
			} else {
				method.parameters, method.body, method.returnType, method.returnNullable, err = p.callableBody(lower == "__construct")
			}
			p.currentCallable = previousCallable
			if err != nil {
				return nil, err
			}
			if lower == "__construct" && method.returnType != "" {
				return nil, &fatalDiagnostic{name.line, "Method " + decl.name + "::__construct() cannot declare a return type"}
			}
			for _, param := range method.parameters {
				if !param.promoted {
					continue
				}
				for _, existing := range decl.properties {
					if existing.name == param.name {
						return nil, &fatalDiagnostic{param.line, "Cannot redeclare " + decl.name + "::$" + param.name}
					}
				}
				decl.properties = append(decl.properties, propertyDecl{name: param.name, visibility: param.visibility, typeName: param.typeName, nullable: param.nullable, readonly: param.readonly, line: param.line, owner: decl})
			}
			decl.methods[lower] = method
			decl.methodOrder = append(decl.methodOrder, method)
			continue
		}
		propertyType, propertyNullable := "", false
		if kind == "enum" {
			return nil, &fatalDiagnostic{p.current().line, "Enum " + decl.name + " cannot include properties"}
		}
		if p.current().kind != "variable" {
			var err error
			propertyType, propertyNullable, err = p.typeDeclaration(true)
			if err != nil {
				return nil, err
			}
			if err = p.validMemberType(propertyType, "property"); err != nil {
				return nil, err
			}
		}
		if kind == "interface" || methodFinal || methodAbstract || (!explicit && !methodStatic && !memberReadonly && propertyType == "") || p.current().kind != "variable" {
			return nil, syntax(p.current().line, "unsupported class member %q", p.current().text)
		}
		hasHooks := false
		for {
			name := p.current()
			if err := p.expect("variable"); err != nil {
				return nil, err
			}
			for _, previous := range decl.properties {
				if previous.name == name.text {
					return nil, syntax(name.line, "duplicate property %s::$%s", decl.name, name.text)
				}
			}
			property := propertyDecl{name: name.text, line: name.line, visibility: visibility, static: methodStatic, readonly: memberReadonly, owner: decl, typeName: propertyType, nullable: propertyNullable, attributes: attributes}
			if memberReadonly && methodStatic {
				return nil, &fatalDiagnostic{name.line, "Static property " + decl.name + "::$" + name.text + " cannot be readonly"}
			}
			if memberReadonly && propertyType == "" {
				return nil, &fatalDiagnostic{name.line, "Readonly property " + decl.name + "::$" + name.text + " must have type"}
			}
			if p.accept("=") {
				if memberReadonly {
					return nil, &fatalDiagnostic{name.line, "Readonly property " + decl.name + "::$" + name.text + " cannot have default value"}
				}
				var err error
				property.initial, err = p.expression(0)
				if err != nil {
					return nil, err
				}
			}
			if p.current().kind == "{" && p.version != PHP83 {
				if err := p.propertySetHook(&property); err != nil {
					return nil, err
				}
				hasHooks = true
			}
			decl.properties = append(decl.properties, property)
			if hasHooks {
				break
			}
			if !p.accept(",") {
				break
			}
		}
		if !hasHooks {
			if err := p.expect(";"); err != nil {
				return nil, err
			}
		}
	}
	return decl, p.expect("}")
}

func (p *parser) validMemberType(name, member string) error {
	if parseDeclaredType(name, false).has("callable") {
		owner := ""
		if p.currentClass != nil {
			owner = p.currentClass.name + "::"
		}
		if member == "property" {
			return &fatalDiagnostic{p.current().line, "Property " + owner + "$" + p.current().text + " cannot have type callable"}
		}
		return &fatalDiagnostic{p.current().line, "Class constant " + owner + p.current().text + " cannot have type callable"}
	}
	if name == "void" || name == "never" {
		return &fatalDiagnostic{p.current().line, fmt.Sprintf("%s cannot be used as a %s type", name, member)}
	}
	return nil
}

func (p *parser) traitUseDeclaration() (traitUse, error) {
	var use traitUse
	for {
		name, err := p.readName()
		if err != nil {
			return use, err
		}
		resolved, _ := p.resolveName(name.text, "class")
		use.names = append(use.names, resolved)
		if !p.accept(",") {
			break
		}
	}
	if p.accept(";") {
		return use, nil
	}
	if err := p.expect("{"); err != nil {
		return use, err
	}
	for !p.accept("}") {
		first, err := p.readName()
		if err != nil {
			return use, err
		}
		adapt := traitAdaptation{line: first.line, method: first.text}
		if p.accept("::") {
			adapt.trait, _ = p.resolveName(first.text, "class")
			method := p.current()
			if err := p.expect("name"); err != nil {
				return use, err
			}
			adapt.method = method.text
		}
		if p.acceptKeyword("insteadof") {
			if adapt.trait == "" {
				return use, syntax(adapt.line, "insteadof requires a trait-qualified method")
			}
			for {
				name, err := p.readName()
				if err != nil {
					return use, err
				}
				resolved, _ := p.resolveName(name.text, "class")
				adapt.insteadOf = append(adapt.insteadOf, resolved)
				if !p.accept(",") {
					break
				}
			}
		} else if p.acceptKeyword("as") {
			for p.keyword("public") || p.keyword("protected") || p.keyword("private") || p.keyword("final") {
				modifier := strings.ToLower(p.advance().text)
				if modifier == "final" {
					adapt.final = true
				} else {
					adapt.visibility = modifier
				}
			}
			if p.current().kind == "name" {
				adapt.alias = p.advance().text
			}
		} else {
			return use, unexpected(p.current(), "insteadof", "as")
		}
		if err := p.expect(";"); err != nil {
			return use, err
		}
		use.adaptations = append(use.adaptations, adapt)
	}
	return use, nil
}

type returned struct {
	value value
	line  int
}

func (*returned) Error() string { return "return" }

type callFrame struct {
	class, method string
	static        bool
	object        value
	line          int
	filename      string
	args          []value
	argKeys       []arrayKey
	callArgs      []callArgument
	boundArgs     []value
	present       []bool
	traceReady    bool
	parameters    []parameter
	locals        variableScope
}

func (e *evaluator) initClasses() {
	if e.classes == nil {
		throwable := &classDecl{name: "Throwable", kind: "interface"}
		unitEnum := &classDecl{name: "UnitEnum", kind: "interface"}
		backedEnum := &classDecl{name: "BackedEnum", kind: "interface", interfaces: []*classDecl{unitEnum}}
		countable := &classDecl{name: "Countable", kind: "interface"}
		stringable := &classDecl{name: "Stringable", kind: "interface", methods: make(map[string]*methodDecl)}
		toStringMethod := &methodDecl{name: "__toString", visibility: "public", abstract: true, returnType: "string", owner: stringable}
		stringable.methods["__tostring"] = toStringMethod
		stringable.methodOrder = []*methodDecl{toStringMethod}
		jsonSerializable := &classDecl{name: "JsonSerializable", kind: "interface", methods: make(map[string]*methodDecl)}
		jsonSerializeMethod := &methodDecl{name: "jsonSerialize", visibility: "public", abstract: true, returnType: "mixed", owner: jsonSerializable}
		jsonSerializable.methods["jsonserialize"] = jsonSerializeMethod
		jsonSerializable.methodOrder = []*methodDecl{jsonSerializeMethod}
		sessionHandlerInterface := &classDecl{name: "SessionHandlerInterface", kind: "interface", methods: make(map[string]*methodDecl)}
		for _, spec := range []struct {
			name, result string
			params       []parameter
		}{
			{"open", "bool", []parameter{{name: "path", typeName: "string"}, {name: "name", typeName: "string"}}},
			{"close", "bool", nil},
			{"read", "string|false", []parameter{{name: "id", typeName: "string"}}},
			{"write", "bool", []parameter{{name: "id", typeName: "string"}, {name: "data", typeName: "string"}}},
			{"destroy", "bool", []parameter{{name: "id", typeName: "string"}}},
			{"gc", "int|false", []parameter{{name: "max_lifetime", typeName: "int"}}},
		} {
			m := &methodDecl{name: spec.name, visibility: "public", abstract: true, parameters: spec.params, returnType: spec.result, owner: sessionHandlerInterface}
			sessionHandlerInterface.methods[spec.name] = m
			sessionHandlerInterface.methodOrder = append(sessionHandlerInterface.methodOrder, m)
		}
		countMethod := &methodDecl{name: "count", visibility: "public", abstract: true, owner: countable}
		countable.methods = map[string]*methodDecl{"count": countMethod}
		countable.methodOrder = []*methodDecl{countMethod}
		traversable := &classDecl{name: "Traversable", kind: "interface"}
		iterator := &classDecl{name: "Iterator", kind: "interface", interfaces: []*classDecl{traversable}, methods: make(map[string]*methodDecl)}
		generator := &classDecl{name: "Generator", kind: "class", interfaces: []*classDecl{iterator}, final: true}
		arrayIterator := &classDecl{name: "ArrayIterator", kind: "class", interfaces: []*classDecl{iterator, countable}}
		arrayIterator.methods = make(map[string]*methodDecl)
		for _, name := range []string{"rewind", "valid", "key", "current", "next", "count", "getArrayCopy"} {
			method := &methodDecl{name: name, visibility: "public", owner: arrayIterator}
			arrayIterator.methods[strings.ToLower(name)] = method
			arrayIterator.methodOrder = append(arrayIterator.methodOrder, method)
		}
		priorityQueue := &classDecl{name: "SplPriorityQueue", kind: "class", interfaces: []*classDecl{iterator, countable}, methods: make(map[string]*methodDecl), constants: make(map[string]classConstantValue)}
		for _, name := range []string{"insert", "extract", "top", "count", "isEmpty", "current", "key", "next", "rewind", "valid", "setExtractFlags", "getExtractFlags"} {
			method := &methodDecl{name: name, visibility: "public", owner: priorityQueue}
			priorityQueue.methods[strings.ToLower(name)] = method
			priorityQueue.methodOrder = append(priorityQueue.methodOrder, method)
		}
		for name, n := range map[string]int64{"EXTR_DATA": 1, "EXTR_PRIORITY": 2, "EXTR_BOTH": 3} {
			priorityQueue.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: priorityQueue}}
		}
		emptyIterator := &classDecl{name: "EmptyIterator", kind: "class", interfaces: []*classDecl{iterator}}
		for _, name := range []string{"current", "key", "next", "rewind", "valid"} {
			method := &methodDecl{name: name, visibility: "public", abstract: true, owner: iterator}
			iterator.methods[name] = method
			iterator.methodOrder = append(iterator.methodOrder, method)
		}
		iteratorAggregate := &classDecl{name: "IteratorAggregate", kind: "interface", interfaces: []*classDecl{traversable}, methods: make(map[string]*methodDecl)}
		getIterator := &methodDecl{name: "getIterator", visibility: "public", abstract: true, owner: iteratorAggregate}
		iteratorAggregate.methods["getiterator"] = getIterator
		iteratorAggregate.methodOrder = []*methodDecl{getIterator}
		outerIterator := &classDecl{name: "OuterIterator", kind: "interface", interfaces: []*classDecl{iterator}, methods: make(map[string]*methodDecl)}
		getInnerIterator := &methodDecl{name: "getInnerIterator", visibility: "public", abstract: true, owner: outerIterator}
		outerIterator.methods["getinneriterator"] = getInnerIterator
		outerIterator.methodOrder = []*methodDecl{getInnerIterator}
		iteratorIterator := &classDecl{name: "IteratorIterator", kind: "class", interfaces: []*classDecl{outerIterator}, methods: make(map[string]*methodDecl)}
		for _, name := range []string{"__construct", "rewind", "valid", "key", "current", "next", "getInnerIterator"} {
			method := &methodDecl{name: name, visibility: "public", owner: iteratorIterator}
			iteratorIterator.methods[strings.ToLower(name)] = method
			iteratorIterator.methodOrder = append(iteratorIterator.methodOrder, method)
		}
		limitIterator := &classDecl{name: "LimitIterator", kind: "class", parent: iteratorIterator, methods: make(map[string]*methodDecl)}
		for _, name := range []string{"__construct", "rewind", "valid", "key", "current", "next", "getPosition", "seek"} {
			method := &methodDecl{name: name, visibility: "public", owner: limitIterator}
			limitIterator.methods[strings.ToLower(name)] = method
			limitIterator.methodOrder = append(limitIterator.methodOrder, method)
		}
		appendIterator := &classDecl{name: "AppendIterator", kind: "class", parent: iteratorIterator, methods: make(map[string]*methodDecl)}
		for _, name := range []string{"__construct", "append", "rewind", "valid", "key", "current", "next", "getInnerIterator", "getIteratorIndex", "getArrayIterator"} {
			method := &methodDecl{name: name, visibility: "public", owner: appendIterator}
			appendIterator.methods[strings.ToLower(name)] = method
			appendIterator.methodOrder = append(appendIterator.methodOrder, method)
		}
		filterIterator := &classDecl{name: "FilterIterator", kind: "class", abstract: true, interfaces: []*classDecl{outerIterator}, methods: make(map[string]*methodDecl)}
		for _, name := range []string{"__construct", "accept", "getInnerIterator", "rewind", "valid", "key", "current", "next"} {
			method := &methodDecl{name: name, visibility: "public", owner: filterIterator, abstract: name == "accept"}
			filterIterator.methods[strings.ToLower(name)] = method
			filterIterator.methodOrder = append(filterIterator.methodOrder, method)
		}
		regexIterator := &classDecl{name: "RegexIterator", kind: "class", parent: filterIterator, methods: make(map[string]*methodDecl), constants: make(map[string]classConstantValue)}
		for name, method := range filterIterator.methods {
			regexIterator.methods[name] = method
		}
		for _, name := range []string{"__construct", "accept", "getRegex", "getMode"} {
			method := &methodDecl{name: name, visibility: "public", owner: regexIterator}
			regexIterator.methods[strings.ToLower(name)] = method
			regexIterator.methodOrder = append(regexIterator.methodOrder, method)
		}
		for name, n := range map[string]int64{"MATCH": 0, "GET_MATCH": 1, "ALL_MATCHES": 2, "SPLIT": 3, "REPLACE": 4, "USE_KEY": 1, "INVERT_MATCH": 2} {
			regexIterator.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: regexIterator}}
		}
		recursiveIterator := &classDecl{name: "RecursiveIterator", kind: "interface", interfaces: []*classDecl{iterator}, methods: make(map[string]*methodDecl)}
		for _, name := range []string{"hasChildren", "getChildren"} {
			method := &methodDecl{name: name, visibility: "public", abstract: true, owner: recursiveIterator}
			recursiveIterator.methods[strings.ToLower(name)] = method
			recursiveIterator.methodOrder = append(recursiveIterator.methodOrder, method)
		}
		recursiveRegexIterator := &classDecl{name: "RecursiveRegexIterator", kind: "class", parent: regexIterator, interfaces: []*classDecl{recursiveIterator}, constants: make(map[string]classConstantValue)}
		for name, c := range regexIterator.constants {
			recursiveRegexIterator.constants[name] = c
		}
		splFileInfo := &classDecl{name: "SplFileInfo", kind: "class", methods: make(map[string]*methodDecl)}
		for _, name := range []string{"__construct", "getPathname", "getFilename", "getBasename", "getExtension", "getPath", "getRealPath", "isFile", "isDir", "isLink", "__toString"} {
			method := &methodDecl{name: name, visibility: "public", owner: splFileInfo}
			splFileInfo.methods[strings.ToLower(name)] = method
			splFileInfo.methodOrder = append(splFileInfo.methodOrder, method)
		}
		directoryIterator := &classDecl{name: "DirectoryIterator", kind: "class", parent: splFileInfo, interfaces: []*classDecl{iterator}, methods: make(map[string]*methodDecl)}
		for _, name := range []string{"__construct", "current", "key", "next", "rewind", "valid", "isDot"} {
			method := &methodDecl{name: name, visibility: "public", owner: directoryIterator}
			directoryIterator.methods[strings.ToLower(name)] = method
			directoryIterator.methodOrder = append(directoryIterator.methodOrder, method)
		}
		filesystemIterator := &classDecl{name: "FilesystemIterator", kind: "class", parent: directoryIterator, methods: make(map[string]*methodDecl), constants: make(map[string]classConstantValue)}
		for _, name := range []string{"__construct", "current", "key", "next", "rewind", "valid", "getFlags", "setFlags"} {
			method := &methodDecl{name: name, visibility: "public", owner: filesystemIterator}
			filesystemIterator.methods[strings.ToLower(name)] = method
			filesystemIterator.methodOrder = append(filesystemIterator.methodOrder, method)
		}
		for name, n := range map[string]int64{"CURRENT_AS_FILEINFO": 0, "CURRENT_AS_SELF": 16, "CURRENT_AS_PATHNAME": 32, "KEY_AS_PATHNAME": 0, "KEY_AS_FILENAME": 256, "FOLLOW_SYMLINKS": 512, "SKIP_DOTS": 4096, "UNIX_PATHS": 8192, "NEW_CURRENT_AND_KEY": 256, "CURRENT_MODE_MASK": 240, "KEY_MODE_MASK": 3840, "OTHER_MODE_MASK": 12288} {
			filesystemIterator.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: filesystemIterator}}
		}
		recursiveDirectoryIterator := &classDecl{name: "RecursiveDirectoryIterator", kind: "class", parent: filesystemIterator, interfaces: []*classDecl{recursiveIterator}, methods: make(map[string]*methodDecl), constants: make(map[string]classConstantValue)}
		for name, method := range filesystemIterator.methods {
			recursiveDirectoryIterator.methods[name] = method
		}
		for name, method := range splFileInfo.methods {
			if recursiveDirectoryIterator.methods[name] == nil {
				recursiveDirectoryIterator.methods[name] = method
			}
		}
		for name, c := range filesystemIterator.constants {
			recursiveDirectoryIterator.constants[name] = c
		}
		for _, name := range []string{"__construct", "getChildren", "hasChildren", "getSubPath", "getSubPathname"} {
			method := &methodDecl{name: name, visibility: "public", owner: recursiveDirectoryIterator}
			recursiveDirectoryIterator.methods[strings.ToLower(name)] = method
			recursiveDirectoryIterator.methodOrder = append(recursiveDirectoryIterator.methodOrder, method)
		}
		recursiveIteratorIterator := &classDecl{name: "RecursiveIteratorIterator", kind: "class", interfaces: []*classDecl{outerIterator}, methods: make(map[string]*methodDecl), constants: make(map[string]classConstantValue)}
		for _, name := range []string{"__construct", "rewind", "valid", "key", "current", "next", "getDepth", "getInnerIterator", "getSubIterator"} {
			method := &methodDecl{name: name, visibility: "public", owner: recursiveIteratorIterator}
			recursiveIteratorIterator.methods[strings.ToLower(name)] = method
			recursiveIteratorIterator.methodOrder = append(recursiveIteratorIterator.methodOrder, method)
		}
		for name, n := range map[string]int64{"LEAVES_ONLY": 0, "SELF_FIRST": 1, "CHILD_FIRST": 2, "CATCH_GET_CHILD": 16} {
			recursiveIteratorIterator.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: recursiveIteratorIterator}}
		}
		arrayAccess := &classDecl{name: "ArrayAccess", kind: "interface", methods: make(map[string]*methodDecl)}
		cachingIterator := &classDecl{name: "CachingIterator", kind: "class", parent: iteratorIterator, interfaces: []*classDecl{arrayAccess, countable, stringable}, methods: make(map[string]*methodDecl), constants: make(map[string]classConstantValue)}
		for _, name := range []string{"__construct", "rewind", "valid", "key", "current", "next", "hasNext", "getInnerIterator", "getFlags", "setFlags", "getCache", "count", "offsetExists", "offsetGet", "offsetSet", "offsetUnset", "__toString"} {
			method := &methodDecl{name: name, visibility: "public", owner: cachingIterator}
			cachingIterator.methods[strings.ToLower(name)] = method
			cachingIterator.methodOrder = append(cachingIterator.methodOrder, method)
		}
		for name, n := range map[string]int64{"CALL_TOSTRING": 1, "TOSTRING_USE_KEY": 2, "TOSTRING_USE_CURRENT": 4, "TOSTRING_USE_INNER": 8, "CATCH_GET_CHILD": 16, "FULL_CACHE": 256} {
			cachingIterator.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: cachingIterator}}
		}
		arrayObject := &classDecl{name: "ArrayObject", kind: "class", interfaces: []*classDecl{iteratorAggregate, arrayAccess, countable}, methods: make(map[string]*methodDecl)}
		for _, name := range []string{"__construct", "offsetExists", "offsetGet", "offsetSet", "offsetUnset", "count", "getIterator", "getArrayCopy", "exchangeArray", "ksort", "append"} {
			method := &methodDecl{name: name, visibility: "public", owner: arrayObject}
			switch name {
			case "offsetExists":
				method.parameters = []parameter{{name: "offset", typeName: "mixed"}}
				method.returnType = "bool"
			case "offsetGet":
				method.parameters = []parameter{{name: "offset", typeName: "mixed"}}
				method.returnType = "mixed"
			case "offsetSet":
				method.parameters = []parameter{{name: "offset", typeName: "mixed"}, {name: "value", typeName: "mixed"}}
				method.returnType = "void"
			case "offsetUnset":
				method.parameters = []parameter{{name: "offset", typeName: "mixed"}}
				method.returnType = "void"
			case "count":
				method.returnType = "int"
			}
			arrayObject.methods[strings.ToLower(name)] = method
			arrayObject.methodOrder = append(arrayObject.methodOrder, method)
		}
		weakMap := &classDecl{name: "WeakMap", kind: "class", final: true, interfaces: []*classDecl{arrayAccess, countable, iteratorAggregate}, methods: make(map[string]*methodDecl)}
		objectStorage := &classDecl{name: "SplObjectStorage", kind: "class", interfaces: []*classDecl{arrayAccess, countable}, methods: make(map[string]*methodDecl)}
		for _, name := range []string{"attach", "detach", "contains", "offsetExists", "offsetGet", "offsetSet", "offsetUnset", "count"} {
			method := &methodDecl{name: name, visibility: "public", owner: objectStorage}
			objectStorage.methods[strings.ToLower(name)] = method
			objectStorage.methodOrder = append(objectStorage.methodOrder, method)
		}
		for _, name := range []string{"offsetExists", "offsetGet", "offsetSet", "offsetUnset", "count", "getIterator"} {
			method := &methodDecl{name: name, visibility: "public", owner: weakMap}
			weakMap.methods[strings.ToLower(name)] = method
			weakMap.methodOrder = append(weakMap.methodOrder, method)
		}
		// The DateTime extension's declaration is needed to load Carbon's
		// interface graph. Date/time operations remain a separate implementation.
		dateTimeInterface := &classDecl{name: "DateTimeInterface", kind: "interface", constants: make(map[string]classConstantValue)}
		for name, format := range map[string]string{"ATOM": `Y-m-d\TH:i:sP`, "RFC3339": `Y-m-d\TH:i:sP`, "W3C": `Y-m-d\TH:i:sP`, "RFC3339_EXTENDED": `Y-m-d\TH:i:s.vP`, "RFC2822": `D, d M Y H:i:s O`} {
			dateTimeInterface.constants[name] = classConstantValue{value: value{kind: stringKind, s: format}, decl: classConstantDecl{name: name, visibility: "public", owner: dateTimeInterface}}
		}
		dateTime := &classDecl{name: "DateTime", kind: "class", interfaces: []*classDecl{dateTimeInterface}}
		dateTimeImmutable := &classDecl{name: "DateTimeImmutable", kind: "class", interfaces: []*classDecl{dateTimeInterface}}
		dateTimeZone := &classDecl{name: "DateTimeZone", kind: "class", methods: make(map[string]*methodDecl)}
		datePeriod := &classDecl{name: "DatePeriod", kind: "class", interfaces: []*classDecl{iteratorAggregate}, methods: make(map[string]*methodDecl), constants: make(map[string]classConstantValue)}
		datePeriod.methods["__construct"] = &methodDecl{name: "__construct", visibility: "public", owner: datePeriod}
		datePeriod.methods["getiterator"] = &methodDecl{name: "getIterator", visibility: "public", owner: datePeriod, returnType: "Traversable"}
		datePeriod.methodOrder = []*methodDecl{datePeriod.methods["__construct"], datePeriod.methods["getiterator"]}
		for name, n := range map[string]int64{"EXCLUDE_START_DATE": 1, "INCLUDE_END_DATE": 2} {
			datePeriod.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: datePeriod}}
		}
		pdo := &classDecl{name: "PDO", kind: "class", constants: make(map[string]classConstantValue)}
		pdo.methods = make(map[string]*methodDecl)
		for _, methodName := range []string{"connect", "__construct", "getAttribute", "setAttribute", "prepare", "query", "exec", "beginTransaction", "inTransaction", "commit", "rollBack", "lastInsertId", "quote", "errorCode", "errorInfo"} {
			method := &methodDecl{name: methodName, visibility: "public", owner: pdo, static: methodName == "connect"}
			pdo.methods[strings.ToLower(methodName)] = method
			pdo.methodOrder = append(pdo.methodOrder, method)
		}
		pdoStatement := &classDecl{name: "PDOStatement", kind: "class", methods: make(map[string]*methodDecl)}
		for _, methodName := range []string{"bindValue", "execute", "fetch", "fetchAll", "fetchColumn", "rowCount", "columnCount", "closeCursor", "setFetchMode", "errorCode", "errorInfo"} {
			method := &methodDecl{name: methodName, visibility: "public", owner: pdoStatement}
			pdoStatement.methods[strings.ToLower(methodName)] = method
			pdoStatement.methodOrder = append(pdoStatement.methodOrder, method)
		}
		for name, n := range map[string]int64{"FETCH_DEFAULT": 0, "FETCH_ASSOC": 2, "FETCH_NUM": 3, "FETCH_BOTH": 4, "FETCH_OBJ": 5, "FETCH_COLUMN": 7, "PARAM_NULL": 0, "PARAM_INT": 1, "PARAM_STR": 2, "PARAM_LOB": 3, "PARAM_BOOL": 5, "ATTR_CLIENT_VERSION": 5, "ATTR_SERVER_VERSION": 4, "ATTR_DRIVER_NAME": 16, "ATTR_TIMEOUT": 2, "ATTR_CASE": 8, "CASE_NATURAL": 0, "CASE_UPPER": 1, "CASE_LOWER": 2, "ATTR_ERRMODE": 3, "ERRMODE_SILENT": 0, "ERRMODE_WARNING": 1, "ERRMODE_EXCEPTION": 2, "ATTR_ORACLE_NULLS": 11, "NULL_NATURAL": 0, "NULL_EMPTY_STRING": 1, "NULL_TO_STRING": 2, "ATTR_STRINGIFY_FETCHES": 17, "ATTR_DEFAULT_FETCH_MODE": 19, "ATTR_EMULATE_PREPARES": 20} {
			pdo.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: pdo}}
		}
		dateTimeZone.methods["__construct"] = &methodDecl{name: "__construct", visibility: "public", owner: dateTimeZone}
		dateTimeZone.methodOrder = []*methodDecl{dateTimeZone.methods["__construct"]}
		fiber := &classDecl{name: "Fiber", kind: "class", final: true, methods: make(map[string]*methodDecl)}
		fiber.methods["getcurrent"] = &methodDecl{name: "getCurrent", visibility: "public", static: true, owner: fiber}
		fiber.methodOrder = []*methodDecl{fiber.methods["getcurrent"]}
		reflectionFunctionAbstract := &classDecl{name: "ReflectionFunctionAbstract", kind: "class", abstract: true}
		reflectionFunction := &classDecl{name: "ReflectionFunction", kind: "class", parent: reflectionFunctionAbstract, methods: make(map[string]*methodDecl)}
		reflectionFunction.methods["__construct"] = &methodDecl{name: "__construct", visibility: "public", owner: reflectionFunction}
		reflectionFunction.methodOrder = []*methodDecl{reflectionFunction.methods["__construct"]}
		reflectionMethod := &classDecl{name: "ReflectionMethod", kind: "class", parent: reflectionFunctionAbstract, constants: make(map[string]classConstantValue)}
		for name, n := range map[string]int64{"IS_PUBLIC": 1, "IS_PROTECTED": 2, "IS_PRIVATE": 4, "IS_STATIC": 16, "IS_FINAL": 32, "IS_ABSTRACT": 64} {
			reflectionMethod.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: reflectionMethod}}
		}
		reflectionProperty := &classDecl{name: "ReflectionProperty", kind: "class", constants: make(map[string]classConstantValue)}
		for name, n := range map[string]int64{"IS_PUBLIC": 1, "IS_PROTECTED": 2, "IS_PRIVATE": 4, "IS_STATIC": 16, "IS_READONLY": 128} {
			reflectionProperty.constants[name] = classConstantValue{value: integer(n), decl: classConstantDecl{name: name, visibility: "public", owner: reflectionProperty}}
		}
		for _, cls := range []*classDecl{dateTime, dateTimeImmutable} {
			cls.methods = make(map[string]*methodDecl)
			for _, name := range []string{"__construct", "format", "getTimestamp", "getTimezone", "setTimezone", "getLastErrors", "createFromFormat", "add", "sub"} {
				method := &methodDecl{name: name, visibility: "public", owner: cls}
				switch name {
				case "add", "sub":
					method.parameters = []parameter{{name: "interval", typeName: "DateInterval"}}
				case "format":
					method.parameters = []parameter{{name: "format", typeName: "string"}}
					method.returnType = "string"
				case "getTimestamp":
					method.returnType = "int"
				case "getTimezone":
					method.returnType = "DateTimeZone|false"
				case "setTimezone":
					method.parameters = []parameter{{name: "timezone", typeName: "DateTimeZone"}}
					method.returnType = cls.name
				case "getLastErrors":
					method.static = true
					method.returnType = "array|false"
				case "createFromFormat":
					method.static = true
					method.parameters = []parameter{{name: "format", typeName: "string"}, {name: "datetime", typeName: "string"}, {name: "timezone", typeName: "DateTimeZone", nullable: true, initial: literal{value: value{kind: nullKind}}}}
				}
				cls.methods[strings.ToLower(name)] = method
				cls.methodOrder = append(cls.methodOrder, method)
			}
		}
		dateInterval := &classDecl{name: "DateInterval", kind: "class", methods: make(map[string]*methodDecl)}
		constructor := &methodDecl{name: "__construct", visibility: "public", owner: dateInterval, parameters: []parameter{{name: "duration", typeName: "string"}}}
		dateInterval.methods["__construct"] = constructor
		dateInterval.methodOrder = []*methodDecl{constructor}
		jsonMethod := &methodDecl{name: "jsonSerialize", visibility: "public", abstract: true, owner: jsonSerializable, returnType: "mixed"}
		jsonSerializable.methods["jsonserialize"] = jsonMethod
		jsonSerializable.methodOrder = []*methodDecl{jsonMethod}
		for _, name := range []string{"offsetExists", "offsetGet", "offsetSet", "offsetUnset"} {
			method := &methodDecl{name: name, visibility: "public", abstract: true, owner: arrayAccess}
			method.parameters = []parameter{{name: "offset", typeName: "mixed"}}
			switch name {
			case "offsetExists":
				method.returnType = "bool"
			case "offsetGet":
				method.returnType = "mixed"
			case "offsetSet":
				method.parameters = append(method.parameters, parameter{name: "value", typeName: "mixed"})
				method.returnType = "void"
			case "offsetUnset":
				method.returnType = "void"
			}
			arrayAccess.methods[strings.ToLower(name)] = method
			arrayAccess.methodOrder = append(arrayAccess.methodOrder, method)
		}
		exception := &classDecl{name: "Exception", kind: "class", interfaces: []*classDecl{throwable}}
		exception.methods = map[string]*methodDecl{"__tostring": {name: "__toString", visibility: "public", returnType: "string", owner: exception}}
		exception.properties = []propertyDecl{
			{name: "message", visibility: "protected", initial: literal{value: value{kind: stringKind}}, owner: exception},
			{name: "code", visibility: "protected", initial: literal{value: integer(0)}, owner: exception},
		}
		logicException := &classDecl{name: "LogicException", kind: "class", parent: exception}
		runtimeException := &classDecl{name: "RuntimeException", kind: "class", parent: exception}
		logicException.properties = append([]propertyDecl(nil), exception.properties...)
		runtimeException.properties = append([]propertyDecl(nil), exception.properties...)
		pdoException := &classDecl{name: "PDOException", kind: "class", parent: runtimeException}
		pdoException.properties = append(append([]propertyDecl(nil), runtimeException.properties...), propertyDecl{name: "errorInfo", visibility: "public", initial: literal{value: value{kind: nullKind}}, owner: pdoException})
		errorException := &classDecl{name: "ErrorException", kind: "class", parent: exception}
		dateException := &classDecl{name: "DateException", kind: "class", parent: exception}
		badFunctionCall := &classDecl{name: "BadFunctionCallException", kind: "class", parent: logicException}
		baseError := &classDecl{name: "Error", kind: "class", interfaces: []*classDecl{throwable}}
		typeError := &classDecl{name: "TypeError", kind: "class", parent: baseError}
		baseError.methods = map[string]*methodDecl{"__tostring": {name: "__toString", visibility: "public", returnType: "string", owner: baseError}}
		compileError := &classDecl{name: "CompileError", kind: "class", parent: baseError}
		domNodeClass := &classDecl{name: "DOMNode", kind: "class"}
		domDocumentClass := &classDecl{name: "DOMDocument", kind: "class", parent: domNodeClass}
		domElementClass := &classDecl{name: "DOMElement", kind: "class", parent: domNodeClass}
		domTextClass := &classDecl{name: "DOMText", kind: "class", parent: domNodeClass}
		domCommentClass := &classDecl{name: "DOMComment", kind: "class", parent: domNodeClass}
		domNodeListClass := &classDecl{name: "DOMNodeList", kind: "class", interfaces: []*classDecl{iterator, countable}}
		e.classes = map[string]*classDecl{
			"domnode": domNodeClass, "domdocument": domDocumentClass, "domelement": domElementClass,
			"domtext": domTextClass, "domcomment": domCommentClass, "domnodelist": domNodeListClass,
			"stdclass": {name: "stdClass"}, "closure": {name: "Closure", final: true},
			"throwable": throwable, "exception": exception, "error": baseError,
			"logicexception": logicException, "runtimeexception": runtimeException,
			"pdoexception":                         pdoException,
			"errorexception":                       errorException,
			"dateexception":                        dateException,
			"datemalformedintervalstringexception": {name: "DateMalformedIntervalStringException", kind: "class", parent: dateException},
			"reflectionexception":                  {name: "ReflectionException", kind: "class", parent: exception},
			"badfunctioncallexception":             badFunctionCall,
			"badmethodcallexception":               {name: "BadMethodCallException", kind: "class", parent: badFunctionCall},
			"domainexception":                      {name: "DomainException", kind: "class", parent: logicException},
			"invalidargumentexception":             {name: "InvalidArgumentException", kind: "class", parent: logicException},
			"lengthexception":                      {name: "LengthException", kind: "class", parent: logicException},
			"outofrangeexception":                  {name: "OutOfRangeException", kind: "class", parent: logicException},
			"outofboundsexception":                 {name: "OutOfBoundsException", kind: "class", parent: runtimeException},
			"overflowexception":                    {name: "OverflowException", kind: "class", parent: runtimeException},
			"rangeexception":                       {name: "RangeException", kind: "class", parent: runtimeException},
			"underflowexception":                   {name: "UnderflowException", kind: "class", parent: runtimeException},
			"unexpectedvalueexception":             {name: "UnexpectedValueException", kind: "class", parent: runtimeException},
			"unitenum":                             unitEnum, "backedenum": backedEnum,
			"countable":                  countable,
			"stringable":                 stringable,
			"traversable":                traversable,
			"iterator":                   iterator,
			"generator":                  generator,
			"arrayiterator":              arrayIterator,
			"arrayobject":                arrayObject,
			"pdo":                        pdo,
			"dateperiod":                 datePeriod,
			"pdostatement":               pdoStatement,
			"splpriorityqueue":           priorityQueue,
			"emptyiterator":              emptyIterator,
			"reflectionclass":            {name: "ReflectionClass", kind: "class"},
			"reflectionenum":             {name: "ReflectionEnum", kind: "class"},
			"reflectionfunctionabstract": reflectionFunctionAbstract,
			"sessionhandlerinterface":    sessionHandlerInterface,
			"reflectionfunction":         reflectionFunction,
			"reflectionmethod":           reflectionMethod,
			"reflectionproperty":         reflectionProperty,
			"reflectionparameter":        {name: "ReflectionParameter", kind: "class"},
			"reflectionattribute":        {name: "ReflectionAttribute", kind: "class", constants: map[string]classConstantValue{"IS_INSTANCEOF": {value: integer(2), decl: classConstantDecl{name: "IS_INSTANCEOF", visibility: "public"}}}},
			"reflectionnamedtype":        {name: "ReflectionNamedType", kind: "class"},
			"iteratoraggregate":          iteratorAggregate,
			"outeriterator":              outerIterator,
			"filteriterator":             filterIterator,
			"regexiterator":              regexIterator,
			"recursiveregexiterator":     recursiveRegexIterator,
			"recursiveiterator":          recursiveIterator,
			"splfileinfo":                splFileInfo,
			"directoryiterator":          directoryIterator,
			"filesystemiterator":         filesystemIterator,
			"recursivedirectoryiterator": recursiveDirectoryIterator,
			"recursiveiteratoriterator":  recursiveIteratorIterator,
			"iteratoriterator":           iteratorIterator,
			"limititerator":              limitIterator,
			"appenditerator":             appendIterator,
			"cachingiterator":            cachingIterator,
			"arrayaccess":                arrayAccess,
			"weakmap":                    weakMap,
			"splobjectstorage":           objectStorage,
			"jsonserializable":           jsonSerializable,
			"datetimeinterface":          dateTimeInterface,
			"datetime":                   dateTime,
			"fiber":                      fiber,
			"datetimeimmutable":          dateTimeImmutable,
			"datetimezone":               dateTimeZone,
			"dateinterval":               dateInterval,
			"typeerror":                  typeError,
			"compileerror":               compileError,
			"parseerror":                 {name: "ParseError", kind: "class", parent: compileError},
			"argumentcounterror":         {name: "ArgumentCountError", kind: "class", parent: typeError},
			"requestparsebodyexception":  {name: "RequestParseBodyException", kind: "class", parent: exception},
			"valueerror":                 {name: "ValueError", kind: "class", parent: baseError},
			"divisionbyzeroerror":        {name: "DivisionByZeroError", kind: "class", parent: baseError},
			"unhandledmatcherror":        {name: "UnhandledMatchError", kind: "class", parent: baseError},
		}
	}
}

func (e *evaluator) declareClass(decl *classDecl) error {
	e.initClasses()
	// PHP supplies the Stringable return contract for legacy __toString
	// declarations during class compilation.
	if method := decl.methods["__tostring"]; method != nil && method.returnType == "" {
		method.returnType = "string"
	}
	key := strings.ToLower(decl.name)
	if _, exists := e.classes[key]; exists {
		return &runtimeError{decl.line, "", "Cannot declare class " + decl.name + ", because the name is already in use"}
	}
	if decl.final && decl.abstract {
		return &runtimeError{decl.line, "", "Cannot use the final modifier on an abstract class"}
	}
	if decl.readonly {
		if err := validateReadonlyClassDeclaration(decl); err != nil {
			return err
		}
	}
	if decl.kind == "enum" {
		unit := e.classes["unitenum"]
		decl.interfaces = append(decl.interfaces, unit)
		if decl.enumBacking != "" {
			decl.interfaces = append(decl.interfaces, e.classes["backedenum"])
		}
	} else if decl.kind == "class" {
		for _, name := range decl.interfaceNames {
			if strings.EqualFold(name, "UnitEnum") || strings.EqualFold(name, "BackedEnum") {
				return &fatalDiagnostic{decl.line, "Non-enum class " + decl.name + " cannot implement interface " + name}
			}
		}
	}
	for _, name := range decl.interfaceNames {
		iface := e.classes[strings.ToLower(name)]
		if iface == nil {
			if err := e.tryAutoload(name, decl.line); err != nil {
				return err
			}
			iface = e.classes[strings.ToLower(name)]
		}
		if iface == nil || iface.kind != "interface" {
			return &runtimeError{decl.line, "Error", "Interface \"" + name + "\" not found"}
		}
		if decl.kind == "class" && (equalName(iface.name, "UnitEnum") || equalName(iface.name, "BackedEnum") || classImplements(iface, "UnitEnum")) {
			return &fatalDiagnostic{decl.line, "Non-enum class " + decl.name + " cannot implement interface " + iface.name}
		}
		if decl.kind == "enum" && decl.enumBacking == "" && (equalName(iface.name, "BackedEnum") || classImplements(iface, "BackedEnum")) {
			return &fatalDiagnostic{decl.line, "Non-backed enum " + decl.name + " cannot implement interface " + iface.name}
		}
		for _, existing := range decl.interfaces {
			if existing == iface {
				if decl.kind == "interface" {
					return &fatalDiagnostic{decl.line, "Interface " + decl.name + " cannot extend previously extended interface " + iface.name}
				}
				return &fatalDiagnostic{decl.line, "Class " + decl.name + " cannot implement previously implemented interface " + iface.name}
			}
		}
		decl.interfaces = append(decl.interfaces, iface)
	}
	if decl.kind == "class" && !decl.abstract && classImplements(decl, "Traversable") && !classImplements(decl, "Iterator") && !classImplements(decl, "IteratorAggregate") {
		return &fatalDiagnostic{decl.line, "Class " + decl.name + " must implement either Iterator or IteratorAggregate to be Traversable"}
	}
	if decl.kind == "interface" {
		methods := make(map[string]*methodDecl)
		for _, iface := range decl.interfaces {
			for _, method := range interfaceContracts(iface) {
				name := strings.ToLower(method.name)
				if existing := methods[name]; existing != nil && (existing.static != method.static || !e.compatibleMethodDeclaration(existing, method)) {
					if err := e.methodCompatibilityDiagnostic(existing, method, existing.line); err != nil {
						return err
					}
				}
				if methods[name] == nil {
					methods[name] = method
				}
			}
		}
		for _, method := range decl.methodOrder {
			name := strings.ToLower(method.name)
			for _, iface := range decl.interfaces {
				for _, inherited := range interfaceContracts(iface) {
					if strings.EqualFold(name, inherited.name) && (method.static != inherited.static || !e.compatibleMethodDeclaration(method, inherited)) {
						if err := e.methodCompatibilityDiagnostic(method, inherited, method.line); err != nil {
							return err
						}
					}
				}
			}
			methods[name] = method
		}
		decl.methods = methods
		e.classes[key] = decl
		if err := e.initClassMembers(decl); err != nil {
			delete(e.classes, key)
			return err
		}
		return nil
	}
	if err := e.composeTraits(decl); err != nil {
		return err
	}
	if decl.kind == "enum" {
		for name, method := range decl.methods {
			if name == "__construct" || name == "cases" || method.traitSource == "" && decl.enumBacking != "" && (name == "from" || name == "tryfrom") {
				return &fatalDiagnostic{decl.line, "Cannot redeclare " + decl.name + "::" + method.name + "()"}
			}
		}
	}
	if decl.kind == "enum" && len(decl.properties) > 0 {
		return &fatalDiagnostic{decl.line, "Enum " + decl.name + " cannot include properties"}
	}
	if decl.readonly {
		for _, property := range decl.properties {
			if property.fromTrait && !property.readonly {
				return &fatalDiagnostic{decl.line, "Readonly class " + decl.name + " cannot use trait with a non-readonly property " + property.traitSource + "::$" + property.name}
			}
			if property.static || property.typeName == "" || property.initial != nil {
				return &fatalDiagnostic{decl.line, "Readonly class " + decl.name + " cannot use incompatible trait property $" + property.name}
			}
		}
	}
	if decl.kind == "trait" {
		e.classes[key] = decl
		if err := e.initClassMembers(decl); err != nil {
			delete(e.classes, key)
			return err
		}
		return nil
	}
	if decl.parentName != "" {
		parent := e.classes[strings.ToLower(decl.parentName)]
		if parent == nil {
			if err := e.tryAutoload(decl.parentName, decl.line); err != nil {
				return err
			}
			parent = e.classes[strings.ToLower(decl.parentName)]
		}
		if parent == nil {
			return &runtimeError{decl.line, "Error", "Class \"" + decl.parentName + "\" not found"}
		}
		if parent.kind == "interface" {
			return &fatalDiagnostic{decl.line, "Class " + decl.name + " cannot extend interface " + parent.name}
		}
		if parent.final {
			return &fatalDiagnostic{decl.line, "Class " + decl.name + " cannot extend final class " + parent.name}
		}
		if decl.readonly && !parent.readonly {
			return &fatalDiagnostic{decl.line, "Readonly class " + decl.name + " cannot extend non-readonly class " + parent.name}
		}
		if !decl.readonly && parent.readonly {
			return &fatalDiagnostic{decl.line, "Non-readonly class " + decl.name + " cannot extend readonly class " + parent.name}
		}
		decl.parent = parent
		properties := append([]propertyDecl(nil), parent.properties...)
		for _, property := range decl.properties {
			replaced := false
			for i, inherited := range parent.properties {
				if property.name == inherited.name && inherited.visibility != "private" {
					if (property.readonly || decl.readonly) != (inherited.readonly || inherited.owner.readonly) {
						return &fatalDiagnostic{property.line, "Cannot change readonly modifier of property " + decl.name + "::$" + property.name}
					}
					if !e.typeAcceptsScoped(property.typeName, property.nullable, property.owner, inherited.typeName, inherited.nullable, inherited.owner) || !e.typeAcceptsScoped(inherited.typeName, inherited.nullable, inherited.owner, property.typeName, property.nullable, property.owner) {
						return &fatalDiagnostic{property.line, fmt.Sprintf("Type of %s::$%s must be %s (as in class %s)", decl.name, property.name, declaredTypeText(inherited.typeName, inherited.nullable), inherited.owner.name)}
					}
					if property.fromTrait && e.version != PHP85 {
						compatible, err := e.compatibleProperty(property, inherited)
						if err != nil {
							return err
						}
						if !compatible {
							return &fatalDiagnostic{decl.line, fmt.Sprintf("%s and %s define the same property ($%s) in the composition of %s. However, the definition differs and is considered incompatible. Class was composed", parent.name, property.traitSource, property.name, decl.name)}
						}
					}
					if property.static != inherited.static {
						return &fatalDiagnostic{property.line, fmt.Sprintf("Cannot redeclare %s property %s::$%s as %s", map[bool]string{true: "static", false: "non static"}[inherited.static], decl.name, property.name, map[bool]string{true: "static", false: "non static"}[property.static])}
					}
					if visibilityRank(property.visibility) < visibilityRank(inherited.visibility) {
						return &fatalDiagnostic{property.line, fmt.Sprintf("Access level to %s::$%s must be %s (as in class %s)%s", decl.name, property.name, inherited.visibility, inherited.owner.name, visibilitySuffix(inherited.visibility))}
					}
					properties[i] = property
					replaced = true
					break
				}
			}
			if !replaced {
				properties = append(properties, property)
			}
		}
		decl.properties = properties
		methods := make(map[string]*methodDecl, len(parent.methods)+len(decl.methods))
		for name, method := range parent.methods {
			methods[name] = method
		}
		for name, method := range decl.methods {
			if isThrowableClass(parent) && (name == "getmessage" || name == "getcode" || name == "getfile" || name == "getline" || name == "gettrace" || name == "gettraceasstring" || name == "getprevious") {
				base := "Exception"
				for c := parent; c != nil; c = c.parent {
					if equalName(c.name, "Error") {
						base = "Error"
						break
					}
					if equalName(c.name, "Exception") {
						break
					}
				}
				return &fatalDiagnostic{method.line, fmt.Sprintf("Cannot override final method %s::%s()", base, method.name)}
			}
			if inherited := parent.methods[name]; inherited != nil && inherited.final && (inherited.visibility != "private" || name == "__construct") {
				return &fatalDiagnostic{method.line, fmt.Sprintf("Cannot override final method %s::%s()", inherited.owner.name, inherited.name)}
			}
			if inherited := parent.methods[name]; inherited != nil && inherited.visibility != "private" && method.static != inherited.static {
				return &fatalDiagnostic{method.line, fmt.Sprintf("Cannot make %s method %s::%s() %s in class %s", map[bool]string{true: "static", false: "non static"}[inherited.static], inherited.owner.name, inherited.name, map[bool]string{true: "non static", false: "static"}[inherited.static], decl.name)}
			}
			if inherited := parent.methods[name]; inherited != nil && inherited.visibility != "private" && name != "__construct" && visibilityRank(method.visibility) < visibilityRank(inherited.visibility) {
				return &fatalDiagnostic{method.line, fmt.Sprintf("Access level to %s::%s() must be %s (as in class %s)%s", decl.name, method.name, inherited.visibility, inherited.owner.name, visibilitySuffix(inherited.visibility))}
			}
			if inherited := parent.methods[name]; inherited != nil && inherited.visibility != "private" && name != "__construct" && !e.compatibleMethodDeclaration(method, inherited) {
				if err := e.methodCompatibilityDiagnostic(method, inherited, method.line); err != nil {
					return err
				}
			}
			methods[name] = method
		}
		decl.methods = methods
	}
	if !isThrowableClass(decl.parent) {
		for _, iface := range decl.interfaces {
			if equalName(iface.name, "Throwable") || classImplements(iface, "Throwable") {
				return &fatalDiagnostic{decl.line, "Class " + decl.name + " cannot implement interface Throwable, extend Exception or Error instead"}
			}
		}
	}
	if !decl.abstract {
		for _, method := range decl.methods {
			if method.abstract {
				remaining := "methods"
				if e.version == PHP85 {
					remaining = "method"
				}
				return &fatalDiagnostic{decl.line, fmt.Sprintf("Class %s contains 1 abstract method and must therefore be declared abstract or implement the remaining %s (%s::%s)", decl.name, remaining, method.owner.name, method.name)}
			}
		}
	}
	// A concrete child must also satisfy contracts deferred by an abstract
	// parent. Those interfaces are inherited even when implements is omitted.
	requiredInterfaces := append([]*classDecl(nil), decl.interfaces...)
	for parent := decl.parent; parent != nil; parent = parent.parent {
		requiredInterfaces = append(requiredInterfaces, parent.interfaces...)
	}
	for _, iface := range requiredInterfaces {
		for _, required := range interfaceContracts(iface) {
			name := strings.ToLower(required.name)
			method := e.lookupMethod(decl, name)
			if method == nil || method.abstract {
				if decl.abstract {
					continue
				}
				remaining := "methods"
				if e.version == PHP85 {
					remaining = "method"
				}
				return &fatalDiagnostic{decl.line, fmt.Sprintf("Class %s contains 1 abstract method and must therefore be declared abstract or implement the remaining %s (%s::%s)", decl.name, remaining, required.owner.name, required.name)}
			}
			if method.visibility != "public" {
				return &fatalDiagnostic{method.line, fmt.Sprintf("Access level to %s::%s() must be public (as in class %s)", decl.name, method.name, required.owner.name)}
			}
			if method.static != required.static {
				return &fatalDiagnostic{method.line, fmt.Sprintf("Cannot make %s method %s::%s() %s in class %s", map[bool]string{true: "static", false: "non static"}[required.static], required.owner.name, required.name, map[bool]string{true: "non static", false: "static"}[required.static], decl.name)}
			}
			if !e.compatibleMethodDeclaration(method, required) {
				if err := e.methodCompatibilityDiagnostic(method, required, method.line); err != nil {
					return err
				}
			}
		}
	}
	if classImplements(decl, "IteratorAggregate") {
		for _, method := range decl.methodOrder {
			if !strings.EqualFold(method.name, "getIterator") || method.returnType != "" {
				continue
			}
			suppressed := false
			for _, attribute := range method.attributes {
				if strings.EqualFold(attribute.name, "ReturnTypeWillChange") {
					suppressed = true
					break
				}
			}
			if suppressed {
				continue
			}
			prefix := ""
			if method.returnsByReference {
				prefix = "& "
			}
			message := fmt.Sprintf("Return type of %s%s::%s() should either be compatible with IteratorAggregate::getIterator(): Traversable, or the #[\\ReturnTypeWillChange] attribute should be used to temporarily suppress the notice", prefix, decl.name, method.name)
			if err := e.diagnostic("Deprecated", method.line, message); err != nil {
				return err
			}
		}
	}
	e.classes[key] = decl
	if decl.kind == "enum" {
		decl.constants = make(map[string]classConstantValue)
		if err := e.initEnumCases(decl); err != nil {
			delete(e.classes, key)
			return err
		}
	}
	if err := e.initClassMembers(decl); err != nil {
		delete(e.classes, key)
		return err
	}
	return nil
}

// Keep independent sibling contracts, but replace an ancestor contract when
// a descendant redeclares it. The redeclaration was checked at registration;
// checking the ancestor against that descendant again reverses variance.
func interfaceContracts(iface *classDecl) []*methodDecl {
	var result []*methodDecl
	seen := make(map[*classDecl]bool)
	var visit func(*classDecl)
	visit = func(current *classDecl) {
		if seen[current] {
			return
		}
		seen[current] = true
		for _, parent := range current.interfaces {
			visit(parent)
		}
		for _, method := range current.methodOrder {
			kept := result[:0]
			for _, previous := range result {
				if equalName(previous.name, method.name) && classImplements(current, previous.owner.name) {
					continue
				}
				kept = append(kept, previous)
			}
			result = append(kept, method)
		}
	}
	visit(iface)
	return result
}

func compatibleUntypedMethod(implementation, contract *methodDecl) bool {
	if requiredParameters(implementation) > requiredParameters(contract) {
		return false
	}
	for i, required := range contract.parameters {
		if i >= len(implementation.parameters) {
			break
		}
		if implementation.parameters[i].byRef != required.byRef {
			return false
		}
	}
	if len(contract.parameters) > 0 && contract.parameters[len(contract.parameters)-1].variadic &&
		(len(implementation.parameters) == 0 || !implementation.parameters[len(implementation.parameters)-1].variadic) {
		return false
	}
	if len(implementation.parameters) >= len(contract.parameters) {
		return true
	}
	return len(implementation.parameters) > 0 && implementation.parameters[len(implementation.parameters)-1].variadic
}

func methodSignature(method *methodDecl) string {
	parts := make([]string, len(method.parameters))
	for i, parameter := range method.parameters {
		parts[i] = "$" + parameter.name
		if parameter.variadic {
			parts[i] = "..." + parts[i]
		}
		if parameter.byRef {
			parts[i] = "&" + parts[i]
		}
		if parameter.typeName != "" {
			parts[i] = declaredTypeText(relativeTypeName(parameter.typeName, method.owner, nil), parameter.nullable) + " " + parts[i]
		}
		if initial, ok := parameter.initial.(literal); ok {
			switch initial.value.kind {
			case nullKind:
				parts[i] += " = null"
			case boolKind:
				if initial.value.truthy() {
					parts[i] += " = true"
				} else {
					parts[i] += " = false"
				}
			case stringKind:
				parts[i] += " = " + exportString(initial.value.s)
			case intKind, floatKind:
				parts[i] += " = " + initial.value.text()
			}
		}
	}
	signature := method.owner.name + "::" + method.name + "(" + strings.Join(parts, ", ") + ")"
	if method.returnType != "" {
		signature += ": " + declaredTypeText(relativeTypeName(method.returnType, method.owner, nil), method.returnNullable)
	}
	return signature
}

func (e *evaluator) initClassMembers(decl *classDecl) error {
	oldScope := e.scope
	e.scope = decl
	defer func() { e.scope = oldScope }()
	if decl.kind != "enum" || decl.constants == nil {
		decl.constants = make(map[string]classConstantValue)
	}
	decl.staticProperties = make(map[string]*staticPropertySlot)
	if decl.parent != nil {
		for name, constant := range decl.parent.constants {
			decl.constants[name] = constant
		}
		for name, slot := range decl.parent.staticProperties {
			decl.staticProperties[name] = slot
		}
	}
	for _, iface := range decl.interfaces {
		for name, constant := range iface.constants {
			if _, exists := decl.constants[name]; !exists {
				decl.constants[name] = constant
			}
		}
	}
	for _, constant := range decl.constantDecls {
		if old, exists := decl.constants[constant.name]; exists && old.decl.final {
			return &fatalDiagnostic{constant.line, decl.name + "::" + constant.name + " cannot override final constant " + old.decl.owner.name + "::" + constant.name}
		}
		if old, exists := decl.constants[constant.name]; exists && old.decl.visibility != "private" && old.decl.typeName != "" && !e.typeAcceptsScoped(old.decl.typeName, old.decl.nullable, old.decl.owner, constant.typeName, constant.nullable, constant.owner) {
			return &fatalDiagnostic{constant.line, fmt.Sprintf("Type of %s::%s must be compatible with %s::%s of type %s", decl.name, constant.name, old.decl.owner.name, old.decl.name, declaredTypeText(old.decl.typeName, old.decl.nullable))}
		}
		if old, exists := decl.constants[constant.name]; exists && constant.traitSource != "" && e.version != PHP85 && old.decl.owner != decl {
			compatible, err := e.compatibleClassConstant(constant, old.decl)
			if err != nil {
				return err
			}
			if !compatible {
				return &fatalDiagnostic{decl.line, fmt.Sprintf("%s and %s define the same constant (%s) in the composition of %s", old.decl.owner.name, constant.traitSource, constant.name, decl.name)}
			}
		}
		if !constantInitializer(constant.initial) {
			return &fatalDiagnostic{constant.line, "Constant expression contains invalid operations"}
		}
		v, err := e.eval(constant.initial)
		if err != nil {
			if missing, ok := err.(*runtimeError); ok && missing.class == "Error" && strings.HasPrefix(missing.message, "Class ") && strings.HasSuffix(missing.message, " not found") {
				// A class constant may name a class that is registered later,
				// including a class whose parent is this declaration.
				decl.constants[constant.name] = classConstantValue{decl: constant, pending: true}
				continue
			}
			return err
		}
		if constant.typeName != "" {
			_, _, objectFits := e.coerceType(v, relativeTypeName(constant.typeName, decl, nil), constant.nullable, true)
			if !literalFitsDeclaredType(v, constant.typeName, constant.nullable) && !(v.kind == objectKind && objectFits) {
				return &fatalDiagnostic{constant.line, fmt.Sprintf("Cannot use %s as value for class constant %s::%s of type %s", v.typeName(), decl.name, constant.name, declaredTypeText(constant.typeName, constant.nullable))}
			}
			if v.kind == intKind && constant.typeName == "float" {
				v = floating(float64(v.i))
			}
		}
		decl.constants[constant.name] = classConstantValue{value: v, decl: constant}
	}
	for _, property := range decl.properties {
		if property.owner != decl {
			continue
		}
		var v value
		if property.initial != nil {
			initial, err := e.eval(property.initial)
			if err != nil {
				return err
			}
			v = initial
			if property.typeName != "" {
				_, _, objectFits := e.coerceType(v, relativeTypeName(property.typeName, decl, nil), property.nullable, true)
				if !literalFitsDeclaredType(v, property.typeName, property.nullable) && !(v.kind == objectKind && objectFits) {
					return &fatalDiagnostic{property.line, fmt.Sprintf("Cannot use %s as default value for property %s::$%s of type %s", v.typeName(), decl.name, property.name, declaredTypeText(property.typeName, property.nullable))}
				}
				if v.kind == intKind && property.typeName == "float" {
					v = floating(float64(v.i))
				}
			}
		}
		if property.static {
			decl.staticProperties[property.name] = &staticPropertySlot{value: v, decl: property, initialized: property.initial != nil || property.typeName == ""}
		}
	}
	return nil
}

// Composition copies declarations into the receiving class. In particular,
// neither a method's declaring scope nor a property's private storage key may
// remain attached to the reusable trait definition.
func (e *evaluator) composeTraits(decl *classDecl) error {
	for _, use := range decl.traitUses {
		traits := make([]*classDecl, 0, len(use.names))
		for _, name := range use.names {
			trait := e.classes[strings.ToLower(name)]
			if trait == nil {
				if err := e.tryAutoload(name, decl.line); err != nil {
					return err
				}
				trait = e.classes[strings.ToLower(name)]
			}
			if trait == nil || trait.kind != "trait" {
				return &runtimeError{decl.line, "Error", "Trait \"" + name + "\" not found"}
			}
			traits = append(traits, trait)
		}
		type candidate struct {
			trait  *classDecl
			method *methodDecl
		}
		methods := make(map[string][]candidate)
		for _, trait := range traits {
			for name, method := range trait.methods {
				methods[name] = append(methods[name], candidate{trait, method})
			}
			for _, property := range trait.properties {
				copy := property
				copy.owner, copy.fromTrait, copy.traitSource = decl, true, trait.name
				compatibleDuplicate := false
				for _, existing := range decl.properties {
					if existing.name == property.name {
						compatible, err := e.compatibleProperty(copy, existing)
						if err != nil {
							return err
						}
						if !compatible {
							return &fatalDiagnostic{decl.line, fmt.Sprintf("%s and %s define the same property ($%s) in the composition of %s. However, the definition differs and is considered incompatible. Class was composed", decl.name, trait.name, property.name, decl.name)}
						}
						compatibleDuplicate = true
						break
					}
				}
				if !compatibleDuplicate {
					decl.properties = append(decl.properties, copy)
				}
			}
			for _, constant := range trait.constantDecls {
				copy := constant
				copy.owner = decl
				copy.traitSource = trait.name
				for _, existing := range decl.constantDecls {
					if existing.name == copy.name {
						compatible, err := e.compatibleClassConstant(copy, existing)
						if err != nil {
							return err
						}
						if !compatible {
							return &fatalDiagnostic{decl.line, fmt.Sprintf("%s and %s define the same constant (%s) in the composition of %s", decl.name, trait.name, copy.name, decl.name)}
						}
					}
				}
				decl.constantDecls = append(decl.constantDecls, copy)
			}
		}
		choose := func(adapt traitAdaptation) (*methodDecl, error) {
			list := methods[strings.ToLower(adapt.method)]
			if adapt.trait == "" && len(list) > 1 {
				return nil, syntax(adapt.line, "ambiguous trait method %s", adapt.method)
			}
			for _, option := range list {
				if adapt.trait == "" || equalName(option.trait.name, adapt.trait) {
					return option.method, nil
				}
			}
			return nil, syntax(adapt.line, "trait method %s not found", adapt.method)
		}
		excluded := make(map[string]map[string]bool)
		for _, adapt := range use.adaptations {
			if len(adapt.insteadOf) == 0 {
				continue
			}
			if _, err := choose(adapt); err != nil {
				return err
			}
			key := strings.ToLower(adapt.method)
			if excluded[key] == nil {
				excluded[key] = make(map[string]bool)
			}
			for _, name := range adapt.insteadOf {
				found := false
				for _, option := range methods[key] {
					if equalName(option.trait.name, name) && !equalName(name, adapt.trait) {
						found = true
						break
					}
				}
				if !found {
					return syntax(adapt.line, "trait method %s not found for insteadof", name)
				}
				excluded[key][strings.ToLower(name)] = true
			}
		}
		for name, list := range methods {
			// Only a method written in the receiving class resolves a trait
			// collision. A method copied by an earlier use statement does not.
			if existing := decl.methods[name]; existing != nil && existing.traitSource == "" {
				continue
			}
			var selected *methodDecl
			var selectedTrait *classDecl
			for _, option := range list {
				if excluded[name][strings.ToLower(option.trait.name)] {
					continue
				}
				if selected != nil {
					return &fatalDiagnostic{decl.line, fmt.Sprintf("Trait method %s::%s has not been applied as %s::%s, because of collision with %s::%s", option.trait.name, option.method.name, decl.name, option.method.name, selectedTrait.name, selected.name)}
				}
				selected = option.method
				selectedTrait = option.trait
			}
			if selected != nil {
				if existing := decl.methods[name]; existing != nil {
					origin := selected.traitSource
					if origin == "" {
						origin = selectedTrait.name
					}
					if equalName(existing.traitSource, origin) && equalName(existing.name, selected.name) {
						continue // The same underlying trait was imported through two paths.
					}
					return &fatalDiagnostic{decl.line, fmt.Sprintf("Trait method %s::%s has not been applied as %s::%s, because of collision with %s::%s", selectedTrait.name, selected.name, decl.name, selected.name, existing.traitSource, existing.name)}
				}
				copy := *selected
				copy.owner, copy.traitSource = decl, selectedTrait.name
				if selected.traitSource != "" {
					copy.traitSource = selected.traitSource
				}
				decl.methods[name] = &copy
				decl.methodOrder = append(decl.methodOrder, &copy)
			}
		}
		for _, adapt := range use.adaptations {
			if len(adapt.insteadOf) != 0 {
				continue
			}
			method, err := choose(adapt)
			if err != nil {
				return err
			}
			name := adapt.alias
			if name == "" {
				name = method.name
			}
			if existing := decl.methods[strings.ToLower(name)]; existing != nil && adapt.alias != "" {
				if existing.traitSource == "" {
					continue // A method written in the class takes precedence.
				}
				if !equalName(existing.traitSource, method.owner.name) || !equalName(existing.name, method.name) {
					return &fatalDiagnostic{decl.line, fmt.Sprintf("Trait method %s::%s has not been applied as %s::%s, because of collision with %s::%s", method.owner.name, method.name, decl.name, name, existing.traitSource, existing.name)}
				}
				continue
			}
			copy := *method
			copy.name, copy.owner = name, decl
			copy.traitSource = method.owner.name
			if adapt.visibility != "" {
				copy.visibility = adapt.visibility
			}
			if adapt.final {
				copy.final = true
			}
			decl.methods[strings.ToLower(name)] = &copy
			decl.methodOrder = append(decl.methodOrder, &copy)
		}
	}
	return nil
}

func (e *evaluator) compatibleProperty(a, b propertyDecl) (bool, error) {
	if a.visibility != b.visibility || a.static != b.static || (a.readonly || a.owner.readonly) != (b.readonly || b.owner.readonly) || !e.typeAcceptsScoped(a.typeName, a.nullable, a.owner, b.typeName, b.nullable, b.owner) || !e.typeAcceptsScoped(b.typeName, b.nullable, b.owner, a.typeName, a.nullable, a.owner) {
		return false, nil
	}
	var left, right value
	var err error
	if a.initial != nil {
		left, err = e.eval(a.initial)
		if err != nil {
			return false, err
		}
	}
	if b.initial != nil {
		right, err = e.eval(b.initial)
		if err != nil {
			return false, err
		}
	}
	return reflect.DeepEqual(left, right), nil
}

func (e *evaluator) compatibleClassConstant(a, b classConstantDecl) (bool, error) {
	if a.visibility != b.visibility || a.final != b.final || !e.typeAcceptsScoped(a.typeName, a.nullable, a.owner, b.typeName, b.nullable, b.owner) || !e.typeAcceptsScoped(b.typeName, b.nullable, b.owner, a.typeName, a.nullable, a.owner) {
		return false, nil
	}
	left, err := e.eval(a.initial)
	if err != nil {
		return false, err
	}
	right, err := e.eval(b.initial)
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(left, right), nil
}

func visibilityRank(visibility string) int {
	switch visibility {
	case "public":
		return 3
	case "protected":
		return 2
	default:
		return 1
	}
}

func visibilitySuffix(visibility string) string {
	if visibility == "public" {
		return ""
	}
	return " or weaker"
}

func (e *evaluator) resolveClassName(name string, line int) (*classDecl, error) {
	e.initClasses()
	switch strings.ToLower(name) {
	case "self":
		if e.scope == nil {
			return nil, &runtimeError{line, "Error", "Cannot access self:: when no class scope is active"}
		}
		return e.scope, nil
	case "parent":
		if e.scope == nil {
			return nil, &runtimeError{line, "Error", "Cannot access parent:: when no class scope is active"}
		}
		if e.scope.parent == nil {
			return nil, &runtimeError{line, "Error", "Cannot access parent:: when current class scope has no parent"}
		}
		return e.scope.parent, nil
	case "static":
		if e.lateClass == nil {
			return nil, &runtimeError{line, "Error", "Cannot access static:: when no class scope is active"}
		}
		return e.lateClass, nil
	}
	class := e.classes[strings.ToLower(name)]
	if class == nil {
		if err := e.tryAutoload(name, line); err != nil {
			return nil, err
		}
		class = e.classes[strings.ToLower(name)]
	}
	if class == nil {
		return nil, &runtimeError{line, "Error", fmt.Sprintf("Class %q not found", name)}
	}
	return class, nil
}

func (e *evaluator) readClassConstant(node classConstantExpr) (value, error) {
	class, err := e.resolveClassName(node.className, node.line)
	if err != nil {
		return value{}, err
	}
	if class.kind == "enum" && class.enumDuplicate != "" {
		return value{}, &runtimeError{node.line, "Error", class.enumDuplicate}
	}
	constant, ok := class.constants[node.name]
	if !ok {
		return value{}, &runtimeError{node.line, "Error", "Undefined constant " + class.name + "::" + node.name}
	}
	owner := constant.decl.owner
	if constant.decl.visibility == "private" && e.scope != owner || constant.decl.visibility == "protected" && !sameHierarchy(e.scope, owner) {
		return value{}, &runtimeError{node.line, "Error", "Cannot access " + constant.decl.visibility + " constant " + class.name + "::" + node.name}
	}
	if constant.pending {
		if owner != nil {
			if current, exists := owner.constants[node.name]; exists {
				constant = current
			}
		}
		if constant.pending {
			if constant.resolving {
				return value{}, &runtimeError{node.line, "Error", "Cannot declare self-referencing constant " + owner.name + "::" + node.name}
			}
			constant.resolving = true
			owner.constants[node.name] = constant
			oldScope := e.scope
			e.scope = owner
			v, err := e.eval(constant.decl.initial)
			e.scope = oldScope
			if err != nil {
				constant.resolving = false
				owner.constants[node.name] = constant
				return value{}, err
			}
			if constant.decl.typeName != "" {
				_, _, objectFits := e.coerceType(v, relativeTypeName(constant.decl.typeName, owner, nil), constant.decl.nullable, true)
				if !literalFitsDeclaredType(v, constant.decl.typeName, constant.decl.nullable) && !(v.kind == objectKind && objectFits) {
					constant.resolving = false
					owner.constants[node.name] = constant
					return value{}, &runtimeError{node.line, "TypeError", fmt.Sprintf("Cannot use %s as value for class constant %s::%s of type %s", v.typeName(), owner.name, node.name, declaredTypeText(constant.decl.typeName, constant.decl.nullable))}
				}
				if v.kind == intKind && constant.decl.typeName == "float" {
					v = floating(float64(v.i))
				}
			}
			constant.value, constant.pending, constant.resolving = v, false, false
			owner.constants[node.name] = constant
		}
	}
	return constant.value, nil
}

func (e *evaluator) resolveStaticClass(name string, dynamic expr, line int) (*classDecl, error) {
	if dynamic == nil {
		return e.resolveClassName(name, line)
	}
	v, err := e.eval(dynamic)
	if err != nil {
		return nil, err
	}
	if v.kind == objectKind && v.o != nil {
		if v.o.class != nil {
			return v.o.class, nil
		}
		return e.resolveClassName("stdClass", line)
	}
	if v.kind != stringKind {
		return nil, &runtimeError{line, "Error", "Class name must be a valid object or a string"}
	}
	// A name supplied at runtime is looked up literally. In particular, a
	// string containing "self", "parent", or "static" is not a scope keyword.
	name = strings.TrimPrefix(v.s, "\\")
	e.initClasses()
	class := e.classes[strings.ToLower(name)]
	if class == nil {
		if err := e.tryAutoload(name, line); err != nil {
			return nil, err
		}
		class = e.classes[strings.ToLower(name)]
	}
	if class == nil {
		return nil, &runtimeError{line, "Error", fmt.Sprintf("Class %q not found", name)}
	}
	return class, nil
}

func (e *evaluator) staticMemberName(name string, dynamic expr, line int) (string, error) {
	if dynamic == nil {
		return name, nil
	}
	v, err := e.eval(dynamic)
	if err != nil {
		return "", err
	}
	if v.kind != stringKind {
		return "", &runtimeError{line, "Error", "Method name must be a string"}
	}
	return v.s, nil
}

func (e *evaluator) prepareStaticProperty(node staticPropertyExpr) (staticPropertyExpr, error) {
	if node.classExpr == nil && node.nameExpr == nil {
		return node, nil
	}
	class, err := e.resolveStaticClass(node.className, node.classExpr, node.line)
	if err != nil {
		return node, err
	}
	name, err := e.staticMemberName(node.name, node.nameExpr, node.line)
	if err != nil {
		return node, err
	}
	node.className, node.classExpr, node.name, node.nameExpr = class.name, nil, name, nil
	return node, nil
}

func (e *evaluator) staticPropertyLocation(node staticPropertyExpr) (location, error) {
	class, err := e.resolveStaticClass(node.className, node.classExpr, node.line)
	if err != nil {
		return location{}, err
	}
	name, err := e.staticMemberName(node.name, node.nameExpr, node.line)
	if err != nil {
		return location{}, err
	}
	slot := class.staticProperties[name]
	if slot == nil {
		return location{}, &runtimeError{node.line, "Error", "Access to undeclared static property " + class.name + "::$" + name}
	}
	owner := slot.decl.owner
	if slot.decl.visibility == "private" && e.scope != owner || slot.decl.visibility == "protected" && !sameHierarchy(e.scope, owner) {
		return location{}, &runtimeError{node.line, "Error", "Cannot access " + slot.decl.visibility + " property " + class.name + "::$" + name}
	}
	return location{func() value { return dereference(slot.value) }, func(v value) error {
		if slot.value.kind == refKind {
			if err := e.writeReference(slot.value.ref, v, node.line); err != nil {
				return err
			}
			slot.initialized = true
			return nil
		}
		if slot.decl.typeName != "" {
			converted, selected, ok := e.coerceType(v, relativeTypeName(slot.decl.typeName, slot.decl.owner, nil), slot.decl.nullable, e.strictTypes)
			if !ok {
				return &runtimeError{node.line, "TypeError", fmt.Sprintf("Cannot assign %s to property %s::$%s of type %s", v.typeName(), slot.decl.owner.name, name, declaredTypeText(slot.decl.typeName, slot.decl.nullable))}
			}
			if err := e.warnLossyIntTypeCoercion(v, selected, e.strictTypes, node.line); err != nil {
				return err
			}
			v = converted
		}
		slot.value = v
		slot.initialized = true
		return nil
	}, slot.decl.typeName != ""}, nil
}

func (e *evaluator) callStaticMethod(node staticCallExpr) (value, error) {
	class, err := e.resolveStaticClass(node.className, node.classExpr, node.line)
	if err != nil {
		return value{}, err
	}
	name, err := e.staticMemberName(node.name, node.nameExpr, node.line)
	if err != nil {
		return value{}, err
	}
	node.name = name
	if strings.EqualFold(class.name, "Closure") && strings.EqualFold(name, "fromCallable") {
		return e.closureFromCallable(node.args, node.line)
	}
	if strings.EqualFold(class.name, "Closure") && strings.EqualFold(name, "bind") {
		return e.closureBind(node.args, node.line)
	}
	if result, handled, err := e.callEnumStatic(class, name, node.args, node.line); handled {
		return result, err
	}
	method := e.lookupMethod(class, name)
	if method == nil {
		if !strings.EqualFold(name, "__callStatic") {
			if fallback := e.lookupMethod(class, "__callStatic"); fallback != nil {
				args, err := e.evalArguments(node.args)
				if err != nil {
					return value{}, err
				}
				return e.invokeWithClass(value{}, class, fallback, []callArgument{{value: value{kind: stringKind, s: name}}, {value: packMagicArguments(args)}}, node.line)
			}
		}
		if isThrowableClass(class) && e.scope != nil {
			if object, ok := e.readVariable("this"); ok && object.kind == objectKind && object.o.exception != nil {
				for owner := object.o.class; owner != nil; owner = owner.parent {
					if owner == class {
						if result, handled, err := e.callExceptionAccessor(object, node.name, node.args, node.line); handled {
							return result, err
						}
						break
					}
				}
			}
		}
		if strings.EqualFold(node.name, "__construct") && isThrowableClass(class) && e.scope != nil {
			if object, ok := e.readVariable("this"); ok && object.kind == objectKind && object.o.exception != nil {
				args, err := e.evalArguments(node.args)
				if err != nil {
					return value{}, err
				}
				return value{}, e.initializeExceptionArguments(object, args, node.line)
			}
		}
		return value{}, &runtimeError{node.line, "Error", "Call to undefined method " + class.name + "::" + node.name + "()"}
	}
	if err := e.methodAccess(class, method, node.line); err != nil {
		return value{}, err
	}
	if !method.static {
		if this, ok := e.readVariable("this"); ok && this.kind == objectKind && sameHierarchy(this.o.class, method.owner) {
			args, err := e.evalArguments(node.args)
			if err != nil {
				return value{}, err
			}
			return e.invoke(this, method, args, node.line)
		}
		return value{}, &runtimeError{node.line, "Error", "Non-static method " + method.owner.name + "::" + method.name + "() cannot be called statically"}
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	called := class
	// self::, parent::, and static:: are forwarding calls: method lookup
	// uses the named scope, but static:: inside the callee still sees the
	// class from the last non-forwarding call.
	switch strings.ToLower(node.className) {
	case "self", "parent", "static":
		if e.lateClass != nil {
			called = e.lateClass
		}
	}
	return e.invokeWithClass(value{}, called, method, args, node.line)
}

func (e *evaluator) newInstance(node newObject) (value, error) {
	className := node.name
	var class *classDecl
	var err error
	if node.anonymous != nil {
		class = node.anonymous
		e.initClasses()
		if registered := e.classes[strings.ToLower(class.name)]; registered != nil {
			if registered == class {
				// Repeated evaluation of one compiled declaration reuses its class.
				class = registered
			} else {
				// Re-including recompiles the file. The new declaration has the
				// same source location but a distinct class identity.
				base := class.name
				for suffix := 2; ; suffix++ {
					candidate := fmt.Sprintf("%s#%d", base, suffix)
					if e.classes[strings.ToLower(candidate)] == nil {
						class.name = candidate
						break
					}
				}
				if err := e.declareClass(class); err != nil {
					return value{}, err
				}
			}
		} else {
			if err := e.declareClass(class); err != nil {
				return value{}, err
			}
		}
	} else if node.classExpr != nil {
		v, evalErr := e.eval(node.classExpr)
		if evalErr != nil {
			return value{}, evalErr
		}
		if v.kind != stringKind {
			return value{}, &runtimeError{node.line, "Error", "Class name must be a valid object or a string"}
		}
		className = strings.TrimPrefix(v.s, "\\")
		class, err = e.resolveClassName(className, node.line)
	} else {
		class, err = e.resolveClassName(className, node.line)
	}
	if err != nil {
		return value{}, err
	}
	if class.kind == "interface" {
		return value{}, &runtimeError{node.line, "Error", "Cannot instantiate interface " + class.name}
	}
	if class.kind == "trait" {
		return value{}, &runtimeError{node.line, "Error", "Cannot instantiate trait " + class.name}
	}
	if class.kind == "enum" {
		return value{}, &runtimeError{node.line, "Error", "Cannot instantiate enum " + class.name}
	}
	if strings.EqualFold(class.name, "Closure") {
		return value{}, &runtimeError{node.line, "Error", "Instantiation of class Closure is not allowed"}
	}
	if class.abstract {
		return value{}, &runtimeError{node.line, "Error", "Cannot instantiate abstract class " + class.name}
	}
	object := value{kind: objectKind, o: &phpObject{class: class, properties: &phpArray{}}}
	if strings.EqualFold(class.name, "DOMDocument") {
		object.o.dom = &domNode{name: "#document", object: object.o}
		return object, nil
	}
	if strings.EqualFold(class.name, "ReflectionMethod") {
		args, argErr := e.evalArguments(node.args)
		if argErr != nil {
			return value{}, argErr
		}
		bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "objectOrMethod"}, {name: "method"}}, node.line)
		if bindErr != nil {
			return value{}, bindErr
		}
		if !present[0] || !present[1] || len(args) != 2 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionMethod::__construct() expects exactly 2 arguments"}
		}
		var target *classDecl
		if bound[0].kind == objectKind {
			target = bound[0].o.class
		} else if bound[0].kind == stringKind {
			target, argErr = e.resolveClassName(bound[0].s, node.line)
			if argErr != nil {
				return value{}, argErr
			}
		} else {
			return value{}, builtinTypeError(node.line, "ReflectionMethod::__construct", 1, "objectOrMethod", "object|string", bound[0])
		}
		methodName, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "ReflectionMethod::__construct", 2, "method", "string", bound[1])
		}
		object.o.reflectedMethod = e.lookupMethod(target, methodName.s)
		if object.o.reflectedMethod == nil {
			return value{}, &runtimeError{node.line, "ReflectionException", "Method " + target.name + "::" + methodName.s + "() does not exist"}
		}
		return object, nil
	}
	if strings.EqualFold(class.name, "ReflectionFunction") {
		args, argErr := e.evalArguments(node.args)
		if argErr != nil {
			return value{}, argErr
		}
		bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "function"}}, node.line)
		if bindErr != nil {
			return value{}, bindErr
		}
		if len(args) != 1 || !present[0] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionFunction::__construct() expects exactly 1 argument"}
		}
		v := bound[0]
		if v.kind == objectKind && v.o != nil && v.o.closure != nil {
			// A first-class callable can wrap another Closure repeatedly (for
			// example $callback(...), followed by another capture). Reflection
			// describes the underlying callable rather than the wrapper depth.
			seen := map[*phpClosure]bool{}
			for current := v.o.closure; current != nil && !seen[current]; {
				seen[current] = true
				object.o.reflectedClosure = current
				if current.method != nil {
					object.o.reflectedMethod = current.method
					break
				}
				if current.target == nil {
					break
				}
				if current.target.method != nil {
					object.o.reflectedMethod = current.target.method
					break
				}
				current = current.target.closure
			}
		} else if v.kind == stringKind {
			object.o.reflectedMethod = e.functions[symbolKey(strings.TrimPrefix(v.s, "\\"))]
		}
		if object.o.reflectedMethod == nil {
			return value{}, &runtimeError{node.line, "ReflectionException", "Function does not exist"}
		}
		return object, nil
	}
	if strings.EqualFold(class.name, "ReflectionProperty") {
		args, argErr := e.evalArguments(node.args)
		if argErr != nil {
			return value{}, argErr
		}
		bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "class"}, {name: "property"}}, node.line)
		if bindErr != nil {
			return value{}, bindErr
		}
		if len(args) != 2 || !present[0] || !present[1] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionProperty::__construct() expects exactly 2 arguments"}
		}
		var target *classDecl
		if bound[0].kind == objectKind {
			target = bound[0].o.class
		} else if bound[0].kind == stringKind {
			target, argErr = e.resolveClassName(bound[0].s, node.line)
			if argErr != nil {
				return value{}, argErr
			}
		}
		if target == nil {
			return value{}, builtinTypeError(node.line, "ReflectionProperty::__construct", 1, "class", "object|string", bound[0])
		}
		if bound[1].kind != stringKind {
			return value{}, builtinTypeError(node.line, "ReflectionProperty::__construct", 2, "property", "string", bound[1])
		}
		for cls := target; cls != nil; cls = cls.parent {
			for i := range cls.properties {
				if cls.properties[i].name == bound[1].s {
					object.o.reflectedProperty = &cls.properties[i]
					return object, nil
				}
			}
		}
		return value{}, &runtimeError{node.line, "ReflectionException", "Property " + target.name + "::$" + bound[1].s + " does not exist"}
	}
	if strings.EqualFold(class.name, "ReflectionClass") || strings.EqualFold(class.name, "ReflectionEnum") {
		args, argErr := e.evalArguments(node.args)
		if argErr != nil {
			return value{}, argErr
		}
		bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "objectOrClass"}}, node.line)
		if bindErr != nil {
			return value{}, bindErr
		}
		if len(args) != 1 || !present[0] {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::__construct() expects exactly 1 argument"}
		}
		if bound[0].kind == objectKind {
			object.o.reflected = bound[0].o.class
		} else if bound[0].kind == stringKind {
			object.o.reflected, argErr = e.resolveClassName(bound[0].s, node.line)
			if argErr != nil {
				return value{}, argErr
			}
		} else {
			return value{}, builtinTypeError(node.line, "ReflectionClass::__construct", 1, "objectOrClass", "object|string", bound[0])
		}
		if strings.EqualFold(class.name, "ReflectionEnum") && object.o.reflected.kind != "enum" {
			return value{}, &runtimeError{node.line, "ReflectionException", "Class " + object.o.reflected.name + " is not an enum"}
		}
		return object, nil
	}
	if strings.EqualFold(class.name, "ArrayIterator") {
		args, argErr := e.evalArguments(node.args)
		if argErr != nil {
			return value{}, argErr
		}
		bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "array"}, {name: "flags"}}, node.line)
		if bindErr != nil {
			return value{}, bindErr
		}
		if len(args) > 2 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "ArrayIterator::__construct() expects at most 2 arguments"}
		}
		items := arrayValue()
		if present[0] {
			if bound[0].kind != arrayKind {
				return value{}, builtinTypeError(node.line, "ArrayIterator::__construct", 1, "array", "array|object", bound[0])
			}
			items = bound[0]
		}
		object.o.iteratorArray = items.a
		return object, nil
	}
	if isThrowableClass(class) {
		object.o.exception = &exceptionData{file: e.filename, line: node.line, frames: e.snapshotCallFrames(), captured: true}
	}
	for _, property := range class.properties {
		if property.static {
			continue
		}
		var v value
		var err error
		if property.initial != nil {
			declaringScope := e.scope
			e.scope = property.owner
			v, err = e.eval(property.initial)
			e.scope = declaringScope
		}
		if err != nil {
			return value{}, err
		}
		if property.typeName != "" && property.initial == nil {
			continue
		}
		if property.typeName == "float" && v.kind == intKind {
			v = floating(float64(v.i))
		}
		object.o.properties = object.o.properties.with(propertyStorageKey(property), v)
	}
	if method := class.methods["__construct"]; method != nil {
		if err := e.methodAccess(class, method, node.line); err != nil {
			// Instantiation diagnostics omit "method" (unlike an explicit call).
			if failure, ok := err.(*runtimeError); ok {
				failure.message = strings.Replace(failure.message, " method ", " ", 1)
			}
			return value{}, err
		}
		args, err := e.evalArguments(node.args)
		if err != nil {
			return value{}, err
		}
		if object.o.exception != nil {
			object.o.exception.frames = e.snapshotCallFrames()
		}
		if _, err := e.invoke(object, method, args, node.line); err != nil {
			return value{}, err
		}
	} else if object.o.exception != nil {
		if err := e.initializeException(object, node.args, node.line); err != nil {
			return value{}, err
		}
	} else {
		// PHP evaluates and ignores positional arguments if no constructor exists.
		if _, err := e.evalArguments(node.args); err != nil {
			return value{}, err
		}
	}
	return object, nil
}

type callArgument struct {
	name   string // Empty for a positional argument.
	named  bool
	value  value
	ref    *referenceCell  // Raw reference returned by an argument call, if any.
	place  *referencePlace // Writable argument binding captured during evaluation, when requested.
	line   int
	source expr // Original expression, for reference parameters.
}

// Return the original expression bound to a parameter after named arguments
// have been resolved. An unpacked value has no independent writable source.
func argumentSource(args []callArgument, parameter string, position int) expr {
	if arg := argumentForParameter(args, parameter, position); arg != nil {
		return arg.source
	}
	return nil
}

func argumentForParameter(args []callArgument, parameter string, position int) *callArgument {
	positional := 0
	for i := range args {
		arg := &args[i]
		if arg.named {
			if arg.name == parameter {
				return arg
			}
			continue
		}
		if positional == position {
			return arg
		}
		positional++
	}
	return nil
}

func argumentValues(args []callArgument) []value {
	values := make([]value, len(args))
	for i, arg := range args {
		values[i] = arg.value
	}
	return values
}

func (e *evaluator) evalArguments(args []expr) ([]callArgument, error) {
	return e.evalArgumentsWithCapture(args, nil)
}

func (e *evaluator) evalArgumentsCapturingArray(args []expr) ([]callArgument, error) {
	return e.evalArgumentsWithCapture(args, func(operand expr, name string, named bool, position int) (value, *referencePlace, bool, error) {
		if (named && name == "array" || !named && position == 0) && referenceTargetSyntax(operand) && !hasNullsafe(operand) {
			place, err := e.referencePlace(operand, true)
			if err != nil {
				return value{}, nil, true, err
			}
			return place.read(), &place, true, nil
		}
		return value{}, nil, false, nil
	})
}

// capture is used by builtins that need a writable argument. It resolves an
// lvalue while the argument is evaluated, preserving receiver/key side effects.
func (e *evaluator) evalArgumentsWithCapture(args []expr, capture func(expr, string, bool, int) (value, *referencePlace, bool, error)) ([]callArgument, error) {
	values := make([]callArgument, 0, len(args))
	seenNamed := false
	positional := 0
	for _, arg := range args {
		if unpack, ok := arg.(unpackArgument); ok {
			v, err := e.eval(unpack.operand)
			if err != nil {
				return nil, err
			}
			if v.kind == objectKind && v.o != nil && classImplements(v.o.class, "Traversable") {
				v, err = e.iteratorToArrayValue(v, true, unpack.line)
				if err != nil {
					return nil, err
				}
			}
			if v.kind != arrayKind {
				return nil, &runtimeError{unpack.line, "Error", "Only arrays and Traversables can be unpacked"}
			}
			for _, entry := range v.a.entries {
				if entry.key.isString {
					seenNamed = true
					values = append(values, callArgument{name: entry.key.text, named: true, value: dereference(entry.value), line: unpack.line})
					continue
				}
				if seenNamed {
					return nil, &runtimeError{unpack.line, "Error", "Cannot use positional argument after named argument"}
				}
				values = append(values, callArgument{value: dereference(entry.value), line: unpack.line})
				positional++
			}
			continue
		}
		name, isNamed := "", false
		operand := arg
		if named, ok := arg.(namedArgument); ok {
			name, operand = named.name, named.operand
			isNamed = true
			seenNamed = true
		}
		var v value
		var place *referencePlace
		var err error
		captured := false
		if capture != nil {
			v, place, captured, err = capture(operand, name, isNamed, positional)
		}
		if !captured && err == nil && !hasNullsafe(operand) {
			switch call := operand.(type) {
			case functionCall:
				v, err = e.callFunction(call)
			case methodCall:
				v, err = e.callMethod(call)
			case staticCallExpr:
				v, err = e.callStaticMethod(call)
			case dynamicCall:
				v, err = e.callDynamic(call)
			default:
				v, err = e.eval(operand)
			}
		} else if !captured && err == nil {
			v, err = e.eval(operand)
		}
		if e.pendingDiagnostic != nil {
			err = e.pendingDiagnostic
			e.pendingDiagnostic = nil
		}
		if err != nil {
			return nil, err
		}
		argument := callArgument{name: name, named: isNamed, value: dereference(v), place: place, line: expressionLine(arg), source: operand}
		if v.kind == refKind {
			argument.ref = v.ref
		}
		values = append(values, argument)
		if !isNamed {
			positional++
		}
	}
	return values, nil
}

// Evaluation order is represented by args; parameter order is applied only
// after all argument expressions have been evaluated.
func bindCallArguments(args []callArgument, parameters []parameter, line int) ([]value, []bool, value, error) {
	bound := make([]value, len(parameters))
	present := make([]bool, len(parameters))
	variadic := -1
	if len(parameters) > 0 && parameters[len(parameters)-1].variadic {
		variadic = len(parameters) - 1
	}
	extra := arrayValue()
	nextPos := 0
	seenNamed := false
	seenNames := make(map[string]bool)
	for _, arg := range args {
		if arg.named {
			seenNamed = true
			if seenNames[arg.name] {
				return nil, nil, value{}, &runtimeError{line, "Error", "Named parameter $" + arg.name + " overwrites previous argument"}
			}
			seenNames[arg.name] = true
			index := -1
			for i, param := range parameters {
				if param.name == arg.name && i != variadic {
					index = i
					break
				}
			}
			if index >= 0 {
				if present[index] {
					return nil, nil, value{}, &runtimeError{line, "Error", "Named parameter $" + arg.name + " overwrites previous argument"}
				}
				bound[index], present[index] = arg.value, true
			} else if variadic >= 0 {
				extra.a = extra.a.with(arrayKey{isString: true, text: arg.name}, arg.value)
			} else {
				return nil, nil, value{}, &runtimeError{line, "Error", "Unknown named parameter $" + arg.name}
			}
			continue
		}
		if seenNamed {
			return nil, nil, value{}, &runtimeError{arg.line, "Error", "Cannot use positional argument after named argument"}
		}
		if nextPos < len(parameters) && nextPos != variadic {
			bound[nextPos], present[nextPos] = arg.value, true
		} else if variadic >= 0 {
			extra.a = extra.a.with(arrayKey{index: int64(len(extra.a.entries))}, arg.value)
		}
		nextPos++
	}
	return bound, present, extra, nil
}

// Reference parameters bind the caller's variable slot before entering the
// callee. Keep this separate from value binding so named arguments still use
// the shared binder's parameter order.
func (e *evaluator) bindVariableReferenceArguments(args []callArgument, parameters []parameter, bound []value) ([]*referenceCell, error) {
	cells := make([]*referenceCell, len(parameters))
	nextPos := 0
	for _, arg := range args {
		index := -1
		if arg.named {
			for i, parameter := range parameters {
				if parameter.name == arg.name && !parameter.variadic {
					index = i
					break
				}
			}
		} else {
			index = nextPos
			nextPos++
		}
		if index < 0 || index >= len(parameters) || !parameters[index].byRef {
			continue
		}
		if parameters[index].variadic {
			return nil, &runtimeError{arg.line, "Error", "gophp: variadic reference parameters are not supported"}
		}
		if arg.ref != nil && referenceSourceSyntax(arg.source) {
			cells[index] = arg.ref
			bound[index] = arg.ref.value
			continue
		}
		var place referencePlace
		if arg.place != nil {
			place = *arg.place
		} else {
			if !referenceTargetSyntax(arg.source) {
				return nil, &runtimeError{arg.line, "Error", "Only variables can be passed by reference"}
			}
			var err error
			place, err = e.referencePlace(arg.source, true)
			if err != nil {
				return nil, err
			}
		}
		stored := place.read()
		if stored.kind != refKind {
			stored = value{kind: refKind, ref: &referenceCell{value: dereference(stored)}}
			if err := place.replace(stored); err != nil {
				return nil, err
			}
		}
		cell := stored.ref
		if place.source != nil {
			e.pruneReferenceSources(cell)
			seen := false
			for _, source := range cell.sources {
				if source.sameBinding(*place.source) {
					seen = true
					break
				}
			}
			if !seen {
				cell.sources = append(cell.sources, *place.source)
			}
		}
		cells[index] = cell
		bound[index] = cell.value
	}
	return cells, nil
}

func (e *evaluator) callMethod(node methodCall) (value, error) {
	object, err := e.eval(node.receiver)
	if err != nil {
		return value{}, err
	}
	if node.nameExpr != nil {
		node.name, err = e.instancePropertyName(node.name, node.nameExpr, node.line)
		if err != nil {
			return value{}, err
		}
	}
	if object.kind != objectKind {
		return value{}, &runtimeError{node.line, "Error", "Call to a member function " + node.name + "() on " + object.typeName()}
	}
	if object.o != nil && strings.EqualFold(object.o.className(), "DOMNodeList") && strings.EqualFold(node.name, "item") {
		args, argErr := e.evalArguments(node.args)
		if argErr != nil {
			return value{}, argErr
		}
		if len(args) != 1 || args[0].value.kind != intKind {
			return value{}, &runtimeError{node.line, "TypeError", "DOMNodeList::item(): Argument #1 ($index) must be of type int"}
		}
		if v, ok := object.o.iteratorArray.get(arrayKey{index: args[0].value.i}); ok {
			return v, nil
		}
		return value{kind: nullKind}, nil
	}
	if object.o != nil && strings.EqualFold(object.o.className(), "DOMNodeList") {
		for _, iteratorName := range []string{"rewind", "valid", "key", "current", "next"} {
			if strings.EqualFold(node.name, iteratorName) {
				return e.iteratorMethod(object, iteratorName, node.line)
			}
		}
	}
	if object.o != nil && object.o.dom != nil {
		if result, handled, err := e.callDOMMethod(object, node); handled {
			return result, err
		}
	}
	if object.o != nil && object.o.closure != nil && strings.EqualFold(node.name, "bindTo") {
		arguments := make([]expr, 0, len(node.args)+1)
		arguments = append(arguments, literal{value: object, line: node.line})
		arguments = append(arguments, node.args...)
		return e.closureBind(arguments, node.line)
	}
	if object.o.reflected != nil {
		args, argErr := e.evalArguments(node.args)
		if argErr != nil {
			return value{}, argErr
		}
		switch strings.ToLower(node.name) {
		case "iscloneable":
			if len(args) != 0 {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::isCloneable() expects exactly 0 arguments"}
			}
			return boolean(object.o.reflected.kind != "enum"), nil
		case "gettraits":
			if len(args) != 0 {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::getTraits() expects exactly 0 arguments"}
			}
			e.initClasses()
			out := arrayValue()
			for _, use := range object.o.reflected.traitUses {
				for _, name := range use.names {
					trait := e.classes[strings.ToLower(strings.TrimPrefix(name, "\\"))]
					if trait != nil {
						out.a = out.a.with(arrayKey{isString: true, text: trait.name}, value{kind: objectKind, o: &phpObject{class: e.classes["reflectionclass"], properties: &phpArray{}, reflected: trait}})
					}
				}
			}
			return out, nil
		case "getinterfacenames":
			if len(args) != 0 {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::getInterfaceNames() expects exactly 0 arguments"}
			}
			out := arrayValue()
			seen := map[string]bool{}
			var add func(*classDecl)
			add = func(c *classDecl) {
				if c == nil {
					return
				}
				for _, iface := range c.interfaces {
					if seen[strings.ToLower(iface.name)] {
						continue
					}
					seen[strings.ToLower(iface.name)] = true
					out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, value{kind: stringKind, s: iface.name})
					add(iface)
				}
				add(c.parent)
			}
			add(object.o.reflected)
			return out, nil
		case "getinterfaces":
			e.initClasses()
			out := arrayValue()
			seen := map[string]bool{}
			var add func(*classDecl)
			add = func(c *classDecl) {
				if c == nil {
					return
				}
				for _, iface := range c.interfaces {
					if seen[strings.ToLower(iface.name)] {
						continue
					}
					seen[strings.ToLower(iface.name)] = true
					out.a = out.a.with(arrayKey{isString: true, text: iface.name}, value{kind: objectKind, o: &phpObject{class: e.classes["reflectionclass"], properties: &phpArray{}, reflected: iface}})
					add(iface)
				}
				add(c.parent)
			}
			add(object.o.reflected)
			return out, nil
		case "getmethods":
			if len(args) > 1 {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::getMethods() expects at most 1 argument"}
			}
			e.initClasses()
			out := arrayValue()
			seen := map[string]bool{}
			for c := object.o.reflected; c != nil; c = c.parent {
				for _, method := range c.methodOrder {
					key := strings.ToLower(method.name)
					if seen[key] {
						continue
					}
					seen[key] = true
					if c != object.o.reflected && method.visibility == "private" {
						continue
					}
					ref := value{kind: objectKind, o: &phpObject{class: e.classes["reflectionmethod"], properties: &phpArray{}, reflectedMethod: method}}
					out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, ref)
				}
			}
			return out, nil
		case "newinstanceargs":
			if len(args) != 1 {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::newInstanceArgs() expects exactly 1 argument"}
			}
			if args[0].value.kind != arrayKind || args[0].value.a == nil {
				return value{}, builtinTypeError(node.line, "ReflectionClass::newInstanceArgs", 1, "args", "array", args[0].value)
			}
			arguments := make([]expr, 0, len(args[0].value.a.entries))
			for _, entry := range args[0].value.a.entries {
				item := expr(literal{value: dereference(entry.value), line: node.line})
				if entry.key.isString {
					item = namedArgument{name: entry.key.text, operand: item, line: node.line}
				}
				arguments = append(arguments, item)
			}
			return e.newInstance(newObject{name: object.o.reflected.name, args: arguments, line: node.line})
		case "newinstancewithoutconstructor":
			if len(args) != 0 {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::newInstanceWithoutConstructor() expects exactly 0 arguments"}
			}
			class := object.o.reflected
			if class.kind != "class" || class.abstract {
				return value{}, &runtimeError{node.line, "ReflectionException", "Class " + class.name + " is not instantiable"}
			}
			instance := value{kind: objectKind, o: &phpObject{class: class, properties: &phpArray{}}}
			for _, property := range class.properties {
				if property.static || property.initial == nil {
					continue
				}
				oldScope := e.scope
				e.scope = property.owner
				initial, evalErr := e.eval(property.initial)
				e.scope = oldScope
				if evalErr != nil {
					return value{}, evalErr
				}
				instance.o.properties = instance.o.properties.with(propertyStorageKey(property), initial)
			}
			return instance, nil
		case "getproperties":
			if len(args) > 1 {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::getProperties() expects at most 1 argument"}
			}
			e.initClasses()
			out := arrayValue()
			for c := object.o.reflected; c != nil; c = c.parent {
				for i := range c.properties {
					property := &c.properties[i]
					if c != object.o.reflected && property.visibility == "private" {
						continue
					}
					ref := value{kind: objectKind, o: &phpObject{class: e.classes["reflectionproperty"], properties: &phpArray{}, reflectedProperty: property}}
					out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, ref)
				}
			}
			return out, nil
		case "hasproperty", "getproperty":
			bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "name"}}, node.line)
			if bindErr != nil {
				return value{}, bindErr
			}
			if len(args) != 1 || !present[0] {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::" + node.name + "() expects exactly 1 argument"}
			}
			name, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "ReflectionClass::"+node.name, 1, "name", "string", bound[0])
			}
			for c := object.o.reflected; c != nil; c = c.parent {
				for i := range c.properties {
					property := &c.properties[i]
					if property.name != name.s || c != object.o.reflected && property.visibility == "private" {
						continue
					}
					if strings.EqualFold(node.name, "hasProperty") {
						return boolean(true), nil
					}
					e.initClasses()
					return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionproperty"], properties: &phpArray{}, reflectedProperty: property}}, nil
				}
			}
			if strings.EqualFold(node.name, "hasProperty") {
				return boolean(false), nil
			}
			return value{}, &runtimeError{node.line, "ReflectionException", "Property " + object.o.reflected.name + "::$" + name.s + " does not exist"}
		case "hasmethod", "getmethod":
			bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "name"}}, node.line)
			if bindErr != nil {
				return value{}, bindErr
			}
			if len(args) != 1 || !present[0] {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::" + node.name + "() expects exactly 1 argument"}
			}
			name, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(node.line, "ReflectionClass::"+node.name, 1, "name", "string", bound[0])
			}
			method := e.lookupMethod(object.o.reflected, name.s)
			if strings.EqualFold(node.name, "hasMethod") {
				return boolean(method != nil), nil
			}
			if method == nil {
				return value{}, &runtimeError{node.line, "ReflectionException", "Method " + object.o.reflected.name + "::" + name.s + "() does not exist"}
			}
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionmethod"], reflectedMethod: method, properties: &phpArray{}}}, nil
		case "implementsinterface", "issubclassof":
			bound, present, _, bindErr := bindCallArguments(args, []parameter{{name: "interface"}}, node.line)
			if bindErr != nil {
				return value{}, bindErr
			}
			if len(args) != 1 || !present[0] {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionClass::" + node.name + "() expects exactly 1 argument"}
			}
			if bound[0].kind != stringKind {
				return value{}, builtinTypeError(node.line, "ReflectionClass::"+node.name, 1, "interface", "string", bound[0])
			}
			if strings.EqualFold(node.name, "implementsInterface") {
				return boolean(classImplements(object.o.reflected, bound[0].s)), nil
			}
			for c := object.o.reflected.parent; c != nil; c = c.parent {
				if strings.EqualFold(c.name, bound[0].s) {
					return boolean(true), nil
				}
			}
			return boolean(false), nil
		case "getname":
			return value{kind: stringKind, s: object.o.reflected.name}, nil
		case "getfilename":
			if object.o.reflected.source.filename == "" {
				return boolean(false), nil
			}
			return value{kind: stringKind, s: object.o.reflected.source.filename}, nil
		case "getshortname":
			parts := strings.Split(object.o.reflected.name, "\\")
			return value{kind: stringKind, s: parts[len(parts)-1]}, nil
		case "isabstract":
			return boolean(object.o.reflected.abstract), nil
		case "isinstantiable":
			return boolean(object.o.reflected.kind != "interface" && object.o.reflected.kind != "trait" && object.o.reflected.kind != "enum" && !object.o.reflected.abstract), nil
		case "isinterface":
			return boolean(object.o.reflected.kind == "interface"), nil
		case "isenum":
			return boolean(object.o.reflected.kind == "enum"), nil
		case "istrait":
			return boolean(object.o.reflected.kind == "trait"), nil
		case "getconstructor":
			method := e.lookupMethod(object.o.reflected, "__construct")
			if method == nil {
				return value{kind: nullKind}, nil
			}
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionmethod"], reflectedMethod: method, properties: &phpArray{}}}, nil
		case "getattributes":
			return e.reflectionAttributes(object.o.reflected.attributes, object.o.reflected, args, node.line)
		case "getparentclass":
			if object.o.reflected.parent == nil {
				return boolean(false), nil
			}
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionclass"], reflected: object.o.reflected.parent, properties: &phpArray{}}}, nil
		}
	}
	if object.o.reflectedAttribute != nil {
		switch strings.ToLower(node.name) {
		case "getname":
			return value{kind: stringKind, s: object.o.reflectedAttribute.name}, nil
		case "getarguments":
			return e.reflectionAttributeArguments(*object.o.reflectedAttribute, object.o.attributeScope)
		case "newinstance":
			return e.reflectionAttributeInstance(*object.o.reflectedAttribute, object.o.attributeScope)
		}
	}
	if object.o.reflectedProperty != nil {
		args, argErr := e.evalArguments(node.args)
		if argErr != nil {
			return value{}, argErr
		}
		switch strings.ToLower(node.name) {
		case "setrawvalue", "getrawvalue":
			if e.version != PHP83 {
				return e.reflectionRawProperty(object.o.reflectedProperty, node.name, args, node.line)
			}
		case "getname":
			return value{kind: stringKind, s: object.o.reflectedProperty.name}, nil
		case "gettype":
			if object.o.reflectedProperty.typeName == "" {
				return value{kind: nullKind}, nil
			}
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionnamedtype"], properties: &phpArray{}, reflectedTypeName: object.o.reflectedProperty.typeName, reflectedNullable: object.o.reflectedProperty.nullable}}, nil
		case "isreadonly":
			return boolean(object.o.reflectedProperty.readonly), nil
		case "getmodifiers":
			property := object.o.reflectedProperty
			flags := int64(1)
			if property.visibility == "protected" {
				flags = 2
			} else if property.visibility == "private" {
				flags = 4
			}
			if property.static {
				flags |= 16
			}
			if property.readonly {
				flags |= 128
			}
			return integer(flags), nil
		case "getdeclaringclass":
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionclass"], properties: &phpArray{}, reflected: object.o.reflectedProperty.owner}}, nil
		case "getattributes":
			return e.reflectionAttributes(object.o.reflectedProperty.attributes, object.o.reflectedProperty.owner, args, node.line)
		case "ispublic":
			return boolean(object.o.reflectedProperty.visibility == "public"), nil
		case "isprotected":
			return boolean(object.o.reflectedProperty.visibility == "protected"), nil
		case "isprivate":
			return boolean(object.o.reflectedProperty.visibility == "private"), nil
		case "isstatic":
			return boolean(object.o.reflectedProperty.static), nil
		}
	}
	if object.o.reflectedMethod != nil {
		if object.o.class != nil && object.o.class.name == "ReflectionMethod" {
			switch strings.ToLower(node.name) {
			case "isconstructor":
				return boolean(strings.EqualFold(object.o.reflectedMethod.name, "__construct")), nil
			case "getmodifiers":
				method := object.o.reflectedMethod
				var flags int64
				switch method.visibility {
				case "private":
					flags = 4
				case "protected":
					flags = 2
				default:
					flags = 1
				}
				if method.static {
					flags |= 16
				}
				if method.final {
					flags |= 32
				}
				if method.abstract {
					flags |= 64
				}
				return integer(flags), nil
			case "getattributes":
				args, argErr := e.evalArguments(node.args)
				if argErr != nil {
					return value{}, argErr
				}
				return e.reflectionAttributes(object.o.reflectedMethod.attributes, object.o.reflectedMethod.owner, args, node.line)
			case "invoke":
				args, argErr := e.evalArguments(node.args)
				if argErr != nil {
					return value{}, argErr
				}
				if len(args) == 0 {
					return value{}, &runtimeError{node.line, "ArgumentCountError", "ReflectionMethod::invoke() expects at least 1 argument"}
				}
				method := object.o.reflectedMethod
				receiver := args[0].value
				if !method.static {
					if receiver.kind != objectKind || receiver.o == nil {
						return value{}, builtinTypeError(node.line, "ReflectionMethod::invoke", 1, "object", "?object", receiver)
					}
				} else {
					receiver = value{}
				}
				called := method.owner
				if receiver.kind == objectKind {
					called = receiver.o.class
				}
				return e.invokeWithClass(receiver, called, method, args[1:], node.line)
			case "hasreturntype":
				return boolean(object.o.reflectedMethod.returnType != ""), nil
			case "getreturntype":
				if object.o.reflectedMethod.returnType == "" {
					return value{kind: nullKind}, nil
				}
				e.initClasses()
				return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionnamedtype"], properties: &phpArray{}, reflectedTypeName: object.o.reflectedMethod.returnType, reflectedNullable: object.o.reflectedMethod.returnNullable}}, nil
			case "getdeclaringclass":
				e.initClasses()
				return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionclass"], properties: &phpArray{}, reflected: object.o.reflectedMethod.owner}}, nil
			case "getnumberofparameters":
				return integer(int64(len(object.o.reflectedMethod.parameters))), nil
			case "getnumberofrequiredparameters":
				count := int64(0)
				for _, parameter := range object.o.reflectedMethod.parameters {
					if parameter.initial == nil && !parameter.variadic {
						count++
					}
				}
				return integer(count), nil
			case "isabstract":
				return boolean(object.o.reflectedMethod.abstract), nil
			case "isfinal":
				return boolean(object.o.reflectedMethod.final), nil
			case "ispublic":
				return boolean(object.o.reflectedMethod.visibility == "public"), nil
			case "isprotected":
				return boolean(object.o.reflectedMethod.visibility == "protected"), nil
			case "isprivate":
				return boolean(object.o.reflectedMethod.visibility == "private"), nil
			case "isstatic":
				return boolean(object.o.reflectedMethod.static), nil
			}
		}
		reflectionFunction := false
		for cls := object.o.class; cls != nil; cls = cls.parent {
			if cls.name == "ReflectionFunction" {
				reflectionFunction = true
				break
			}
		}
		if reflectionFunction {
			switch strings.ToLower(node.name) {
			case "getclosureusedvariables":
				out := arrayValue()
				if object.o.reflectedClosure != nil {
					names := make([]string, 0, len(object.o.reflectedClosure.captures))
					for name := range object.o.reflectedClosure.captures {
						names = append(names, name)
					}
					sort.Strings(names)
					for _, name := range names {
						out.a = out.a.with(arrayKey{isString: true, text: name}, dereference(object.o.reflectedClosure.captures[name]))
					}
				}
				return out, nil
			case "getclosurecalledclass":
				if object.o.reflectedClosure == nil || object.o.reflectedClosure.called == nil {
					return value{kind: nullKind}, nil
				}
				e.initClasses()
				return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionclass"], reflected: object.o.reflectedClosure.called, properties: &phpArray{}}}, nil
			case "getfilename":
				if object.o.reflectedMethod == nil || object.o.reflectedMethod.source.filename == "" {
					return boolean(false), nil
				}
				return value{kind: stringKind, s: object.o.reflectedMethod.source.filename}, nil
			case "getstartline":
				if object.o.reflectedMethod == nil {
					return boolean(false), nil
				}
				return integer(int64(object.o.reflectedMethod.line)), nil
			case "isanonymous":
				return boolean(object.o.reflectedClosure != nil && object.o.reflectedClosure.target == nil), nil
			case "getclosurescopeclass":
				var scope *classDecl
				if object.o.reflectedClosure != nil {
					if object.o.reflectedClosure.target != nil {
						scope = object.o.reflectedClosure.target.called
					} else {
						scope = object.o.reflectedClosure.method.owner
					}
				}
				if scope == nil {
					return value{kind: nullKind}, nil
				}
				e.initClasses()
				return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionclass"], reflected: scope, properties: &phpArray{}}}, nil
			}
		}
		switch strings.ToLower(node.name) {
		case "getparameters":
			e.initClasses()
			out := arrayValue()
			for i := range object.o.reflectedMethod.parameters {
				param := &object.o.reflectedMethod.parameters[i]
				props := (&phpArray{}).with(arrayKey{isString: true, text: "name"}, value{kind: stringKind, s: param.name})
				wrapped := value{kind: objectKind, o: &phpObject{class: e.classes["reflectionparameter"], reflectedParameter: param, reflectedOwner: object.o.reflectedMethod.owner, properties: props}}
				out.a = out.a.with(arrayKey{index: int64(i)}, wrapped)
			}
			return out, nil
		case "getname":
			return value{kind: stringKind, s: object.o.reflectedMethod.name}, nil
		}
	}
	if object.o.reflectedParameter != nil {
		switch strings.ToLower(node.name) {
		case "hastype":
			return boolean(object.o.reflectedParameter.typeName != ""), nil
		case "getattributes":
			args, argErr := e.evalArguments(node.args)
			if argErr != nil {
				return value{}, argErr
			}
			return e.reflectionAttributes(object.o.reflectedParameter.attributes, object.o.reflectedOwner, args, node.line)
		case "getname":
			return value{kind: stringKind, s: object.o.reflectedParameter.name}, nil
		case "isvariadic":
			return boolean(object.o.reflectedParameter.variadic), nil
		case "isdefaultvalueavailable":
			return boolean(object.o.reflectedParameter.initial != nil), nil
		case "getdefaultvalue":
			if object.o.reflectedParameter.initial == nil {
				return value{}, &runtimeError{node.line, "ReflectionException", "Internal error: Failed to retrieve the default value"}
			}
			return e.eval(object.o.reflectedParameter.initial)
		case "gettype":
			param := object.o.reflectedParameter
			if param.typeName == "" {
				return value{kind: nullKind}, nil
			}
			if strings.ContainsAny(param.typeName, "|&()") {
				return value{}, &runtimeError{node.line, "Error", "gophp: Reflection of union and intersection types is not supported"}
			}
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionnamedtype"], reflectedTypeName: param.typeName, reflectedNullable: param.nullable, properties: &phpArray{}}}, nil
		case "getdeclaringclass":
			if object.o.reflectedOwner == nil {
				return value{kind: nullKind}, nil
			}
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["reflectionclass"], reflected: object.o.reflectedOwner, properties: &phpArray{}}}, nil
		}
	}
	if object.o.class != nil && object.o.class.name == "ReflectionNamedType" {
		switch strings.ToLower(node.name) {
		case "getname":
			return value{kind: stringKind, s: object.o.reflectedTypeName}, nil
		case "allowsnull":
			return boolean(object.o.reflectedNullable), nil
		case "isbuiltin":
			switch strings.ToLower(object.o.reflectedTypeName) {
			case "int", "float", "string", "bool", "array", "object", "callable", "iterable", "mixed", "null", "false", "true":
				return boolean(true), nil
			}
			return boolean(false), nil
		}
	}
	if object.o.generator != nil {
		if strings.EqualFold(node.name, "throw") {
			args, err := e.evalArguments(node.args)
			if err != nil {
				return value{}, err
			}
			bound, present, _, err := bindCallArguments(args, []parameter{{name: "exception"}}, node.line)
			if err != nil {
				return value{}, err
			}
			if len(args) != 1 || !present[0] {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "Generator::throw() expects exactly 1 argument"}
			}
			if bound[0].kind != objectKind || bound[0].o == nil || !isThrowableClass(bound[0].o.class) {
				return value{}, builtinTypeError(node.line, "Generator::throw", 1, "exception", "Throwable", bound[0])
			}
			generator := object.o.generator
			if !generator.started {
				if err := generator.step(); err != nil {
					return value{}, err
				}
			}
			injected := e.throwValue(bound[0], node.line)
			if generator.finished {
				return value{}, injected
			}
			if err := generator.stepSignal(generatorResume{err: injected}); err != nil {
				return value{}, err
			}
			if generator.finished {
				return value{kind: nullKind}, nil
			}
			return dereference(generator.current.item), nil
		}
		if strings.EqualFold(node.name, "send") {
			args, err := e.evalArguments(node.args)
			if err != nil {
				return value{}, err
			}
			bound, present, _, err := bindCallArguments(args, []parameter{{name: "value"}}, node.line)
			if err != nil {
				return value{}, err
			}
			if len(args) != 1 || !present[0] {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "Generator::send() expects exactly 1 argument"}
			}
			generator := object.o.generator
			if !generator.started {
				if err := generator.step(); err != nil {
					return value{}, err
				}
			}
			if err := generator.stepWith(bound[0]); err != nil {
				return value{}, err
			}
			if generator.finished {
				return value{kind: nullKind}, nil
			}
			return dereference(generator.current.item), nil
		}
		for _, name := range []string{"rewind", "valid", "key", "current", "next"} {
			if strings.EqualFold(node.name, name) {
				if len(node.args) != 0 {
					return value{}, &runtimeError{node.line, "ArgumentCountError", "Generator::" + name + "() expects exactly 0 arguments"}
				}
				return e.iteratorMethod(object, name, node.line)
			}
		}
		if strings.EqualFold(node.name, "getReturn") {
			if len(node.args) != 0 {
				return value{}, &runtimeError{node.line, "ArgumentCountError", "Generator::getReturn() expects exactly 0 arguments"}
			}
			generator := object.o.generator
			if !generator.started {
				if err := generator.step(); err != nil {
					return value{}, err
				}
			}
			if !generator.hasReturn {
				return value{}, &runtimeError{node.line, "Exception", "Cannot get return value of a generator that hasn't returned"}
			}
			return generator.returned, nil
		}
	}
	method := e.lookupMethod(object.o.class, node.name)
	if method == nil && object.o.exception != nil {
		if result, handled, err := e.callExceptionMethod(object, node); handled {
			return result, err
		}
	}
	var accessErr error
	if method != nil {
		accessErr = e.methodAccess(object.o.class, method, node.line)
	}
	if method == nil || accessErr != nil {
		if !strings.EqualFold(node.name, "__call") {
			if fallback := e.lookupMethod(object.o.class, "__call"); fallback != nil {
				args, err := e.evalArguments(node.args)
				if err != nil {
					return value{}, err
				}
				return e.invoke(object, fallback, []callArgument{{value: value{kind: stringKind, s: node.name}}, {value: packMagicArguments(args)}}, node.line)
			}
		}
		if accessErr != nil {
			return value{}, accessErr
		}
		return value{}, &runtimeError{node.line, "Error", "Call to undefined method " + object.o.className() + "::" + node.name + "()"}
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	if method.static {
		return e.invokeWithClass(value{}, object.o.class, method, args, node.line)
	}
	return e.invoke(object, method, args, node.line)
}

func packMagicArguments(args []callArgument) value {
	packed := arrayValue()
	for i, arg := range args {
		key := arrayKey{index: int64(i)}
		if arg.named {
			key = arrayKey{isString: true, text: arg.name}
		}
		packed.a = packed.a.with(key, arg.value)
	}
	return packed
}

func (e *evaluator) lookupMethod(class *classDecl, name string) *methodDecl {
	if class == nil {
		return nil
	}
	key := strings.ToLower(name)
	// A method inherited from the parent keeps calls to its own private
	// methods bound to that parent, even if the child reused the name.
	if e.scope != nil {
		for c := class; c != nil; c = c.parent {
			if c == e.scope {
				for _, own := range c.methodOrder {
					if strings.EqualFold(own.name, name) && own.visibility == "private" {
						return own
					}
				}
				break
			}
		}
	}
	for current := class; current != nil; current = current.parent {
		if method := current.methods[key]; method != nil {
			return method
		}
	}
	return nil
}

func (e *evaluator) methodAccess(class *classDecl, method *methodDecl, line int) error {
	owner := method.owner
	if owner == nil {
		owner = class
	}
	if method.visibility == "public" || e.scope == owner || method.visibility == "protected" && sameHierarchy(e.scope, owner) {
		return nil
	}
	from := "global scope"
	if e.scope != nil {
		from = "scope " + e.scope.name
	}
	return &runtimeError{line, "Error", "Call to " + method.visibility + " method " + owner.name + "::" + method.name + "() from " + from}
}

func sameHierarchy(a, b *classDecl) bool {
	for c := a; c != nil; c = c.parent {
		if c == b {
			return true
		}
	}
	for c := b; c != nil; c = c.parent {
		if c == a {
			return true
		}
	}
	return false
}

func classImplements(class *classDecl, name string) bool {
	if class == nil {
		return false
	}
	for _, iface := range class.interfaces {
		if equalName(iface.name, name) || classImplements(iface, name) {
			return true
		}
	}
	return classImplements(class.parent, name)
}

func propertyStorageKey(property propertyDecl) arrayKey {
	if property.visibility == "private" {
		return arrayKey{text: "\x00" + property.owner.name + "\x00" + property.name, isString: true}
	}
	return arrayKey{text: property.name, isString: true}
}

func propertyForKey(object *phpObject, key arrayKey) *propertyDecl {
	if object.class == nil {
		return nil
	}
	for i := range object.class.properties {
		p := &object.class.properties[i]
		if !p.static && propertyStorageKey(*p) == key {
			return p
		}
	}
	return nil
}

func uninitializedPropertyError(property propertyDecl, line int) error {
	kind := "Typed property"
	if property.static {
		kind = "Typed static property"
	}
	return &runtimeError{line, "Error", fmt.Sprintf("%s %s::$%s must not be accessed before initialization", kind, property.owner.name, property.name)}
}

func (e *evaluator) canInitializeReadonlyProperty(property *propertyDecl) bool {
	if e.scope == property.owner {
		return true
	}
	if e.version == PHP83 || property.visibility == "private" {
		return false
	}
	for current := e.scope; current != nil; current = current.parent {
		if current == property.owner {
			return true
		}
	}
	return false
}

func (e *evaluator) readonlyScopeError(property *propertyDecl, line int, action string) error {
	scope := "global scope"
	if e.scope != nil {
		scope = "scope " + e.scope.name
	}
	return &runtimeError{line, "Error", "Cannot " + action + " readonly property " + property.owner.name + "::$" + property.name + " from " + scope}
}

func (e *evaluator) storeObjectProperty(object *phpObject, key arrayKey, property *propertyDecl, v value, line int) error {
	if property != nil && property.setHook != nil && !(e.currentMethod == property.setHook && e.variables["this"].o == object) {
		_, err := e.invoke(value{kind: objectKind, o: object}, property.setHook, []callArgument{{value: v}}, line)
		return err
	}
	if object.enumCase != "" {
		return &runtimeError{line, "Error", "Cannot modify readonly property " + object.className() + "::$" + key.text}
	}
	if property == nil && object.class != nil && object.class.readonly {
		return &runtimeError{line, "Error", "Cannot create dynamic property " + object.class.name + "::$" + key.text}
	}
	if property != nil && (property.readonly || property.owner.readonly) {
		if _, initialized := object.properties.get(key); initialized {
			return &runtimeError{line, "Error", "Cannot modify readonly property " + property.owner.name + "::$" + property.name}
		}
		if !e.canInitializeReadonlyProperty(property) {
			return e.readonlyScopeError(property, line, "initialize")
		}
	}
	if current, exists := object.properties.getRaw(key); exists && current.kind == refKind {
		return e.writeReference(current.ref, v, line)
	}
	if property != nil && property.typeName != "" {
		converted, selected, ok := e.coerceType(v, relativeTypeName(property.typeName, property.owner, nil), property.nullable, e.strictTypes)
		if !ok {
			return &runtimeError{line, "TypeError", fmt.Sprintf("Cannot assign %s to property %s::$%s of type %s", v.typeName(), property.owner.name, property.name, declaredTypeText(property.typeName, property.nullable))}
		}
		if err := e.warnLossyIntTypeCoercion(v, selected, e.strictTypes, line); err != nil {
			return err
		}
		v = converted
	}
	object.properties = object.properties.with(key, v)
	if object.exception != nil && property != nil && property.owner != nil && property.owner.name == "Exception" {
		switch property.name {
		case "message":
			if v.kind == stringKind {
				object.exception.message = v.s
			}
		case "code":
			if v.kind == intKind {
				object.exception.code = v.i
				object.exception.codeText = ""
			} else if v.kind == stringKind {
				object.exception.codeText = v.s
			}
		}
	}
	return nil
}

// Evaluate a dynamic instance member name exactly once before resolving its
// declared slot. Static member names follow different PHP rules.
func (e *evaluator) instancePropertyName(name string, dynamic expr, line int) (string, error) {
	if dynamic == nil {
		return name, nil
	}
	v, err := e.eval(dynamic)
	if err != nil {
		return "", err
	}
	return e.toText(v, line)
}

// A method sees its declaring class's private slot even when the receiver is
// a child object with another property of the same source-level name.
func (e *evaluator) propertyKey(object *phpObject, name string, line int) (arrayKey, error) {
	plain := arrayKey{text: name, isString: true}
	if object.class == nil {
		return plain, nil
	}
	for _, property := range object.class.properties {
		if property.static {
			continue
		}
		if property.name == name && property.visibility == "private" && property.owner == e.scope {
			return propertyStorageKey(property), nil
		}
	}
	for i := len(object.class.properties) - 1; i >= 0; i-- {
		property := object.class.properties[i]
		if property.static {
			continue
		}
		if property.name != name || property.visibility == "private" {
			continue
		}
		if property.visibility == "protected" && !sameHierarchy(e.scope, property.owner) {
			return arrayKey{}, &runtimeError{line, "Error", "Cannot access protected property " + object.className() + "::$" + name}
		}
		return plain, nil
	}
	for _, property := range object.class.properties {
		if property.static {
			continue
		}
		if property.name == name && property.visibility == "private" {
			// A parent's private member is not visible in a child method. PHP
			// treats a write from that scope as a new dynamic property.
			for scope := e.scope; scope != nil && scope != property.owner; scope = scope.parent {
				if scope.parent == property.owner {
					return plain, nil
				}
			}
			if object.properties.find(plain) >= 0 {
				return plain, nil
			}
			return arrayKey{}, &runtimeError{line, "Error", "Cannot access private property " + object.className() + "::$" + name}
		}
	}
	return plain, nil
}

func (e *evaluator) invoke(object value, method *methodDecl, args []callArgument, line int) (result value, err error) {
	var called *classDecl
	if object.kind == objectKind {
		called = object.o.class
	}
	return e.invokeWithClass(object, called, method, args, line)
}

func (e *evaluator) invokeWithClass(object value, called *classDecl, method *methodDecl, args []callArgument, line int) (result value, err error) {
	return e.invokeWithCapture(object, called, method, args, line, nil)
}

func (e *evaluator) invokeWithCapture(object value, called *classDecl, method *methodDecl, args []callArgument, line int, closure *phpClosure) (result value, err error) {
	if method.owner != nil && method.owner.name == "ReflectionFunction" && strings.EqualFold(method.name, "__construct") {
		if len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "ReflectionFunction::__construct() expects exactly 1 argument"}
		}
		reflected, reflectErr := e.newInstance(newObject{name: "ReflectionFunction", args: []expr{literal{value: args[0].value, line: line}}, line: line})
		if reflectErr != nil {
			return value{}, reflectErr
		}
		object.o.reflectedMethod = reflected.o.reflectedMethod
		object.o.reflectedClosure = reflected.o.reflectedClosure
		return value{kind: nullKind}, nil
	}
	if method.owner != nil && method.owner.name == "DatePeriod" {
		return value{}, &runtimeError{line, "Error", "gophp: DatePeriod execution is not implemented"}
	}
	if method.owner != nil && (method.owner.name == "PDO" || method.owner.name == "PDOStatement") {
		return e.invokePDO(object, method, args, line)
	}
	if method.owner != nil && (method.owner.name == "ArrayObject" || method.owner.name == "ArrayIterator") {
		return e.invokeArrayObject(object, method, args, line)
	}
	if method.owner != nil && (method.owner.name == "Exception" || method.owner.name == "Error") && strings.EqualFold(method.name, "__toString") && object.kind == objectKind && object.o != nil && object.o.exception != nil {
		return value{kind: stringKind, s: exceptionDiagnostic(object)}, nil
	}
	if method.owner != nil && method.owner.name == "Fiber" && strings.EqualFold(method.name, "getCurrent") {
		if len(args) != 0 {
			return value{}, &runtimeError{line, "ArgumentCountError", "Fiber::getCurrent() expects exactly 0 arguments"}
		}
		return value{kind: nullKind}, nil
	}
	if object.kind == objectKind && object.o != nil && object.o.class != nil && strings.EqualFold(object.o.class.name, "RegexIterator") {
		return e.invokeRegexIterator(object, method, args, line)
	}
	if method.owner != nil && method.owner.name == "WeakMap" {
		return e.invokeWeakMap(object, method, args, line)
	}
	if method.owner != nil && method.owner.name == "SplObjectStorage" {
		return e.invokeObjectStorage(object, method, args, line)
	}
	if method.owner != nil && method.owner.name == "SplPriorityQueue" {
		return e.invokePriorityQueue(object, method, args, line)
	}
	if method.owner != nil && method.owner.name == "FilterIterator" {
		return e.invokeFilterIterator(object, method, args, line)
	}
	if method.owner != nil && (method.owner.name == "SplFileInfo" || method.owner.name == "DirectoryIterator" || method.owner.name == "FilesystemIterator" || method.owner.name == "RecursiveDirectoryIterator") {
		return e.invokeSPLDirectory(object, method, args, line)
	}
	if method.owner != nil && method.owner.name == "RecursiveIteratorIterator" {
		return e.invokeRecursiveIteratorIterator(object, method, args, line)
	}
	if method.owner != nil && (method.owner.name == "IteratorIterator" || method.owner.name == "LimitIterator") {
		result, invokeErr := e.invokeIteratorWrapper(object, method, args, line)
		if invokeErr != nil && e.failureFrames == nil {
			for _, frame := range e.frames {
				e.failureFrames = append(e.failureFrames, snapshotCallFrame(frame))
			}
			e.failureFrames = append(e.failureFrames, callFrame{class: method.owner.name, method: method.name, object: object, filename: e.filename, line: line, args: argumentValues(args)})
		}
		return result, invokeErr
	}
	if method.owner != nil && method.owner.name == "AppendIterator" {
		return e.invokeAppendIterator(object, method, args, line)
	}
	if method.owner != nil && method.owner.name == "CachingIterator" {
		return e.invokeCachingIterator(object, method, args, line)
	}
	if method.owner != nil && (method.owner.name == "DateTime" || method.owner.name == "DateTimeImmutable") {
		return e.invokeDateTime(object, called, method, args, line)
	}
	if method.owner != nil && method.owner.name == "DateInterval" {
		return e.invokeDateInterval(object, args, line)
	}
	if method.owner != nil && method.owner.name == "DateTimeZone" {
		return e.invokeDateTimeZone(object, method, args, line)
	}
	isGenerator := methodContainsYield(method.body)
	if isGenerator && e.generatorEntering != method {
		if _, _, _, bindErr := bindCallArguments(args, method.parameters, line); bindErr != nil {
			return value{}, bindErr
		}
		worker := *e
		worker.frames = append([]callFrame(nil), e.frames...)
		worker.suspendedScopes = append([]variableScope(nil), e.suspendedScopes...)
		worker.generatorRunning = true
		worker.generatorEntering = method
		generator := newPHPGenerator(func(g *phpGenerator) (value, error) {
			worker.activeGenerator = g
			result, invokeErr := worker.invokeWithCapture(object, called, method, args, line, closure)
			if invokeErr != nil {
				g.failureFrames = append([]callFrame(nil), worker.failureFrames...)
			}
			return result, invokeErr
		})
		generator.yieldsByReference = method.returnsByReference
		if e.generators != nil {
			e.generators.add(generator)
		}
		e.initClasses()
		result := value{kind: objectKind, o: &phpObject{class: e.classes["generator"], generator: generator, properties: &phpArray{}}}
		if method.returnType != "" {
			if _, _, ok := e.coerceType(result, method.returnType, method.returnNullable, e.strictTypes); !ok {
				return value{}, &runtimeError{method.line, "TypeError", method.name + "(): Return value must be of type " + declaredTypeText(method.returnType, method.returnNullable) + ", Generator returned"}
			}
		}
		return result, nil
	}
	if e.generatorEntering == method {
		e.generatorEntering = nil
	}
	// Bound recursive calls so untrusted PHP cannot overflow the Go stack.
	if len(e.frames) >= 256 {
		kind := "method"
		if object.kind != objectKind {
			kind = "function"
		}
		return value{}, &runtimeError{line, "Error", "gophp: maximum " + kind + " call depth exceeded (256)"}
	}
	// Named-argument binding is part of preparing the call. A failure here
	// belongs to the caller: PHP has not entered the callee frame yet.
	bound, present, extra, bindErr := bindCallArguments(args, method.parameters, line)
	if bindErr != nil {
		return value{}, bindErr
	}
	refCells, bindErr := e.bindVariableReferenceArguments(args, method.parameters, bound)
	if bindErr != nil {
		return value{}, bindErr
	}
	outer, scope, lateClass, sourceFilename, outerBindings := e.variables, e.scope, e.lateClass, e.filename, e.bindings
	callerStrict := e.strictTypes
	outerStaticScope, outerMethod := e.staticScope, e.currentMethod
	e.suspendedScopes = append(e.suspendedScopes, outer)
	e.variables = make(map[string]value)
	if closure != nil {
		for name, v := range closure.captures {
			e.variables[name] = v
		}
	}
	e.bindings = nil
	if e.staticScopes == nil {
		e.staticScopes = make(map[*methodDecl]variableScope)
	}
	if e.staticScopes[method] == nil {
		e.staticScopes[method] = make(variableScope)
	}
	e.staticScope, e.currentMethod = e.staticScopes[method], method
	e.scope = nil
	e.lateClass = called
	if method.source.filename != "" {
		e.filename = method.source.filename
	}
	e.strictTypes = method.source.strictTypes
	className, callableName := "", method.name
	if object.kind == objectKind {
		e.variables["this"] = object
		className = object.o.className()
		callableName = className + "::" + method.name
	} else if called != nil {
		className = called.name
		callableName = className + "::" + method.name
	}
	if called != nil {
		e.scope = method.owner
	}
	traceMethod := method.name
	if method.name == "{closure}" && e.version != PHP83 {
		traceMethod = fmt.Sprintf("{closure:%s:%d}", method.source.filename, method.line)
		callableName = traceMethod
	}
	traceClass := className
	if method.owner != nil {
		traceClass = method.owner.name
	}
	e.frames = append(e.frames, callFrame{
		class: traceClass, method: traceMethod, static: method.static && called != nil, line: line,
		object:   object,
		filename: sourceFilename,
		args:     argumentValues(args), callArgs: args, boundArgs: bound, present: present,
		parameters: method.parameters, locals: e.variables,
	})
	defer func() {
		if err != nil && e.failureFrames == nil {
			for _, frame := range e.frames {
				e.failureFrames = append(e.failureFrames, snapshotCallFrame(frame))
			}
		}
		e.frames = e.frames[:len(e.frames)-1]
		if err != nil {
			err = withSource(err, e.filename)
		}
		e.variables, e.scope, e.lateClass, e.filename, e.bindings = outer, scope, lateClass, sourceFilename, outerBindings
		e.suspendedScopes = e.suspendedScopes[:len(e.suspendedScopes)-1]
		e.strictTypes = callerStrict
		e.staticScope, e.currentMethod = outerStaticScope, outerMethod
	}()
	required := requiredParameters(method)
	missingRequired := false
	for i := 0; i < required; i++ {
		if !present[i] {
			missingRequired = true
			break
		}
	}
	if missingRequired {
		count := "exactly"
		if required < len(method.parameters) {
			count = "at least"
		}
		return value{}, &runtimeError{method.line, "ArgumentCountError", fmt.Sprintf("Too few arguments to function %s(), %d passed in %s on line %d and %s %d expected", callableName, len(args), sourceFilename, line, count, required)}
	}
	for i, parameter := range method.parameters {
		if parameter.variadic {
			packed := arrayValue()
			for n, entry := range extra.a.entries {
				v := dereference(entry.value)
				if parameter.typeName != "" {
					coerced, selected, ok := e.coerceType(v, parameter.typeName, parameter.nullable, callerStrict)
					if !ok {
						return value{}, &runtimeError{method.line, "TypeError", fmt.Sprintf("%s(): Argument #%d must be of type %s, %s given, called in %s on line %d and defined", callableName, i+n+1, declaredTypeText(parameter.typeName, parameter.nullable), v.diagnosticName(), sourceFilename, line)}
					}
					if err := e.warnLossyIntTypeCoercion(v, selected, callerStrict, method.line); err != nil {
						return value{}, err
					}
					v = coerced
				}
				packed.a = packed.a.with(entry.key, v)
			}
			e.variables[parameter.name] = packed
			// Trace arguments are the values passed to the call. The variadic
			// local is an array and may be changed before an exception is made.
			e.frames[len(e.frames)-1].boundArgs[i] = packed
			break
		}
		v := value{}
		if present[i] {
			v = bound[i]
		} else {
			v, err = e.eval(parameter.initial)
			if err != nil {
				return value{}, err
			}
		}
		if parameter.typeName != "" {
			coerced, selected, ok := e.coerceType(v, parameter.typeName, parameter.nullable, callerStrict)
			if !ok {
				suffix := " and defined"
				if method.propertyHook {
					suffix = ""
				}
				return value{}, &runtimeError{method.line, "TypeError", fmt.Sprintf("%s(): Argument #%d ($%s) must be of type %s, %s given, called in %s on line %d%s", callableName, i+1, parameter.name, declaredTypeText(parameter.typeName, parameter.nullable), v.diagnosticName(), sourceFilename, line, suffix)}
			}
			if present[i] {
				if err := e.warnLossyIntTypeCoercion(v, selected, callerStrict, method.line); err != nil {
					return value{}, err
				}
			}
			v = coerced
		}
		if refCells[i] != nil {
			if err := e.writeReference(refCells[i], v, method.line); err != nil {
				return value{}, err
			}
			e.variables[parameter.name] = value{kind: refKind, ref: refCells[i]}
		} else {
			e.variables[parameter.name] = v
		}
	}
	if object.kind == objectKind {
		for _, parameter := range method.parameters {
			if !parameter.promoted {
				continue
			}
			property := propertyForKey(object.o, propertyStorageKey(propertyDecl{name: parameter.name, visibility: parameter.visibility, owner: method.owner}))
			if property == nil {
				return value{}, &runtimeError{method.line, "Error", "Promoted property $" + parameter.name + " not found"}
			}
			key := propertyStorageKey(*property)
			if err := e.storeObjectProperty(object.o, key, property, e.variables[parameter.name], method.line); err != nil {
				return value{}, err
			}
		}
	}
	e.frames[len(e.frames)-1].traceReady = true
	err = e.execute(method.body)
	if r, ok := err.(*returned); ok {
		if isGenerator && e.activeGenerator != nil {
			return r.value, nil
		}
		if method.returnsByReference {
			if r.value.kind != refKind {
				return value{}, &runtimeError{r.line, "Error", "gophp: return by reference requires a variable"}
			}
			actual := dereference(r.value)
			if method.returnType != "" {
				coerced, selected, ok := e.coerceType(actual, method.returnType, method.returnNullable, e.strictTypes)
				if !ok {
					return value{}, &runtimeError{r.line, "TypeError", fmt.Sprintf("%s(): Return value must be of type %s, %s returned", callableName, declaredTypeText(method.returnType, method.returnNullable), actual.diagnosticName())}
				}
				if err := e.warnLossyIntTypeCoercion(actual, selected, e.strictTypes, r.line); err != nil {
					return value{}, err
				}
				if err := e.writeReference(r.value.ref, coerced, r.line); err != nil {
					return value{}, err
				}
			}
			return r.value, nil
		}
		if method.returnType == "never" && method.arrow {
			return value{}, &runtimeError{r.line, "TypeError", callableName + "(): never-returning function must not implicitly return"}
		}
		if method.returnType == "void" {
			return value{}, nil
		}
		if method.returnType != "" {
			coerced, selected, ok := e.coerceType(r.value, method.returnType, method.returnNullable, e.strictTypes)
			if !ok {
				return value{}, &runtimeError{r.line, "TypeError", fmt.Sprintf("%s(): Return value must be of type %s, %s returned", callableName, declaredTypeText(method.returnType, method.returnNullable), r.value.diagnosticName())}
			}
			if err := e.warnLossyIntTypeCoercion(r.value, selected, e.strictTypes, r.line); err != nil {
				return value{}, err
			}
			return coerced, nil
		}
		return r.value, nil
	}
	if isGenerator && e.activeGenerator != nil && err == nil {
		return value{kind: nullKind}, nil
	}
	if err == nil && method.returnType == "never" {
		kind := "function"
		if method.owner != nil && e.version != PHP83 {
			kind = "method"
		}
		return value{}, &runtimeError{method.line, "TypeError", callableName + "(): never-returning " + kind + " must not implicitly return"}
	}
	if err == nil && method.returnType != "" && method.returnType != "void" {
		return value{}, &runtimeError{method.line, "TypeError", fmt.Sprintf("%s(): Return value must be of type %s, none returned", callableName, declaredTypeText(method.returnType, method.returnNullable))}
	}
	return value{}, err
}

func traceValue(v value) string {
	switch v.kind {
	case nullKind:
		return "NULL"
	case boolKind:
		if v.i != 0 {
			return "true"
		}
		return "false"
	case arrayKind:
		return "Array"
	case objectKind:
		return "Object(" + v.o.className() + ")"
	case floatKind:
		s := v.text()
		if !math.IsNaN(v.f) && !math.IsInf(v.f, 0) && !strings.Contains(s, ".") {
			s += ".0"
		}
		return s
	case stringKind:
		s := v.s
		if len(s) > 15 {
			s = s[:15] + "..."
		}
		var out strings.Builder
		out.WriteByte('\'')
		for i := 0; i < len(s); i++ {
			c := s[i]
			switch c {
			case '\n':
				out.WriteString(`\n`)
			case '\r':
				out.WriteString(`\r`)
			case '\t':
				out.WriteString(`\t`)
			case '\f':
				out.WriteString(`\f`)
			case '\v':
				out.WriteString(`\v`)
			case 27:
				out.WriteString(`\e`)
			case '\\':
				out.WriteString(`\\`)
			default:
				if c < 32 || c >= 127 {
					fmt.Fprintf(&out, "\\x%02X", c)
				} else {
					out.WriteByte(c)
				}
			}
		}
		out.WriteByte('\'')
		return out.String()
	}
	return v.text()
}

func formatStackTrace(frames []callFrame) string {
	var result strings.Builder
	for i := len(frames) - 1; i >= 0; i-- {
		frame := frames[i]
		args := make([]string, len(frame.args))
		for n, arg := range frame.args {
			args[n] = traceValue(arg)
			if n < len(frame.argKeys) && frame.argKeys[n].isString {
				args[n] = frame.argKeys[n].text + ": " + args[n]
			}
		}
		name := frame.method
		if frame.class != "" {
			separator := "->"
			if frame.static {
				separator = "::"
			}
			name = frame.class + separator + name
		}
		fmt.Fprintf(&result, "#%d %s(%d): %s(%s)\n", len(frames)-i-1, frame.filename, frame.line, name, strings.Join(args, ", "))
	}
	fmt.Fprintf(&result, "#%d {main}\n", len(frames))
	return result.String()
}

func (e *evaluator) stackTrace() string { return formatStackTrace(e.failureFrames) }
