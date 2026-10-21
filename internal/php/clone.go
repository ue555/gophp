package php

func (e *evaluator) evalClone(node cloneExpr) (value, error) {
	original, err := e.eval(node.operand)
	if err != nil {
		return value{}, err
	}
	if original.kind != objectKind {
		return value{}, &runtimeError{node.line, "Error", "__clone method called on non-object"}
	}
	if original.o.enumCase != "" {
		return value{}, &runtimeError{node.line, "Error", "Trying to clone an uncloneable object of class " + original.o.className()}
	}
	copyObject := *original.o
	if original.o.properties != nil {
		properties := *original.o.properties
		properties.entries = append([]arrayEntry(nil), original.o.properties.entries...)
		copyObject.properties = &properties
	}
	copyValue := value{kind: objectKind, o: &copyObject}
	if method := e.lookupMethod(copyObject.class, "__clone"); method != nil {
		if _, err = e.invoke(copyValue, method, nil, node.line); err != nil {
			return value{}, err
		}
	}
	return copyValue, nil
}
