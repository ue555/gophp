package php

import (
	"fmt"
	"math"
	"strings"
)

type propertyDecl struct {
	name, visibility string
	initial          expr
	line             int
}
type parameter struct {
	name    string
	initial expr
	line    int
}
type methodDecl struct {
	name, visibility string
	parameters       []parameter
	body             []statement
	line             int
}
type classDecl struct {
	name        string
	properties  []propertyDecl
	methods     map[string]*methodDecl
	methodOrder []*methodDecl
	line        int
}

func (p *parser) classDeclaration() (*classDecl, error) {
	name := p.current()
	if err := p.expect("name"); err != nil {
		return nil, err
	}
	decl := &classDecl{name: name.text, line: name.line, methods: make(map[string]*methodDecl)}
	if err := p.expect("{"); err != nil {
		return nil, err
	}
	for p.current().kind != "}" {
		visibility := "public"
		explicit := false
		for _, modifier := range []string{"public", "protected", "private", "var"} {
			if p.acceptKeyword(modifier) {
				explicit = true
				if modifier != "var" {
					visibility = modifier
				}
				break
			}
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
			method := &methodDecl{name: name.text, visibility: visibility, line: name.line}
			var err error
			method.parameters, method.body, err = p.callableBody()
			if err != nil {
				return nil, err
			}
			decl.methods[lower] = method
			decl.methodOrder = append(decl.methodOrder, method)
			continue
		}
		if !explicit || p.current().kind != "variable" {
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
			property := propertyDecl{name: name.text, line: name.line, visibility: visibility}
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

type returned struct{ value value }

func (*returned) Error() string { return "return" }

type callFrame struct {
	class, method string
	line          int
	args          []value
	parameters    []parameter
	locals        map[string]value
}

func (e *evaluator) initClasses() {
	if e.classes == nil {
		e.classes = map[string]*classDecl{"stdclass": {name: "stdClass"}}
	}
}

func (e *evaluator) declareClass(decl *classDecl) error {
	e.initClasses()
	key := strings.ToLower(decl.name)
	if _, exists := e.classes[key]; exists {
		return &runtimeError{decl.line, "", "Cannot declare class " + decl.name + ", because the name is already in use"}
	}
	e.classes[key] = decl
	return nil
}

func (e *evaluator) newInstance(node newObject) (value, error) {
	e.initClasses()
	class := e.classes[strings.ToLower(node.name)]
	if class == nil {
		return value{}, &runtimeError{node.line, "Error", fmt.Sprintf("Class %q not found", node.name)}
	}
	object := value{kind: objectKind, o: &phpObject{class: class, properties: &phpArray{}}}
	for _, property := range class.properties {
		var v value
		var err error
		if property.initial != nil {
			v, err = e.eval(property.initial)
		}
		if err != nil {
			return value{}, err
		}
		object.o.properties = object.o.properties.with(arrayKey{text: property.name, isString: true}, v)
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
	} else {
		// PHP evaluates and ignores positional arguments if no constructor exists.
		if _, err := e.evalArguments(node.args); err != nil {
			return value{}, err
		}
	}
	return object, nil
}

func (e *evaluator) evalArguments(args []expr) ([]value, error) {
	values := make([]value, 0, len(args))
	for _, arg := range args {
		v, err := e.eval(arg)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func (e *evaluator) callMethod(node methodCall) (value, error) {
	object, err := e.eval(node.receiver)
	if err != nil {
		return value{}, err
	}
	if object.kind != objectKind {
		return value{}, &runtimeError{node.line, "Error", "Call to a member function " + node.name + "() on " + object.typeName()}
	}
	var method *methodDecl
	if object.o.class != nil {
		method = object.o.class.methods[strings.ToLower(node.name)]
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
	return e.invoke(object, method, args, node.line)
}

func (e *evaluator) methodAccess(class *classDecl, method *methodDecl, line int) error {
	if method.visibility == "public" || e.scope == class {
		return nil
	}
	from := "global scope"
	if e.scope != nil {
		from = "scope " + e.scope.name
	}
	return &runtimeError{line, "Error", "Call to " + method.visibility + " method " + class.name + "::" + method.name + "() from " + from}
}

func (e *evaluator) propertyAccess(object *phpObject, name string, line int) error {
	if object.class == nil {
		return nil
	}
	for _, property := range object.class.properties {
		if property.name == name && property.visibility != "public" && e.scope != object.class {
			return &runtimeError{line, "Error", "Cannot access " + property.visibility + " property " + object.className() + "::$" + name}
		}
	}
	return nil
}

func (e *evaluator) invoke(object value, method *methodDecl, args []value, line int) (result value, err error) {
	// Bound recursive calls so untrusted PHP cannot overflow the Go stack.
	if len(e.frames) >= 256 {
		kind := "method"
		if object.kind != objectKind {
			kind = "function"
		}
		return value{}, &runtimeError{line, "Error", "gophp: maximum " + kind + " call depth exceeded (256)"}
	}
	outer, scope := e.variables, e.scope
	e.variables = make(map[string]value)
	e.scope = nil
	className, callableName := "", method.name
	if object.kind == objectKind {
		e.variables["this"] = object
		e.scope = object.o.class
		className = object.o.className()
		callableName = className + "::" + method.name
	}
	e.frames = append(e.frames, callFrame{
		class: className, method: method.name, line: line,
		args: args, parameters: method.parameters, locals: e.variables,
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
		e.variables, e.scope = outer, scope
	}()
	required := requiredParameters(method)
	if len(args) < required {
		count := "exactly"
		if required < len(method.parameters) {
			count = "at least"
		}
		return value{}, &runtimeError{method.line, "ArgumentCountError", fmt.Sprintf("Too few arguments to function %s(), %d passed in %s on line %d and %s %d expected", callableName, len(args), e.filename, line, count, required)}
	}
	for i, parameter := range method.parameters {
		v := value{}
		if i < len(args) {
			v = args[i]
		} else {
			v, err = e.eval(parameter.initial)
			if err != nil {
				return value{}, err
			}
		}
		e.variables[parameter.name] = v
	}
	err = e.execute(method.body)
	if r, ok := err.(*returned); ok {
		return r.value, nil
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
			name = frame.class + "->" + name
		}
		fmt.Fprintf(&result, "#%d %s(%d): %s(%s)\n", len(e.failureFrames)-i-1, e.filename, frame.line, name, strings.Join(args, ", "))
	}
	fmt.Fprintf(&result, "#%d {main}\n", len(e.failureFrames))
	return result.String()
}
