package php

import "unicode"

func validExtractName(name string) bool {
	if name == "" || name == "this" || name == "GLOBALS" {
		return false
	}
	for i, r := range name {
		if r == '_' || unicode.IsLetter(r) || i > 0 && unicode.IsDigit(r) {
			continue
		}
		return false
	}
	return true
}

func (e *evaluator) callExtract(node functionCall) (value, error) {
	args, err := e.evalArgumentsCapturingArray(node.args)
	if err != nil {
		return value{}, err
	}
	bound, present, _, err := bindCallArguments(args, []parameter{{name: "array"}, {name: "flags"}, {name: "prefix"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !present[0] || len(args) > 3 {
		return value{}, &runtimeError{node.line, "ArgumentCountError", "extract() expects 1 to 3 arguments"}
	}
	if bound[0].kind != arrayKind {
		return value{}, builtinTypeError(node.line, "extract", 1, "array", "array", bound[0])
	}
	flags := int64(0)
	if present[1] {
		v, _, ok := e.coerceType(bound[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, "extract", 2, "flags", "int", bound[1])
		}
		flags = v.i
	}
	byReference := flags&256 != 0
	flags &^= 256
	if flags != 0 && flags != 1 {
		return value{}, &runtimeError{node.line, "Error", "gophp: extract() flag is not implemented"}
	}
	var arrayPlace referencePlace
	if byReference {
		argument := argumentForParameter(args, "array", 0)
		if argument == nil || argument.source == nil || !referenceSourceSyntax(argument.source) {
			return value{}, &runtimeError{node.line, "Error", "gophp: extract() with EXTR_REFS requires a referenceable array"}
		}
		if original, ok := argument.source.(variable); ok {
			// extract() may import a key with the same name as its array
			// variable. Keep the original array binding alive while the symbol
			// table name is rebound to that element's reference.
			owner := e.variableOwner(original.token.text)
			stored := owner[original.token.text]
			if stored.kind != refKind {
				stored = value{kind: refKind, ref: &referenceCell{value: stored}}
				owner[original.token.text] = stored
			}
			cell := stored.ref
			arrayPlace = referencePlace{
				read:   func() value { return value{kind: refKind, ref: cell} },
				assign: func(v value) error { return e.writeReference(cell, v, node.line) },
			}
		} else if argument.place != nil {
			arrayPlace = *argument.place
		} else if argument.ref != nil {
			cell := argument.ref
			arrayPlace = referencePlace{
				read:   func() value { return value{kind: refKind, ref: cell} },
				assign: func(v value) error { return e.writeReference(cell, v, node.line) },
			}
		} else {
			arrayPlace, err = e.referencePlace(argument.source, true)
			if err != nil {
				return value{}, err
			}
		}
	}
	count := int64(0)
	for _, entry := range bound[0].a.entries {
		if !entry.key.isString || !validExtractName(entry.key.text) {
			continue
		}
		name := entry.key.text
		if flags == 1 {
			if _, exists := e.readVariable(name); exists {
				continue
			}
		}
		if byReference {
			key := value{kind: stringKind, s: name}
			refSource := offsetExpr{key: literal{value: key, line: node.line}, line: node.line}
			refTarget := variable{token{kind: "variable", text: name, line: node.line}}
			target, err := e.referencePlace(refTarget, false)
			if err != nil {
				return value{}, err
			}
			source, err := e.offsetReferencePlace(arrayPlace, refSource)
			if err != nil {
				return value{}, err
			}
			if _, err := e.bindReferencePlaces(referenceAssignment{target: refTarget, line: node.line}, target, source); err != nil {
				return value{}, err
			}
		} else if err := e.variableSlotAt(name, node.line).set(dereference(entry.value)); err != nil {
			return value{}, err
		}
		count++
	}
	return integer(count), nil
}
