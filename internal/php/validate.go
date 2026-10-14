package php

// Like PHP's compiler, validate read/write context before executing any branch.
// In particular an append dimension may only occur in an assignment target.
func validateStatements(statements []statement) *runtimeError { return validateBlock(statements, 0) }

func validateBlock(statements []statement, depth int) *runtimeError {
	for _, s := range statements {
		if s.constant != nil {
			if !constantInitializer(s.constant.initial) {
				return &runtimeError{s.constant.line, "", "Constant expression contains invalid operations"}
			}
			if err := validateExpression(s.constant.initial, false); err != nil {
				return err
			}
		}
		if s.kind == "switch" {
			defaults := 0
			for _, c := range s.cases {
				if c.condition == nil {
					defaults++
				}
				if defaults > 1 {
					return &runtimeError{c.line, "", "Switch statements may only contain one default clause"}
				}
				if err := validateExpression(c.condition, false); err != nil {
					return err
				}
				if err := validateBlock(c.body, depth+1); err != nil {
					return err
				}
			}
		}
		if s.kind == "unset" {
			for _, x := range s.expressions {
				if err := validateUnset(x); err != nil {
					return err
				}
			}
			continue
		}
		if s.kind == "break" || s.kind == "continue" {
			if _, err := loopLevel(s, depth); err != nil {
				return err
			}
		}
		if s.function != nil {
			if err := validateCallable(s.function); err != nil {
				return err
			}
		}
		if s.class != nil {
			for _, constant := range s.class.constantDecls {
				if !constantInitializer(constant.initial) {
					return &runtimeError{constant.line, "", "Constant expression contains invalid operations"}
				}
				if err := validateExpression(constant.initial, false); err != nil {
					return err
				}
			}
			for _, property := range s.class.properties {
				if !constantInitializer(property.initial) {
					return &runtimeError{property.line, "", "Constant expression contains invalid operations"}
				}
				if err := validateExpression(property.initial, false); err != nil {
					return err
				}
			}
			for _, method := range s.class.methodOrder {
				if err := validateCallable(method); err != nil {
					return err
				}
			}
		}
		for _, x := range append(append(append([]expr{}, s.initial...), s.expressions...), s.update...) {
			if err := validateExpression(x, false); err != nil {
				return err
			}
		}
		if err := validateExpression(s.condition, false); err != nil {
			return err
		}
		for _, target := range []expr{s.keyTarget, s.valueTarget} {
			if target != nil {
				if err := validateExpression(assignment{target, nil, expressionLine(target)}, false); err != nil {
					return err
				}
			}
		}
		childDepth := depth
		if s.kind == "while" || s.kind == "for" || s.kind == "do" || s.kind == "foreach" {
			childDepth++
		}
		if err := validateBlock(s.thenBranch, childDepth); err != nil {
			return err
		}
		if err := validateBlock(s.elseBranch, depth); err != nil {
			return err
		}
		for _, clause := range s.catches {
			if err := validateBlock(clause.body, depth); err != nil {
				return err
			}
		}
		if err := validateBlock(s.finally, depth); err != nil {
			return err
		}
	}
	return nil
}

