package php

import (
	"math"
	"sort"
	"strings"
)

// Parameter coercion happens in the callee's file at its declaration line.
// Keep the deprecation separate from value conversion so failed type checks
// cannot emit a warning for a conversion that never happened.
func (e *evaluator) warnLossyIntTypeCoercion(v value, name string, strict bool, line int) error {
	if name != "int" || strict {
		return nil
	}
	var n float64
	var source string
	switch v.kind {
	case floatKind:
		n, source = v.f, "float "+floatText(v.f, -1)
	case stringKind:
		parsed, trailing, ok := v.numeric()
		if !ok || trailing || parsed.kind != floatKind {
			return nil
		}
		n, source = parsed.f, "float-string \""+v.s+"\""
	default:
		return nil
	}
	if math.IsNaN(n) || math.IsInf(n, 0) || math.Trunc(n) == n {
		return nil
	}
	return e.diagnostic("Deprecated", line, "Implicit conversion from "+source+" to int loses precision")
}

func declaredTypeText(name string, nullable bool) string {
	if !strings.ContainsAny(name, "|&") {
		if nullable {
			return "?" + name
		}
		return name
	}
	typeInfo := parseDeclaredType(name, nullable)
	// Zend keeps class names in declaration order, but stores builtin types in
	// a bitmask. zend_type_to_string_resolved() prints that mask in a fixed
	// order, independently of the spelling of the union in the source.
	ordered := make([]string, 0, len(typeInfo.alternatives)+1)
	seen := make(map[string]bool, len(typeInfo.alternatives)+1)
	builtinOrder := []string{"static", "callable", "object", "array", "string", "int", "float", "bool", "false", "true", "void", "never"}
	for _, alternative := range typeInfo.alternatives {
		if len(alternative.types) > 1 {
			text := strings.Join(alternative.types, "&")
			if len(typeInfo.alternatives) > 1 || nullable {
				text = "(" + text + ")"
			}
			ordered = append(ordered, text)
			continue
		}
		term := alternative.types[0]
		lower := strings.ToLower(term)
		isBuiltin := lower == "null"
		for _, builtin := range builtinOrder {
			if lower == builtin {
				isBuiltin = true
				break
			}
		}
		if !isBuiltin {
			ordered = append(ordered, term)
		} else {
			seen[lower] = true
		}
	}
	for _, term := range builtinOrder {
		if seen[term] {
			ordered = append(ordered, term)
		}
	}
	if seen["null"] || nullable {
		if len(ordered) == 1 && !typeInfo.hasIntersection() {
			return "?" + ordered[0]
		}
		ordered = append(ordered, "null")
	}
	return strings.Join(ordered, "|")
}

// Each alternative is one union arm; all types within it must match. The
// spelling is retained separately for diagnostics, while checks use this tree.
type typeAlternative struct {
	types []string
}

type declaredType struct {
	alternatives []typeAlternative
	nullable     bool
}

// Names referenced by a declaration's signatures. Looking them up must not
// predeclare their classes, load another file, or invoke autoload just to
// prove a variance relationship.
func classTypeReferences(decl *classDecl) []string {
	seen := make(map[string]bool)
	var names []string
	add := func(typeName string) {
		for _, arm := range parseDeclaredType(typeName, false).alternatives {
			for _, name := range arm.types {
				lower := strings.ToLower(name)
				if name == "" || builtinTypeName(lower) || lower == "self" || lower == "parent" || lower == "static" || lower == "callable" || lower == "iterable" || seen[lower] {
					continue
				}
				seen[lower] = true
				names = append(names, name)
			}
		}
	}
	addMethods := func(methods []*methodDecl) {
		for _, method := range methods {
			add(method.returnType)
			for _, parameter := range method.parameters {
				add(parameter.typeName)
			}
		}
	}
	addMethods(decl.methodOrder)
	if decl.parent != nil {
		addMethods(decl.parent.methodOrder)
	}
	for _, iface := range decl.interfaces {
		addMethods(interfaceContracts(iface))
	}
	return names
}

func methodTypeReferences(method *methodDecl) []string {
	var names []string
	for _, typeName := range append([]string{method.returnType}, parameterTypeNames(method.parameters)...) {
		names = append(names, typeReferences(typeName)...)
	}
	return names
}

