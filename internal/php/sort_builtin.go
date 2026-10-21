package php

func (e *evaluator) callSort(node functionCall) (value, error) { return e.callBuiltinSort(node, false) }
