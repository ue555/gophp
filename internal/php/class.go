package php

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

type propertyDecl struct {
	name, visibility string
	initial          expr
	static           bool
	fromTrait        bool
	traitSource      string
	line             int
	owner            *classDecl
}
type parameter struct {
	name     string
	initial  expr
	line     int
	typeName string
	nullable bool
	variadic bool
}
type methodDecl struct {
	name, visibility        string
	final, abstract, static bool
	traitSource             string
	parameters              []parameter
	returnType              string
	returnNullable          bool
	body                    []statement
	line                    int
	source                  sourceContext
	owner                   *classDecl
}
type classConstantDecl struct {
	name, visibility string
	initial          expr
	final            bool
	traitSource      string
	line             int
	owner            *classDecl
}
type classConstantValue struct {
	value value
	decl  classConstantDecl
}
type staticPropertySlot struct {
	value value
	decl  propertyDecl
}
type classDecl struct {
	name             string
	kind             string // class, interface, or trait
	final, abstract  bool
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

func (p *parser) classDeclaration(kind string) (*classDecl, error) {
	name := p.current()
	if err := p.expect("name"); err != nil {
		return nil, err
	}
	if err := p.registerDeclaration(name, "class"); err != nil {
		return nil, err
	}
	decl := &classDecl{name: qualify(p.namespace, name.text), kind: kind, line: name.line, methods: make(map[string]*methodDecl), source: sourceContext{p.filename, p.namespace, p.strictTypes}}
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
	if kind == "class" && p.acceptKeyword("implements") {
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
	defer func() { p.classDepth-- }()
	for p.current().kind != "}" {
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
		for {
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
			if methodAbstract || methodStatic {
				return nil, syntax(p.current().line, "unsupported class constant modifier")
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
				decl.constantDecls = append(decl.constantDecls, classConstantDecl{name: name.text, visibility: visibility, initial: initial, final: methodFinal, line: name.line, owner: decl})
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
			name := p.current()
			if err := p.expect("name"); err != nil {
				return nil, err
			}
			lower := strings.ToLower(name.text)
			if strings.HasPrefix(lower, "__") && lower != "__construct" {
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
			method := &methodDecl{name: name.text, visibility: visibility, final: methodFinal, abstract: methodAbstract || kind == "interface", static: methodStatic, line: name.line, owner: decl}
			var err error
			method.source = sourceContext{p.filename, p.namespace, p.strictTypes}
			if method.abstract {
				method.parameters, err = p.parameterList()
				if err == nil {
					method.returnType, method.returnNullable, err = p.returnDeclaration()
				}
				if err == nil {
					err = p.expect(";")
				}
			} else {
				method.parameters, method.body, method.returnType, method.returnNullable, err = p.callableBody()
			}
			if err != nil {
				return nil, err
			}
			decl.methods[lower] = method
			decl.methodOrder = append(decl.methodOrder, method)
			continue
		}
		if kind == "interface" || methodFinal || methodAbstract || (!explicit && !methodStatic) || p.current().kind != "variable" {
			return nil, syntax(p.current().line, "unsupported class member %q", p.current().text)
		}
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
			property := propertyDecl{name: name.text, line: name.line, visibility: visibility, static: methodStatic, owner: decl}
			if p.accept("=") {
				var err error
				property.initial, err = p.expression(0)
				if err != nil {
					return nil, err
				}
			}
			decl.properties = append(decl.properties, property)
			if !p.accept(",") {
				break
			}
		}
		if err := p.expect(";"); err != nil {
			return nil, err
		}
	}
	return decl, p.expect("}")
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
	line          int
	filename      string
	args          []value
	parameters    []parameter
	locals        variableScope
}

func (e *evaluator) initClasses() {
	if e.classes == nil {
		throwable := &classDecl{name: "Throwable", kind: "interface"}
		exception := &classDecl{name: "Exception", kind: "class", interfaces: []*classDecl{throwable}}
		baseError := &classDecl{name: "Error", kind: "class", interfaces: []*classDecl{throwable}}
		e.classes = map[string]*classDecl{
			"stdclass": {name: "stdClass"}, "closure": {name: "Closure", final: true},
			"throwable": throwable, "exception": exception, "error": baseError,
			"typeerror":           {name: "TypeError", kind: "class", parent: baseError},
			"argumentcounterror":  {name: "ArgumentCountError", kind: "class", parent: baseError},
			"valueerror":          {name: "ValueError", kind: "class", parent: baseError},
			"divisionbyzeroerror": {name: "DivisionByZeroError", kind: "class", parent: baseError},
			"unhandledmatcherror": {name: "UnhandledMatchError", kind: "class", parent: baseError},
		}
	}
}

func (e *evaluator) declareClass(decl *classDecl) error {
	e.initClasses()
	key := strings.ToLower(decl.name)
	if _, exists := e.classes[key]; exists {
		return &runtimeError{decl.line, "", "Cannot declare class " + decl.name + ", because the name is already in use"}
	}
	if decl.final && decl.abstract {
		return &runtimeError{decl.line, "", "Cannot use the final modifier on an abstract class"}
	}
	for _, name := range decl.interfaceNames {
		iface := e.classes[strings.ToLower(name)]
		if iface == nil || iface.kind != "interface" {
			return &runtimeError{decl.line, "Error", "Interface \"" + name + "\" not found"}
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
	if decl.kind == "interface" {
		methods := make(map[string]*methodDecl)
		for _, iface := range decl.interfaces {
			for _, method := range interfaceContracts(iface) {
				name := strings.ToLower(method.name)
				if existing := methods[name]; existing != nil && (existing.static != method.static || !e.compatibleMethodDeclaration(existing, method)) {
					return &fatalDiagnostic{existing.line, "Declaration of " + methodSignature(existing) + " must be compatible with " + methodSignature(method)}
				}
				methods[name] = method
			}
		}
		for name, method := range decl.methods {
			for _, iface := range decl.interfaces {
				for _, inherited := range interfaceContracts(iface) {
					if strings.EqualFold(name, inherited.name) && (method.static != inherited.static || !e.compatibleMethodDeclaration(method, inherited)) {
						return &fatalDiagnostic{method.line, "Declaration of " + methodSignature(method) + " must be compatible with " + methodSignature(inherited)}
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
			return &runtimeError{decl.line, "Error", "Class \"" + decl.parentName + "\" not found"}
		}
		if parent.kind == "interface" {
			return &fatalDiagnostic{decl.line, "Class " + decl.name + " cannot extend interface " + parent.name}
		}
		if parent.final {
			return &fatalDiagnostic{decl.line, "Class " + decl.name + " cannot extend final class " + parent.name}
		}
		decl.parent = parent
		properties := append([]propertyDecl(nil), parent.properties...)
		for _, property := range decl.properties {
			replaced := false
			for i, inherited := range parent.properties {
				if property.name == inherited.name && inherited.visibility != "private" {
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
				return &fatalDiagnostic{method.line, "Declaration of " + methodSignature(method) + " must be compatible with " + methodSignature(inherited)}
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
	for _, iface := range decl.interfaces {
		for _, required := range interfaceContracts(iface) {
			name := strings.ToLower(required.name)
			method := decl.methods[name]
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
				return &fatalDiagnostic{method.line, "Declaration of " + methodSignature(method) + " must be compatible with " + methodSignature(required)}
			}
		}
	}
	e.classes[key] = decl
	if err := e.initClassMembers(decl); err != nil {
		delete(e.classes, key)
		return err
	}
	return nil
}

// Interface lookup keeps one method per name, but declaration checking must
// retain every inherited contract. A class implementing a diamond may need
// to satisfy two signatures even when its interface's lookup map has only one.
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
		result = append(result, current.methodOrder...)
	}
	visit(iface)
	return result
}

func compatibleUntypedMethod(implementation, contract *methodDecl) bool {
	if requiredParameters(implementation) > requiredParameters(contract) {
		return false
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
		if parameter.typeName != "" {
			parts[i] = declaredTypeText(parameter.typeName, parameter.nullable) + " " + parts[i]
		}
	}
	signature := method.owner.name + "::" + method.name + "(" + strings.Join(parts, ", ") + ")"
	if method.returnType != "" {
		signature += ": " + declaredTypeText(method.returnType, method.returnNullable)
	}
	return signature
}

func (e *evaluator) initClassMembers(decl *classDecl) error {
	oldScope := e.scope
	e.scope = decl
	defer func() { e.scope = oldScope }()
	decl.constants = make(map[string]classConstantValue)
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
			return err
		}
		decl.constants[constant.name] = classConstantValue{v, constant}
	}
	for _, property := range decl.properties {
		if !property.static || property.owner != decl {
			continue
		}
		var v value
		if property.initial != nil {
			initial, err := e.eval(property.initial)
			if err != nil {
				return err
			}
			v = initial
		}
		decl.staticProperties[property.name] = &staticPropertySlot{v, property}
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
							return &fatalDiagnostic{decl.line, fmt.Sprintf("%s and %s define the same property ($%s) in the composition of %s", decl.name, trait.name, property.name, decl.name)}
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
					return &fatalDiagnostic{decl.line, fmt.Sprintf("Trait method %s::%s has not been applied as %s::%s, because of collision with %s::%s", selectedTrait.name, selected.name, decl.name, selected.name, existing.traitSource, existing.name)}
				}
				copy := *selected
				copy.owner, copy.traitSource = decl, selectedTrait.name
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
			if _, own := decl.methods[strings.ToLower(name)]; own && adapt.alias != "" {
				continue // A class member takes precedence over a trait alias.
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
	if a.visibility != b.visibility || a.static != b.static {
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
	if a.visibility != b.visibility || a.final != b.final {
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
		return nil, &runtimeError{line, "Error", fmt.Sprintf("Class %q not found", name)}
	}
	return class, nil
}

func (e *evaluator) readClassConstant(node classConstantExpr) (value, error) {
	class, err := e.resolveClassName(node.className, node.line)
	if err != nil {
		return value{}, err
	}
	constant, ok := class.constants[node.name]
	if !ok {
		return value{}, &runtimeError{node.line, "Error", "Undefined constant " + class.name + "::" + node.name}
	}
	owner := constant.decl.owner
	if constant.decl.visibility == "private" && e.scope != owner || constant.decl.visibility == "protected" && !sameHierarchy(e.scope, owner) {
		return value{}, &runtimeError{node.line, "Error", "Cannot access " + constant.decl.visibility + " constant " + class.name + "::" + node.name}
	}
	return constant.value, nil
}

func (e *evaluator) staticPropertyLocation(node staticPropertyExpr) (location, error) {
	class, err := e.resolveClassName(node.className, node.line)
	if err != nil {
		return location{}, err
	}
	slot := class.staticProperties[node.name]
	if slot == nil {
		return location{}, &runtimeError{node.line, "Error", "Access to undeclared static property " + class.name + "::$" + node.name}
	}
	owner := slot.decl.owner
	if slot.decl.visibility == "private" && e.scope != owner || slot.decl.visibility == "protected" && !sameHierarchy(e.scope, owner) {
		return location{}, &runtimeError{node.line, "Error", "Cannot access " + slot.decl.visibility + " property " + class.name + "::$" + node.name}
	}
	return location{func() value { return slot.value }, func(v value) error { slot.value = v; return nil }}, nil
}

func (e *evaluator) callStaticMethod(node staticCallExpr) (value, error) {
	class, err := e.resolveClassName(node.className, node.line)
	if err != nil {
		return value{}, err
	}
	method := e.lookupMethod(class, node.name)
	if method == nil {
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
	if node.classExpr != nil {
		v, evalErr := e.eval(node.classExpr)
		if evalErr != nil {
			return value{}, evalErr
		}
		if v.kind != stringKind {
			return value{}, &runtimeError{node.line, "Error", "Class name must be a valid object or a string"}
		}
		className = strings.TrimPrefix(v.s, "\\")
		e.initClasses()
		class = e.classes[strings.ToLower(className)]
		if class == nil {
			return value{}, &runtimeError{node.line, "Error", fmt.Sprintf("Class %q not found", className)}
		}
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
	if strings.EqualFold(class.name, "Closure") {
		return value{}, &runtimeError{node.line, "Error", "Instantiation of class Closure is not allowed"}
	}
	if class.abstract {
		return value{}, &runtimeError{node.line, "Error", "Cannot instantiate abstract class " + class.name}
	}
	object := value{kind: objectKind, o: &phpObject{class: class, properties: &phpArray{}}}
	if isThrowableClass(class) {
		object.o.exception = &exceptionData{file: e.filename, line: node.line}
	}
	for _, property := range class.properties {
		if property.static {
			continue
		}
		var v value
		var err error
		if property.initial != nil {
			v, err = e.eval(property.initial)
		}
		if err != nil {
			return value{}, err
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
	name  string // Empty for a positional argument.
	named bool
	value value
	line  int
}

func argumentValues(args []callArgument) []value {
	values := make([]value, len(args))
	for i, arg := range args {
		values[i] = arg.value
	}
	return values
}

func (e *evaluator) evalArguments(args []expr) ([]callArgument, error) {
	values := make([]callArgument, 0, len(args))
	seenNamed := false
	for _, arg := range args {
		if unpack, ok := arg.(unpackArgument); ok {
			v, err := e.eval(unpack.operand)
			if err != nil {
				return nil, err
			}
			if v.kind != arrayKind {
				return nil, &runtimeError{unpack.line, "Error", "Only arrays and Traversables can be unpacked"}
			}
			for _, entry := range v.a.entries {
				if entry.key.isString {
					seenNamed = true
					values = append(values, callArgument{name: entry.key.text, named: true, value: entry.value, line: unpack.line})
					continue
				}
				if seenNamed {
					return nil, &runtimeError{unpack.line, "Error", "Cannot use positional argument after named argument"}
				}
				values = append(values, callArgument{value: entry.value, line: unpack.line})
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
		v, err := e.eval(operand)
		if err != nil {
			return nil, err
		}
		values = append(values, callArgument{name: name, named: isNamed, value: v, line: expressionLine(arg)})
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

func (e *evaluator) callMethod(node methodCall) (value, error) {
	object, err := e.eval(node.receiver)
	if err != nil {
		return value{}, err
	}
	if object.kind != objectKind {
		return value{}, &runtimeError{node.line, "Error", "Call to a member function " + node.name + "() on " + object.typeName()}
	}
	method := e.lookupMethod(object.o.class, node.name)
	if method == nil && object.o.exception != nil {
		if result, handled, err := e.callExceptionMethod(object, node); handled {
			return result, err
		}
	}
	if method == nil {
		return value{}, &runtimeError{node.line, "Error", "Call to undefined method " + object.o.className() + "::" + node.name + "()"}
	}
	if err := e.methodAccess(object.o.class, method, node.line); err != nil {
		return value{}, err
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
	return class.methods[key]
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
	for _, iface := range class.interfaces {
		if equalName(iface.name, name) || classImplements(iface, name) {
			return true
		}
	}
	return false
}

func propertyStorageKey(property propertyDecl) arrayKey {
	if property.visibility == "private" {
		return arrayKey{text: "\x00" + property.owner.name + "\x00" + property.name, isString: true}
	}
	return arrayKey{text: property.name, isString: true}
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
	outer, scope, lateClass, sourceFilename, outerBindings := e.variables, e.scope, e.lateClass, e.filename, e.bindings
	callerStrict := e.strictTypes
	outerStaticScope, outerMethod := e.staticScope, e.currentMethod
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
	e.frames = append(e.frames, callFrame{
		class: className, method: method.name, static: method.static && called != nil, line: line,
		filename: sourceFilename,
		args:     argumentValues(args), parameters: method.parameters, locals: e.variables,
	})
	defer func() {
		if err != nil && e.failureFrames == nil {
			for _, frame := range e.frames {
				frame.args = append([]value(nil), frame.args...)
				// PHP traces contain the current values of passed parameters.
				for i, parameter := range frame.parameters {
					if i >= len(frame.args) {
						break
					}
					if v, bound := frame.locals[parameter.name]; bound {
						frame.args[i] = v
					}
				}
				frame.locals = nil
				e.failureFrames = append(e.failureFrames, frame)
			}
		}
		e.frames = e.frames[:len(e.frames)-1]
		if err != nil {
			err = withSource(err, e.filename)
		}
		e.variables, e.scope, e.lateClass, e.filename, e.bindings = outer, scope, lateClass, sourceFilename, outerBindings
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
				v := entry.value
				if parameter.typeName != "" {
					coerced, ok := e.coerceDeclaredType(v, parameter.typeName, parameter.nullable, callerStrict)
					if !ok {
						return value{}, &runtimeError{method.line, "TypeError", fmt.Sprintf("%s(): Argument #%d must be of type %s, %s given, called in %s on line %d and defined", callableName, i+n+1, declaredTypeText(parameter.typeName, parameter.nullable), v.diagnosticName(), sourceFilename, line)}
					}
					v = coerced
				}
				packed.a = packed.a.with(entry.key, v)
			}
			e.variables[parameter.name] = packed
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
			coerced, ok := e.coerceDeclaredType(v, parameter.typeName, parameter.nullable, callerStrict)
			if !ok {
				return value{}, &runtimeError{method.line, "TypeError", fmt.Sprintf("%s(): Argument #%d ($%s) must be of type %s, %s given, called in %s on line %d and defined", callableName, i+1, parameter.name, declaredTypeText(parameter.typeName, parameter.nullable), v.diagnosticName(), sourceFilename, line)}
			}
			v = coerced
		}
		e.variables[parameter.name] = v
	}
	err = e.execute(method.body)
	if r, ok := err.(*returned); ok {
		if method.returnType != "" {
			coerced, ok := e.coerceDeclaredType(r.value, method.returnType, method.returnNullable, e.strictTypes)
			if !ok {
				return value{}, &runtimeError{r.line, "TypeError", fmt.Sprintf("%s(): Return value must be of type %s, %s returned", callableName, declaredTypeText(method.returnType, method.returnNullable), r.value.diagnosticName())}
			}
			return coerced, nil
		}
		return r.value, nil
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

func (e *evaluator) stackTrace() string {
	var result strings.Builder
	for i := len(e.failureFrames) - 1; i >= 0; i-- {
		frame := e.failureFrames[i]
		args := make([]string, len(frame.args))
		for n, arg := range frame.args {
			args[n] = traceValue(arg)
		}
		name := frame.method
		if frame.class != "" {
			separator := "->"
			if frame.static {
				separator = "::"
			}
			name = frame.class + separator + name
		}
		fmt.Fprintf(&result, "#%d %s(%d): %s(%s)\n", len(e.failureFrames)-i-1, frame.filename, frame.line, name, strings.Join(args, ", "))
	}
	fmt.Fprintf(&result, "#%d {main}\n", len(e.failureFrames))
	return result.String()
}