func typeReferences(typeName string) []string {
	var names []string
	for _, arm := range parseDeclaredType(typeName, false).alternatives {
		for _, name := range arm.types {
			lower := strings.ToLower(name)
			if name != "" && !builtinTypeName(lower) && lower != "self" && lower != "parent" && lower != "static" && lower != "callable" && lower != "iterable" {
				names = append(names, name)
			}
		}
	}
	return names
}

func parameterTypeNames(parameters []parameter) []string {
	names := make([]string, len(parameters))
	for i, parameter := range parameters {
		names[i] = parameter.typeName
	}
	return names
}

func parseDeclaredType(name string, nullable bool) declaredType {
	if name == "" {
		return declaredType{}
	}
	result := declaredType{nullable: nullable}
	start, depth := 0, 0
	for index, char := range name {
		switch char {
		case '(':
			depth++
		case ')':
			depth--
		case '|':
			if depth == 0 {
				result.addAlternative(name[start:index])
				start = index + 1
			}
		}
	}
	result.addAlternative(name[start:])
	return result
}

func (t *declaredType) addAlternative(text string) {
	text = strings.TrimPrefix(strings.TrimSuffix(text, ")"), "(")
	t.alternatives = append(t.alternatives, typeAlternative{types: strings.Split(text, "&")})
}

func builtinTypeName(name string) bool {
	switch name {
	case "int", "float", "string", "bool", "array", "object", "callable", "iterable", "mixed", "null", "true", "false", "void", "never":
		return true
	}
	return false
}

func (t declaredType) has(name string) bool {
	for _, alternative := range t.alternatives {
		if len(alternative.types) == 1 && strings.EqualFold(alternative.types[0], name) {
			return true
		}
	}
	return false
}

func (t declaredType) members() []string {
	members := make([]string, 0, len(t.alternatives)+1)
	for _, alternative := range t.alternatives {
		members = append(members, alternative.key())
	}
	if t.nullable && !t.has("null") {
		members = append(members, "null")
	}
	return members
}

func (t declaredType) hasIntersection() bool {
	for _, alternative := range t.alternatives {
		if len(alternative.types) > 1 {
			return true
		}
	}
	return false
}

func (a typeAlternative) key() string {
	parts := make([]string, len(a.types))
	for i, part := range a.types {
		parts[i] = strings.ToLower(part)
	}
	sort.Strings(parts)
	return strings.Join(parts, "&")
}

// A conjunction with more named requirements accepts a strict subset of the
// values accepted by another conjunction. Compare only resolved names: PHP's
// declaration check does not load classes to infer inheritance here.
func (a typeAlternative) moreRestrictiveThan(other typeAlternative) bool {
	if len(a.types) <= len(other.types) {
		return false
	}
	provided := make(map[string]bool, len(a.types))
	for _, name := range a.types {
		provided[strings.ToLower(name)] = true
	}
	for _, name := range other.types {
		if !provided[strings.ToLower(name)] {
			return false
		}
	}
	return true
}

// Return both the converted value and the selected member so precision-loss
// diagnostics refer to the conversion actually performed.
func (e *evaluator) coerceType(v value, name string, nullable, strict bool) (value, string, bool) {
	name = relativeTypeName(name, e.scope, e.lateClass)
	t := parseDeclaredType(name, nullable)
	if len(t.alternatives) == 0 {
		return v, "", true
	}
	if v.kind == nullKind && (t.nullable || t.has("null")) {
		return v, "null", true
	}
	// An existing member wins over every scalar conversion, regardless of the
	// union's spelling. An integer is widened to float only when int is absent.
	for _, alternative := range t.alternatives {
		matches := true
		for _, term := range alternative.types {
			if !e.exactTypeMatch(v, term) {
				matches = false
				break
			}
		}
		if matches {
			return v, strings.Join(alternative.types, "&"), true
		}
	}
	if v.kind == intKind && t.has("float") {
		return floating(float64(v.i)), "float", true
	}
	if strict {
		return value{}, "", false
	}
	// PHP's weak union preference is int, float, string, bool, with numeric
	// strings choosing float when their numeric form is floating point.
	order := []string{"int", "float", "string", "bool"}
	if v.kind == stringKind && t.has("int") && t.has("float") {
		if n, trailing, ok := v.numeric(); ok && !trailing && n.kind == floatKind {
			order[0], order[1] = order[1], order[0]
		}
	}
	for _, term := range order {
		if !t.has(term) {
			continue
		}
		if converted, ok := e.coerceDeclaredType(v, term, false, false); ok {
			return converted, term, true
		}
	}
	return value{}, "", false
}

