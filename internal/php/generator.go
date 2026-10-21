package php

import (
	"fmt"
	"reflect"
	"sync"
)

// A generator owns a suspended invocation. The worker and its caller only
// access interpreter state on opposite sides of the yield handshake.
type phpGenerator struct {
	start             func(*phpGenerator) (value, error)
	resume            chan generatorResume
	events            chan generatorEvent
	stop              chan struct{}
	done              chan struct{}
	once              sync.Once
	started, finished bool
	// rewind is permitted only before or while suspended at the first yield.
	pastFirstYield    bool
	yielded           bool
	current           generatorEvent
	nextKey           int64
	returned          value
	hasReturn         bool
	yieldsByReference bool
	failureFrames     []callFrame
}

type generatorEvent struct {
	key, item value
	err       error
	done      bool
}

type generatorResume struct {
	item value
	err  error
}

type generatorRegistry struct {
	mu  sync.Mutex
	all []*phpGenerator
}

func (r *generatorRegistry) add(g *phpGenerator) {
	r.mu.Lock()
	r.all = append(r.all, g)
	r.mu.Unlock()
}

func (r *generatorRegistry) closeAll() {
	r.mu.Lock()
	all := append([]*phpGenerator(nil), r.all...)
	r.mu.Unlock()
	for _, g := range all {
		g.once.Do(func() { close(g.stop) })
	}
	for _, g := range all {
		if g.started {
			<-g.done
		}
	}
}

func newPHPGenerator(start func(*phpGenerator) (value, error)) *phpGenerator {
	return &phpGenerator{start: start, resume: make(chan generatorResume), events: make(chan generatorEvent), stop: make(chan struct{}), done: make(chan struct{})}
}

func (g *phpGenerator) step() error {
	return g.stepWith(value{kind: nullKind})
}

func (g *phpGenerator) stepWith(input value) error {
	return g.stepSignal(generatorResume{item: input})
}

func (g *phpGenerator) stepSignal(input generatorResume) error {
	if g.finished {
		return nil
	}
	if !g.started {
		g.started = true
		go func() {
			defer close(g.done)
			select {
			case <-g.resume:
			case <-g.stop:
				return
			}
			returned, err := g.start(g)
			g.returned = returned
			g.hasReturn = err == nil
			select {
			case g.events <- generatorEvent{err: err, done: true}:
			case <-g.stop:
			}
		}()
	}
	select {
	case g.resume <- input:
	case <-g.stop:
		return nil
	}
	select {
	case event := <-g.events:
		if !event.done {
			if g.yielded {
				g.pastFirstYield = true
			}
			g.yielded = true
		} else if g.yielded {
			g.pastFirstYield = true
		}
		g.current = event
		g.finished = event.done
		return event.err
	case <-g.stop:
		return nil
	}
}

func (g *phpGenerator) emit(key, item value) (value, error) {
	select {
	case g.events <- generatorEvent{key: key, item: item}:
	case <-g.stop:
		return value{}, &exitSignal{}
	}
	select {
	case received := <-g.resume:
		return received.item, received.err
	case <-g.stop:
		return value{}, &exitSignal{}
	}
}

func (e *evaluator) evalYield(node yieldExpr) (value, error) {
	if e.activeGenerator == nil {
		return value{}, &runtimeError{node.line, "Error", "Cannot use yield outside a function"}
	}
	g := e.activeGenerator
	if node.from {
		if g.yieldsByReference {
			return value{}, &runtimeError{node.line, "Error", "gophp: yield from in reference-yielding generators is not implemented"}
		}
		iterable, err := e.eval(node.item)
		if err != nil {
			return value{}, err
		}
		switch iterable.kind {
		case arrayKind:
			for _, entry := range iterable.a.entries {
				key := integer(entry.key.index)
				if entry.key.isString {
					key = value{kind: stringKind, s: entry.key.text}
				}
				if _, err := g.emit(key, dereference(entry.value)); err != nil {
					return value{}, err
				}
			}
		case objectKind:
			return e.yieldFromTraversable(g, iterable, node.line, 0, true)
		default:
			return value{}, &runtimeError{node.line, "Error", "Can use \"yield from\" only with arrays and Traversables"}
		}
		return value{kind: nullKind}, nil
	}
	key := integer(g.nextKey)
	if node.key != nil {
		var err error
		key, err = e.eval(node.key)
		if err != nil {
			return value{}, err
		}
		if key.kind == intKind && key.i >= g.nextKey {
			g.nextKey = key.i + 1
		}
	} else {
		g.nextKey++
	}
	item := value{kind: nullKind}
	if node.item != nil {
		if g.yieldsByReference {
			if !referenceTargetSyntax(node.item) {
				return value{}, &runtimeError{node.line, "Error", "gophp: yielding an expression by reference is not implemented"}
			}
			e.foreachTempCounter++
			name := fmt.Sprintf("\x00yield_reference_%d", e.foreachTempCounter)
			target := variable{token{kind: "variable", text: name, line: node.line}}
			if _, err := e.evalReferenceAssignment(referenceAssignment{target: target, source: node.item, line: node.line}); err != nil {
				return value{}, err
			}
			item = e.variables[name]
			delete(e.variables, name)
		} else {
			var err error
			item, err = e.eval(node.item)
			if err != nil {
				return value{}, err
			}
		}
	}
	return g.emit(key, item)
}

