package php

// PHP reverses >/>= to </<=. This matters for arrays with differing key sets:
// they need not form a total order. Spaceship returns 1 for incomparable values.
func (e *evaluator) comparisonValue(op string, a, b value, line int) (value, error) {
	if op == ">" || op == ">=" {
		a, b = b, a
	}
	order, err := e.valueOrder(a, b, line, make(map[*phpArray]bool))
	if err != nil {
		return value{}, err
	}
	switch op {
	case "<=>":
		return integer(int64(order)), nil
	case "==":
		return boolean(order == 0), nil
	case "!=", "<>":
		return boolean(order != 0), nil
	case "<", ">":
		return boolean(order < 0), nil
	case "<=", ">=":
		return boolean(order <= 0), nil
	}
	panic("invalid comparison")
}

func (e *evaluator) valueOrder(a, b value, line int, active map[*phpArray]bool) (int, error) {
	if a.kind == boolKind || b.kind == boolKind || a.kind == nullKind || b.kind == nullKind {
		order, _ := scalarOrder(a, b)
		return order, nil
	}
	if a.kind == objectKind && b.kind == objectKind {
		if a.o == b.o {
			return 0, nil
		}
		if a.o.enumCase != "" || b.o.enumCase != "" {
			// Enum objects are identity-comparable but have no value ordering.
			return 1, nil
		}
		if a.o.className() != b.o.className() {
			return 1, nil
		}
		return e.arrayOrder(a.o.properties, b.o.properties, line, active)
	}
	if a.kind == objectKind || b.kind == objectKind {
		other := b
		object := a
		if b.kind == objectKind {
			other = a
			object = b
		}
		if other.kind == intKind || other.kind == floatKind {
			e.diagnostic("Notice", line, "Object of class "+object.o.className()+" could not be converted to "+other.typeName())
			one := integer(1)
			if other.kind == floatKind {
				one = floating(1)
			}
			if a.kind == objectKind {
				a = one
			} else {
				b = one
			}
			return e.valueOrder(a, b, line, active)
		}
		if a.kind == objectKind {
			return 1, nil
		}
		return -1, nil
	}
	if a.kind == arrayKind && b.kind == arrayKind {
		return e.arrayOrder(a.a, b.a, line, active)
	}
	if a.kind == arrayKind {
		return 1, nil
	}
	if b.kind == arrayKind {
		return -1, nil
	}
	order, ordered := scalarOrder(a, b)
	if !ordered {
		return 1, nil
	}
	return order, nil
}

func (e *evaluator) arrayOrder(a, b *phpArray, line int, active map[*phpArray]bool) (int, error) {
	if a == b {
		return 0, nil
	}
	if active[a] {
		class := "Error"
		if e.version == PHP83 {
			class = ""
		}
		return 0, &runtimeError{line, class, "Nesting level too deep - recursive dependency?"}
	}
	active[a] = true
	defer delete(active, a)
	if len(a.entries) != len(b.entries) {
		return integerOrder(int64(len(a.entries)), int64(len(b.entries))), nil
	}
	for _, entry := range a.entries {
		other, found := b.get(entry.key)
		if !found {
			return 1, nil
		}
		order, err := e.valueOrder(dereference(entry.value), other, line, active)
		if err != nil || order != 0 {
			return order, err
		}
	}
	return 0, nil
}
