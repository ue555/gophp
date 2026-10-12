package php

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestClassScopeAndCallState(t *testing.T) {
	source := `<?php
class Box {
    private $value = 1;
    function read($other) { return $other->value; }
    function fail() { return 1 / 0; }
}
$o = new Box; $same = $o;
echo $o->read($same);
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		// Every Run starts with fresh classes, including repeated declarations.
		for i := 0; i < 2; i++ {
			var out, err bytes.Buffer
			if code := RunWithVersion(source, "scope.php", &out, &err, version); code != 0 || out.String() != "1" || err.Len() != 0 {
				t.Fatalf("PHP %s: exit=%d stdout=%q stderr=%q", version, code, &out, &err)
			}
		}
	}
	tokens, err := lex(source)
	if err != nil {
		t.Fatal(err)
	}
	program, err := parse(tokens)
	if err != nil {
		t.Fatal(err)
	}
	e := evaluator{version: PHP84, stdout: io.Discard, stderr: io.Discard, variables: make(map[string]value)}
	if err := e.executeProgram(program); err != nil {
		t.Fatal(err)
	}
	object := e.variables["o"]
	if _, err := e.invoke(object, object.o.class.methods["fail"], nil, 1); err == nil {
		t.Fatal("expected method failure")
	}
	if e.scope != nil || len(e.frames) != 0 || e.variables["o"] != object || len(e.failureFrames) != 1 {
		t.Fatal("method failure did not restore caller state")
	}
}

func TestUnsupportedClassFeatures(t *testing.T) {
	for _, source := range []string{
		"class C extends P {}", "class C implements I {}", "trait T {}", "interface I {}",
		"class C { public int $x; }", "class C { static function f() {} }",
		"class C { function f(&$x) {} }", "class C { function f(...$args) {} }",
		"class C { function f(): int { return 1; } }", "class C { function __get($x) {} }",
		"class C { function __destruct() {} }", "class C { function f() {} } $o = new C; $o->f(x: 1);",
	} {
		var out, err bytes.Buffer
		if code := Run("<?php echo 'must not execute'; "+source, "class.php", &out, &err); code != 255 || out.Len() != 0 || err.Len() == 0 {
			t.Fatalf("unsupported feature: %s; code=%d out=%q err=%q", source, code, &out, &err)
		}
	}
}

func TestMethodRecursionLimit(t *testing.T) {
	var out, err bytes.Buffer
	code := Run("<?php class C { function f() { return $this->f(); } } $o = new C; $o->f();", "recursive.php", &out, &err)
	if code != 255 || !strings.Contains(err.String(), "maximum method call depth exceeded") {
		t.Fatalf("exit=%d stderr=%q", code, &err)
	}
}

func TestClassCompileErrorsBeforeExecution(t *testing.T) {
	for _, source := range []string{
		"$this = 1;", "class C { public $x = $undefined; }",
		"class C { function f($a = $bad) {} }", "class C { function f() { $this = null; } }",
		"class C {} class c {}",
	} {
		var out, err bytes.Buffer
		if code := Run("<?php echo 'must not execute'; "+source, "class.php", &out, &err); code != 255 || out.Len() != 0 {
			t.Fatalf("source=%s exit=%d out=%q err=%q", source, code, &out, &err)
		}
	}
}
