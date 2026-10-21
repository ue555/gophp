package php

func (e *evaluator) dynamicVariableName(node dynamicVariable) (string, error) {
	v, err := e.eval(node.nameExpr)
	if err != nil {
		return "", err
	}
	return e.toText(v, node.line)
}

func (e *evaluator) evalDynamicVariable(node dynamicVariable) (value, error) {
	name, err := e.dynamicVariableName(node)
	if err != nil {
		return value{}, err
	}
	v, exists := e.readVariable(name)
	if !exists {
		if err := e.warn(node.line, "Undefined variable $"+name); err != nil {
			return value{}, err
		}
	}
	return v, nil
}