// Resolve relative type names from the declaration, never from the object on
// which an inherited method happens to run. Only static uses the called class.
func relativeTypeName(name string, owner, called *classDecl) string {
	if name == "" {
		return name
	}
	t := parseDeclaredType(name, false)
	for i := range t.alternatives {
		for j, term := range t.alternatives[i].types {
			switch strings.ToLower(term) {
			case "self":
				if owner != nil {
					t.alternatives[i].types[j] = owner.name
				}
			case "parent":
				if owner != nil && owner.parent != nil {
					t.alternatives[i].types[j] = owner.parent.name
				}
			case "static":
				if called != nil {
					t.alternatives[i].types[j] = called.name
				}
			}
		}
	}
	var arms []string
	for _, arm := range t.alternatives {
		text := strings.Join(arm.types, "&")
		if len(arm.types) > 1 && len(t.alternatives) > 1 {
			text = "(" + text + ")"
		}
		arms = append(arms, text)
	}
	return strings.Join(arms, "|")
}

func (e *evaluator) exactTypeMatch(v value, term string) bool {
	switch term {
	case "null":
		return v.kind == nullKind
	case "false":
		return v.kind == boolKind && v.i == 0
	case "true":
		return v.kind == boolKind && v.i != 0
	case "int":
		return v.kind == intKind
	case "float":
		return v.kind == floatKind
	case "string":
		return v.kind == stringKind
	case "bool":
		return v.kind == boolKind
	case "array":
		return v.kind == arrayKind
	case "iterable":
		return v.kind == arrayKind || v.kind == objectKind && classImplements(v.o.class, "Traversable")
	case "object":
		return v.kind == objectKind
	case "callable":
		_, ok := e.resolveCallable(v, 0)
		return ok
	case "mixed":
		return true
	}
	_, ok := e.coerceDeclaredType(v, term, false, true)
	return ok
}

func literalFitsDeclaredType(v value, name string, nullable bool) bool {
	t := parseDeclaredType(name, nullable)
	if v.kind == nullKind {
		return t.nullable || t.has("null") || t.has("mixed")
	}
	for _, alternative := range t.alternatives {
		if len(alternative.types) != 1 {
			continue
		}
		term := alternative.types[0]
		switch term {
		case "int":
			if v.kind == intKind {
				return true
			}
		case "float":
			if v.kind == floatKind || v.kind == intKind {
				return true
			}
		case "string":
			if v.kind == stringKind {
				return true
			}
		case "bool":
			if v.kind == boolKind {
				return true
			}
		case "true":
			if v.kind == boolKind && v.i != 0 {
				return true
			}
		case "false":
			if v.kind == boolKind && v.i == 0 {
				return true
			}
		case "array":
			if v.kind == arrayKind {
				return true
			}
		case "mixed":
			return true
		}
	}
	return false
}

func (e *evaluator) compatibleMethodDeclaration(implementation, contract *methodDecl) bool {
	if !compatibleUntypedMethod(implementation, contract) {
		return false
	}
	for i, required := range contract.parameters {
		providedIndex := i
		if providedIndex >= len(implementation.parameters) {
			providedIndex = len(implementation.parameters) - 1
		}
		provided := implementation.parameters[providedIndex]
		if !e.typeAcceptsScoped(provided.typeName, provided.nullable, implementation.owner, required.typeName, required.nullable, contract.owner) {
			return false
		}
	}
	return e.typeAcceptsScoped(contract.returnType, contract.returnNullable, contract.owner, implementation.returnType, implementation.returnNullable, implementation.owner)
}

// The receiver must accept every value allowed by the source declaration.
// A declaration is a union of alternatives, each of which is an intersection.
// This is a declaration-time relation: weak scalar coercion is deliberately
// excluded, and no class is loaded just to prove a subtype relationship.
func (e *evaluator) typeAccepts(receiver string, receiverNullable bool, source string, sourceNullable bool) bool {
	return e.typeAcceptsScoped(receiver, receiverNullable, nil, source, sourceNullable, nil)
}

