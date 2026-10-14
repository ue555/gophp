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
		"class C extends P {}", "class C implements I {}",
		"class C { public int $x; }",
		"class C { function f(&$x) {} }", "class C { function f(...$args, $tail) {} }",
		"class C { function f(): int|float { return 1; } }", "class C { function __get($x) {} }",
		"class C { function __destruct() {} }",
	} {
		var out, err bytes.Buffer
		if code := Run("<?php echo 'must not execute'; "+source, "class.php", &out, &err); code != 255 || out.Len() != 0 || err.Len() == 0 {
			t.Fatalf("unsupported feature: %s; code=%d out=%q err=%q", source, code, &out, &err)
		}
	}
}

func TestInterfaceInheritanceAndInstanceof(t *testing.T) {
	source := `<?php
interface A { public function a(); }
interface B extends A { public function b(); }
interface C { public function c(); }
class Base implements B { public function a() { return 1; } public function b() { return 2; } }
class Child extends Base implements C { public function c() { return 3; } }
$x = new Child;
echo $x->a(), $x->b(), $x->c(), '|', $x instanceof A, $x instanceof B, $x instanceof C;
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, err bytes.Buffer
		if code := RunWithVersion(source, "interface.php", &out, &err, version); code != 0 || out.String() != "123|111" || err.Len() != 0 {
			t.Fatalf("PHP %s: exit=%d stdout=%q stderr=%q", version, code, &out, &err)
		}
	}
}

func TestTraitCompositionIsPerClass(t *testing.T) {
	source := `<?php
trait BaseTrait { private $value = 1; private function read() { return $this->value; } public function get() { return $this->read(); } }
trait DerivedTrait { use BaseTrait; }
class A { use DerivedTrait; public function set($v) { $this->value = $v; } }
class B { use DerivedTrait; public function set($v) { $this->value = $v; } }
$a = new A; $b = new B; $a->set(5); echo $a->get(), '|', $b->get();
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, err bytes.Buffer
		if code := RunWithVersion(source, "trait.php", &out, &err, version); code != 0 || out.String() != "5|1" || err.Len() != 0 {
			t.Fatalf("PHP %s: exit=%d stdout=%q stderr=%q", version, code, &out, &err)
		}
	}
}

func TestTraitAdaptations(t *testing.T) {
	source := `<?php
trait A { public function value() { return 'A'; } }
trait B { public function value() { return 'B'; } }
class C { use A, B { A::value insteadof B; B::value as private other; } public function both() { return $this->value() . $this->other(); } }
echo (new C)->both();
`
	var out, err bytes.Buffer
	if code := RunWithVersion(source, "adapt.php", &out, &err, PHP84); code != 0 || out.String() != "AB" || err.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &err)
	}
}

func TestClassStaticMembersConstantsAndLateBinding(t *testing.T) {
	source := `<?php
class Base {
    public static $count = 1;
    public const LABEL = 'base';
    public static function bump() { static::$count++; return static::LABEL; }
    public static function make() { return new static; }
}

class Child extends Base {
    public static $count = 10;
    public const LABEL = 'child';
}
echo Base::bump(), ':', Base::$count, '|', Child::bump(), ':', Child::$count, '|', (Child::make() instanceof Child);
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, err bytes.Buffer
		if code := RunWithVersion(source, "static.php", &out, &err, version); code != 0 || out.String() != "base:2|child:11|1" || err.Len() != 0 {
			t.Fatalf("PHP %s: exit=%d stdout=%q stderr=%q", version, code, &out, &err)
		}
	}
}

func TestStaticForwardingKeepsCalledClass(t *testing.T) {
	source := `<?php
class A {
    public static function who() { return static::class; }
    public static function viaSelf() { return self::who(); }
    public static function viaExplicit() { return A::who(); }
}
class B extends A {
    public static function viaParent() { return parent::who(); }
    public static function viaStatic() { return static::who(); }
}
class C extends B {}
echo C::viaSelf(), '|', C::viaParent(), '|', C::viaStatic(), '|', C::viaExplicit();
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, err bytes.Buffer
		if code := RunWithVersion(source, "forwarding.php", &out, &err, version); code != 0 || out.String() != "C|C|C|A" || err.Len() != 0 {
			t.Fatalf("PHP %s: exit=%d stdout=%q stderr=%q", version, code, &out, &err)
		}
	}
}

