package php

import (
	"math"
	"strings"
)

func declaredTypeText(name string, nullable bool) string {
	if nullable {
		return "?" + name
	}
	return name
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
		if !e.typeAccepts(provided.typeName, provided.nullable, required.typeName, required.nullable) {
			return false
		}
	}
	return e.typeAccepts(contract.returnType, contract.returnNullable, implementation.returnType, implementation.returnNullable)
}

// The receiver must accept every value allowed by the source declaration.
// This covers simple covariance/contravariance without treating typed and
// untyped methods as interchangeable. Compound declarations remain rejected
// by the parser until their variance rules are implemented.
func (e *evaluator) typeAccepts(receiver string, receiverNullable bool, source string, sourceNullable bool) bool {
	if receiver == "" || receiver == "mixed" {
		return true
	}
	if source == "" || source == "mixed" || sourceNullable && !receiverNullable {
		return false
	}
	if strings.EqualFold(receiver, source) {
		return true
	}
	if receiver == "float" && source == "int" {
		return true
	}
	derived := e.classes[strings.ToLower(source)]
	if derived == nil {
		return false
	}
	for class := derived; class != nil; class = class.parent {
		if strings.EqualFold(class.name, receiver) || classImplements(class, receiver) {
			return true
		}
	}
	return false
}

// coerceDeclaredType implements the first, simple type-declaration unit.
// The strictness argument belongs to the *calling* file for parameters and
// to the declaring file for return values. It never reads the callee's file
// setting implicitly, because includes and cross-file calls differ here.
func (e *evaluator) coerceDeclaredType(v value, name string, nullable, strict bool) (value, bool) {
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
	case "object":
		if v.kind == objectKind {
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