func (e *evaluator) typeAcceptsScoped(receiver string, receiverNullable bool, receiverOwner *classDecl, source string, sourceNullable bool, sourceOwner *classDecl) bool {
	r, s := parseDeclaredType(receiver, receiverNullable), parseDeclaredType(source, sourceNullable)
	// never is the bottom return type. void denotes an absent return value,
	// not the null value produced when a void call is used as an expression.
	if source == "never" {
		return true
	}
	if receiver == "" {
		return true
	}
	if receiver == "void" || source == "void" {
		return receiver == source
	}
	if receiver == "mixed" {
		return true
	}
	if source == "" || source == "mixed" {
		return false
	}
	if s.nullable && !r.nullable && !r.has("null") {
		return false
	}
	for _, sourceArm := range s.alternatives {
		if len(sourceArm.types) == 1 && sourceArm.types[0] == "null" && r.nullable {
			continue
		}
		accepted := false
		for _, receiverArm := range r.alternatives {
			if e.intersectionAcceptsScoped(receiverArm, receiverOwner, sourceArm, sourceOwner) {
				accepted = true
				break
			}
		}
		if !accepted {
			return false
		}
	}
	return true
}

// Every requirement of the receiver must follow from at least one source
// requirement. Thus A&B&C is a subtype of A&B, and a concrete class that
// implements both interfaces is also a subtype of A&B.
func (e *evaluator) intersectionAccepts(receiver, source typeAlternative) bool {
	return e.intersectionAcceptsScoped(receiver, nil, source, nil)
}