func TestClassMethodResolvesTraitCollision(t *testing.T) {
	source := `<?php
trait T1 { public function value() { return 1; } }
trait T2 { public function value() { return 2; } }
class C { use T1, T2; public function value() { return 3; } }
echo (new C)->value();
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, err bytes.Buffer
		if code := RunWithVersion(source, "trait-collision.php", &out, &err, version); code != 0 || out.String() != "3" || err.Len() != 0 {
			t.Fatalf("PHP %s: exit=%d stdout=%q stderr=%q", version, code, &out, &err)
		}
	}
}

func TestSeparateTraitUseKeepsMethodOrigin(t *testing.T) {
	traits := `<?php
trait T1 { public function value() { return 1; } }
trait T2 { public function value() { return 2; } }
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for _, body := range []string{
			`class C { use T1; use T2; } echo (new C)->value();`,
			`trait Combined { use T1; use T2; } class C { use Combined; }`,
		} {
			var out, err bytes.Buffer
			code := RunWithVersion(traits+body, "separate.php", &out, &err, version)
			if code != 255 || out.Len() != 0 || !strings.Contains(err.String(), "Trait method T2::value") {
				t.Fatalf("PHP %s collision %q: exit=%d stdout=%q stderr=%q", version, body, code, &out, &err)
			}
		}
		var out, err bytes.Buffer
		code := RunWithVersion(traits+`class C { use T1; use T2; public function value() { return 3; } } echo (new C)->value();`, "resolved.php", &out, &err, version)
		if code != 0 || out.String() != "3" || err.Len() != 0 {
			t.Fatalf("PHP %s resolved: exit=%d stdout=%q stderr=%q", version, code, &out, &err)
		}
	}
}

func TestStaticStateAndClassConstantsDoNotLeakAcrossRuns(t *testing.T) {
	source := `<?php
class Cell {
    protected const START = 4;
    public static $value = self::START;
    public static function next() { return ++self::$value; }
}

echo Cell::next();
`
	for i := 0; i < 2; i++ {
		var out, err bytes.Buffer
		if code := RunWithVersion(source, "cell.php", &out, &err, PHP84); code != 0 || out.String() != "5" || err.Len() != 0 {
			t.Fatalf("run %d: exit=%d stdout=%q stderr=%q", i, code, &out, &err)
		}
	}
}

func TestTraitStaticCellsAndConstantVisibility(t *testing.T) {
	source := `<?php
trait CounterTrait {
    private static $n = 1;
    private const STEP = 2;
    public static function next() { self::$n += self::STEP; return self::$n; }
}
class CounterA { use CounterTrait; }
class CounterB { use CounterTrait; }
echo CounterA::next(), '|', CounterA::next(), '|', CounterB::next();
`
	var out, err bytes.Buffer
	if code := RunWithVersion(source, "counter.php", &out, &err, PHP84); code != 0 || out.String() != "3|5|3" || err.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &err)
	}
}

func TestClassFeatureFailuresDoNotRunPriorEcho(t *testing.T) {
	for _, source := range []string{
		"class A { final public const X = 1; } class B extends A { public const X = 2; }",
	} {
		var out, err bytes.Buffer
		if code := RunWithVersion("<?php echo 'before'; "+source, "failure.php", &out, &err, PHP84); code != 255 || out.Len() != 0 || err.Len() == 0 {
			t.Fatalf("source=%q exit=%d stdout=%q stderr=%q", source, code, &out, &err)
		}
	}
}

func TestClassCompositionErrorsAfterPriorEcho(t *testing.T) {
	for _, source := range []string{
		"interface Work { public function run(); } class Missing implements Work {}",
		"trait A { function f() {} } trait B { function f() {} } class C { use A, B; }",
	} {
		var out, err bytes.Buffer
		if code := RunWithVersion("<?php echo 'before'; "+source, "failure.php", &out, &err, PHP84); code != 255 || out.String() != "before" || err.Len() == 0 {
			t.Fatalf("source=%q exit=%d stdout=%q stderr=%q", source, code, &out, &err)
		}
	}
}

func TestUnselectedRelativeClassAndStaticUnset(t *testing.T) {
	for _, source := range []string{
		"if (false) { echo self::class; }",
		"class A { static $x = 1; } if (false) { unset(A::$x); }",
	} {
		var out, err bytes.Buffer
		if code := RunWithVersion("<?php echo 'before'; "+source+" echo 'after';", "unselected.php", &out, &err, PHP84); code != 0 || out.String() != "beforeafter" || err.Len() != 0 {
			t.Fatalf("source=%q exit=%d stdout=%q stderr=%q", source, code, &out, &err)
		}
	}
}

func TestSelectedStaticUnsetStillFails(t *testing.T) {
	var out, err bytes.Buffer
	code := RunWithVersion("<?php class A { static $x = 1; } echo 'before'; unset(A::$x); echo 'after';", "selected.php", &out, &err, PHP84)
	if code != 255 || out.String() != "before" || !strings.Contains(err.String(), "Attempt to unset static property") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &err)
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