func validateExpression(x expr, write bool) *runtimeError {
	if write && hasNullsafe(x) {
		return &runtimeError{expressionLine(x), "", "Can't use nullsafe operator in write context"}
	}
	switch node := x.(type) {
	case unpackArgument:
		return validateExpression(node.operand, false)
	case namedArgument:
		return validateExpression(node.operand, false)
	case throwExpr:
		return validateExpression(node.operand, false)
	case dynamicCall:
		if err := validateExpression(node.callee, false); err != nil {
			return err
		}
		for _, arg := range node.args {
			if err := validateExpression(arg, false); err != nil {
				return err
			}
		}
	case closureExpr:
		return validateCallable(node.method)
	case instanceofExpr:
		if err := validateExpression(node.operand, false); err != nil {
			return err
		}
		return validateExpression(node.classExpr, false)
	case includeExpr:
		return validateExpression(node.path, false)
	case interpolatedString:
		for _, part := range node.parts {
			if err := validateExpression(part, false); err != nil {
				return err
			}
		}
	case ternaryExpr:
		for _, branch := range []expr{node.condition, node.yes, node.no} {
			if err := validateExpression(branch, false); err != nil {
				return err
			}
		}
	case coalesceAssignment:
		// ??= reads before writing; append dimensions are never valid here.
		if err := validateExpression(node.target, false); err != nil {
			return err
		}
		return validateExpression(assignment{node.target, node.source, node.line}, false)
	case nullsafeExpr:
		if err := validateExpression(node.base, false); err != nil {
			return err
		}
		for _, arg := range node.args {
			if err := validateExpression(arg, false); err != nil {
				return err
			}
		}
	case matchExpr:
		if err := validateExpression(node.subject, false); err != nil {
			return err
		}
		defaults := 0
		for _, arm := range node.arms {
			if len(arm.conditions) == 0 {
				defaults++
			}
			if defaults > 1 {
				return &runtimeError{arm.line, "", "Match expressions may only contain one default arm"}
			}
			for _, condition := range arm.conditions {
				if err := validateExpression(condition, false); err != nil {
					return err
				}
			}
			if err := validateExpression(arm.result, false); err != nil {
				return err
			}
		}
	case castExpr:
		return validateExpression(node.operand, false)
	case languageExpr:
		for _, arg := range node.args {
			if node.name == "isset" && !issetTarget(arg) {
				return &runtimeError{node.line, "", "Cannot use isset() on the result of an expression (you can use \"null !== expression\" instead)"}
			}
			if err := validateExpression(arg, false); err != nil {
				return err
			}
		}
	case compoundAssignment:
		return validateExpression(assignment{node.target, node.source, node.op.line}, false)
	case increment:
		return validateExpression(assignment{node.target, nil, node.op.line}, false)
	case functionCall:
		for _, arg := range node.args {
			if err := validateExpression(arg, false); err != nil {
				return err
			}
		}
	case assignment:
		if v, ok := node.target.(variable); ok && v.token.text == "this" {
			return &runtimeError{node.line, "", "Cannot re-assign $this"}
		}
		if err := validateExpression(node.target, true); err != nil {
			return err
		}
		return validateExpression(node.source, false)
	case offsetExpr:
		if node.key == nil && !write {
			return &runtimeError{node.line, "", "Cannot use [] for reading"}
		}
		if err := validateExpression(node.base, write); err != nil {
			return err
		}
		return validateExpression(node.key, false)
	case propertyExpr:
		return validateExpression(node.base, write)
	case methodCall:
		if err := validateExpression(node.receiver, false); err != nil {
			return err
		}
		for _, arg := range node.args {
			if err := validateExpression(arg, false); err != nil {
				return err
			}
		}
	case newObject:
		if err := validateExpression(node.classExpr, false); err != nil {
			return err
		}
		for _, arg := range node.args {
			if err := validateExpression(arg, false); err != nil {
				return err
			}
		}
	case binary:
		if err := validateExpression(node.left, false); err != nil {
			return err
		}
		return validateExpression(node.right, false)
	case unary:
		return validateExpression(node.operand, false)
	case arrayLiteral:
		// A fully literal array with an array key is rejected by PHP's
		// constant-array compiler, even if the containing branch is skipped.
		if literalArray(node) {
			for _, entry := range node.elements {
				if _, invalid := entry.key.(arrayLiteral); invalid {
					return &runtimeError{node.line, "", "Illegal offset type"}
				}
			}
		}
		for _, entry := range node.elements {
			if err := validateExpression(entry.key, false); err != nil {
				return err
			}
			if err := validateExpression(entry.value, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func constantInitializer(x expr) bool {
	switch node := x.(type) {
	case nil, literal, constantExpr:
		return true
	case classConstantExpr:
		return true
	case unary:
		return constantInitializer(node.operand)
	case binary:
		return constantInitializer(node.left) && constantInitializer(node.right)
	case ternaryExpr:
		return constantInitializer(node.condition) && constantInitializer(node.yes) && constantInitializer(node.no)
	case arrayLiteral:
		for _, entry := range node.elements {
			if !constantInitializer(entry.key) || !constantInitializer(entry.value) {
				return false
			}
		}
		return true
	}
	return false
}

func literalArray(node arrayLiteral) bool {
	for _, entry := range node.elements {
		for _, x := range []expr{entry.key, entry.value} {
			switch child := x.(type) {
			case nil, literal:
			case arrayLiteral:
				if !literalArray(child) {
					return false
				}
			default:
				return false
			}
		}
	}
	return true
}

func validateCallable(decl *methodDecl) *runtimeError {
	for _, param := range decl.parameters {
		if !constantInitializer(param.initial) {
			return &runtimeError{param.line, "", "Constant expression contains invalid operations"}
		}
		if err := validateExpression(param.initial, false); err != nil {
			return err
		}
	}
	return validateStatements(decl.body)
}