func (e *evaluator) intersectionAcceptsScoped(receiver typeAlternative, receiverOwner *classDecl, source typeAlternative, sourceOwner *classDecl) bool {
	for _, want := range receiver.types {
		found := false
		for _, have := range source.types {
			if e.typeTermAcceptsScoped(want, receiverOwner, have, sourceOwner) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (e *evaluator) typeTermAccepts(receiver, source string) bool {
	return e.typeTermAcceptsScoped(receiver, nil, source, nil)
}

func (e *evaluator) typeTermAcceptsScoped(receiver string, receiverOwner *classDecl, source string, sourceOwner *classDecl) bool {
	// static is a return-only type. It remains late-bound in descendants, so
	// replacing it with self or a concrete class is normally not covariant.
	if strings.EqualFold(receiver, "static") {
		if strings.EqualFold(source, "static") {
			return true
		}
		// PHP 8.5 treats self and the class name as static in a final
		// overriding class: no further descendant can change called scope.
		return e.version == PHP85 && sourceOwner != nil && sourceOwner.final &&
			(strings.EqualFold(source, "self") || strings.EqualFold(source, sourceOwner.name))
	}
	if strings.EqualFold(source, "static") {
		if sourceOwner != nil {
			source = sourceOwner.name
		}
	}
	receiver = relativeTypeName(receiver, receiverOwner, nil)
	source = relativeTypeName(source, sourceOwner, nil)
	if strings.EqualFold(receiver, source) || receiver == "mixed" {
		return true
	}
	if receiver == "bool" && (source == "true" || source == "false") {
		return true
	}
	if receiver == "object" && !builtinTypeName(source) {
		return true
	}
	if receiver == "callable" && strings.EqualFold(source, "Closure") {
		return true
	}
	if receiver == "iterable" {
		if source == "array" || strings.EqualFold(source, "Traversable") {
			return true
		}
		if !builtinTypeName(source) {
			if class := e.classes[strings.ToLower(source)]; class != nil && classImplements(class, "Traversable") {
				return true
			}
		}
	}
	if builtinTypeName(receiver) || builtinTypeName(source) {
		return false
	}
	derived := e.classes[strings.ToLower(source)]
	if derived == nil && sourceOwner != nil && strings.EqualFold(source, sourceOwner.name) {
		derived = sourceOwner
	}
	for class := derived; class != nil; class = class.parent {
		if strings.EqualFold(class.name, receiver) || classImplements(class, receiver) {
			return true
		}
	}
	return false
}

func (e *evaluator) methodCompatibilityDiagnostic(implementation, contract *methodDecl, line int) error {
	// PHP may consult registered autoloaders when a declaration cannot be
	// compared without the hierarchy of a named return/parameter class.
	if len(e.autoloaders) != 0 {
		for _, method := range []*methodDecl{implementation, contract} {
			types := []string{method.returnType}
			for _, param := range method.parameters {
				types = append(types, param.typeName)
			}
			for _, declared := range types {
				for _, name := range typeReferences(declared) {
					if e.classes[strings.ToLower(name)] == nil {
						if err := e.tryAutoload(name, line); err != nil {
							return err
						}
					}
				}
			}
		}
		if e.compatibleMethodDeclaration(implementation, contract) {
			return nil
		}
	}
	signatures := "between " + methodSignature(implementation) + " and " + methodSignature(contract)
	// A missing type is relevant only to a type relation that failed. For
	// example, an unrelated unresolved parameter must not obscure a proven
	// arity or return-type incompatibility.
	if compatibleUntypedMethod(implementation, contract) {
		check := func(receiver string, receiverNullable bool, receiverOwner *classDecl, source string, sourceNullable bool, sourceOwner *classDecl) string {
			if e.typeAcceptsScoped(receiver, receiverNullable, receiverOwner, source, sourceNullable, sourceOwner) {
				return ""
			}
			for _, declaration := range []struct {
				name  string
				owner *classDecl
			}{{source, sourceOwner}, {receiver, receiverOwner}} {
				for _, name := range typeReferences(declaration.name) {
					if declaration.owner != nil && strings.EqualFold(name, declaration.owner.name) {
						continue
					}
					if e.classes[strings.ToLower(name)] == nil {
						return name
					}
				}
			}
			return ""
		}
		for i, required := range contract.parameters {
			providedIndex := i
			if providedIndex >= len(implementation.parameters) {
				providedIndex = len(implementation.parameters) - 1
			}
			provided := implementation.parameters[providedIndex]
			if name := check(provided.typeName, provided.nullable, implementation.owner, required.typeName, required.nullable, contract.owner); name != "" {
				return &fatalDiagnostic{line, "Could not check compatibility " + signatures + ", because class " + name + " is not available"}
			}
		}
		if name := check(contract.returnType, contract.returnNullable, contract.owner, implementation.returnType, implementation.returnNullable, implementation.owner); name != "" {
			return &fatalDiagnostic{line, "Could not check compatibility " + signatures + ", because class " + name + " is not available"}
		}
	}
	return &fatalDiagnostic{line, "Declaration of " + methodSignature(implementation) + " must be compatible with " + methodSignature(contract)}
}

// coerceDeclaredType implements the first, simple type-declaration unit.
// The strictness argument belongs to the *calling* file for parameters and
// to the declaring file for return values. It never reads the callee's file
// setting implicitly, because includes and cross-file calls differ here.
func (e *evaluator) coerceDeclaredType(v value, name string, nullable, strict bool) (value, bool) {
	name = relativeTypeName(name, e.scope, e.lateClass)
	if name == "mixed" || nullable && v.kind == nullKind {
		return v, true
	}
	switch name {
	case "int":
		if v.kind == intKind {
			return v, true
		}
		if strict {
			return value{}, false
		}
		switch v.kind {
		case boolKind:
			return integer(v.i), true
		case floatKind:
			if math.IsNaN(v.f) || math.IsInf(v.f, 0) {
				return value{}, false
			}
			return integer(floatToInt(v.f)), true
		case stringKind:
			n, trailing, ok := v.numeric()
			if !ok || trailing {
				return value{}, false
			}
			if n.kind == intKind {
				return n, true
			}
			if math.IsNaN(n.f) || math.IsInf(n.f, 0) {
				return value{}, false
			}
			return integer(floatToInt(n.f)), true
		}
	case "float":
		if v.kind == floatKind {
			return v, true
		}
		if v.kind == intKind {
			return floating(float64(v.i)), true
		}
		if strict {
			return value{}, false
		}
		switch v.kind {
		case boolKind:
			return floating(float64(v.i)), true
		case stringKind:
			n, trailing, ok := v.numeric()
			if ok && !trailing {
				return floating(n.asFloat()), true
			}
		}
	case "string":
		if v.kind == stringKind {
			return v, true
		}
		if !strict && (v.kind == intKind || v.kind == floatKind || v.kind == boolKind) {
			return value{kind: stringKind, s: v.text()}, true
		}
	case "bool":
		if v.kind == boolKind {
			return v, true
		}
		if !strict && (v.kind == intKind || v.kind == floatKind || v.kind == stringKind) {
			return boolean(v.truthy()), true
		}
	case "array":
		if v.kind == arrayKind {
			return v, true
		}
	case "iterable":
		if v.kind == arrayKind || v.kind == objectKind && classImplements(v.o.class, "Traversable") {
			return v, true
		}
	case "object":
		if v.kind == objectKind {
			return v, true
		}
	case "callable":
		if _, ok := e.resolveCallable(v, 0); ok {
			return v, true
		}
	default:
		if v.kind != objectKind {
			return value{}, false
		}
		for class := v.o.class; class != nil; class = class.parent {
			if strings.EqualFold(class.name, name) || classImplements(class, name) {
				return v, true
			}
		}
	}
	return value{}, false
}