func (e *evaluator) yieldFromTraversable(g *phpGenerator, iterable value, line, depth int, forwardSend bool) (value, error) {
	// IteratorAggregate may delegate through many finite getIterator() calls.
	// Retain a high safety bound for cyclic chains until their identity can be
	// tracked across all Traversable adapters.
	if depth >= 1024 {
		return value{}, &runtimeError{line, "Error", "gophp: iterator nesting limit exceeded"}
	}
	if iterable.o == nil || !classImplements(iterable.o.class, "Traversable") {
		return value{}, &runtimeError{line, "Error", "Can use \"yield from\" only with arrays and Traversables"}
	}
	if inner := iterable.o.generator; inner != nil {
		// The yield expression belongs to the delegated generator. A value
		// sent (or an exception thrown) into the outer generator resumes it.
		resume := generatorResume{item: value{kind: nullKind}}
		for {
			if err := inner.stepSignal(resume); err != nil {
				return value{}, err
			}
			if inner.finished {
				if forwardSend {
					return inner.returned, nil
				}
				return value{kind: nullKind}, nil
			}
			item, err := g.emit(inner.current.key, dereference(inner.current.item))
			if forwardSend {
				resume = generatorResume{item: item, err: err}
			} else {
				resume = generatorResume{item: value{kind: nullKind}, err: err}
			}
		}
	}
	if classImplements(iterable.o.class, "IteratorAggregate") {
		inner, err := e.iteratorMethod(iterable, "getIterator", line)
		if err != nil {
			return value{}, err
		}
		if inner.kind != objectKind || inner.o == nil || !classImplements(inner.o.class, "Traversable") {
			return value{}, &runtimeError{line, "TypeError", iterable.o.className() + "::getIterator() must return Traversable"}
		}
		return e.yieldFromTraversable(g, inner, line, depth+1, false)
	}
	if !classImplements(iterable.o.class, "Iterator") {
		return value{}, &runtimeError{line, "Error", "Can use \"yield from\" only with arrays and Traversables"}
	}
	if _, err := e.iteratorMethod(iterable, "rewind", line); err != nil {
		return value{}, err
	}
	for {
		valid, err := e.iteratorMethod(iterable, "valid", line)
		if err != nil {
			return value{}, err
		}
		if !valid.truthy() {
			return value{kind: nullKind}, nil
		}
		key, err := e.iteratorMethod(iterable, "key", line)
		if err != nil {
			return value{}, err
		}
		item, err := e.iteratorMethod(iterable, "current", line)
		if err != nil {
			return value{}, err
		}
		if _, err := g.emit(key, item); err != nil {
			return value{}, err
		}
		if _, err := e.iteratorMethod(iterable, "next", line); err != nil {
			return value{}, err
		}
	}
}

func methodContainsYield(body []statement) bool {
	return methodYieldLine(body, false) != 0
}

func methodYieldFromLine(body []statement) int {
	return methodYieldLine(body, true)
}

// Nested callables compile independently of the enclosing generator.
func methodYieldLine(body []statement, fromOnly bool) int {
	target := reflect.TypeOf(yieldExpr{})
	var visit func(reflect.Value) int
	visit = func(v reflect.Value) int {
		if !v.IsValid() {
			return 0
		}
		if v.Type() == target {
			if !fromOnly || v.FieldByName("from").Bool() {
				return int(v.FieldByName("line").Int())
			}
			return 0
		}
		if v.Kind() == reflect.Interface {
			if v.IsNil() {
				return 0
			}
			return visit(v.Elem())
		}
		if v.Kind() == reflect.Pointer {
			return 0
		}
		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				if line := visit(v.Index(i)); line != 0 {
					return line
				}
			}
		case reflect.Struct:
			// A nested function or closure does not turn the enclosing function
			// into a generator.
			if v.Type().Name() == "closureExpr" || v.Type().Name() == "arrowExpr" {
				return 0
			}
			for i := 0; i < v.NumField(); i++ {
				if line := visit(v.Field(i)); line != 0 {
					return line
				}
			}
		}
		return 0
	}
	return visit(reflect.ValueOf(body))
}
