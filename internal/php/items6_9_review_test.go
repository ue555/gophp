package php

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClassParentLoadedBeforeDeclaration(t *testing.T) {
	dir := t.TempDir()
	main := filepath.Join(dir, "main.php")
	base := filepath.Join(dir, "base.inc")
	if err := os.WriteFile(base, []byte("<?php class Base { public function value() { return 7; } }"), 0600); err != nil {
		t.Fatal(err)
	}
	for name, tc := range map[string]struct {
		source string
		out    string
		exit   int
	}{
		"require then child":     {"<?php require __DIR__ . '/base.inc'; class Child extends Base {} echo (new Child)->value();", "7", 0},
		"function loads parent":  {"<?php function loadBase() { require __DIR__ . '/base.inc'; } loadBase(); class Child extends Base {} echo (new Child)->value();", "7", 0},
		"switch loads parent":    {"<?php switch (1) { case 1: require __DIR__ . '/base.inc'; break; } class Child extends Base {} echo (new Child)->value();", "7", 0},
		"switch skips parent":    {"<?php switch (0) { case 1: require __DIR__ . '/base.inc'; break; } class Child extends Base {}", "", 255},
		"call skips parent":      {"<?php function loadBase() { require __DIR__ . '/base.inc'; } class Child extends Base {}", "", 255},
		"same file forward":      {"<?php class Child extends Base {} class Base { public function value() { return 7; } } echo (new Child)->value();", "7", 0},
		"same file parent first": {"<?php class Base { public function value() { return 7; } } class Child extends Base {} echo (new Child)->value();", "7", 0},
		"conditional parent":     {"<?php if (true) { class Base { public function value() { return 7; } } } class Child extends Base {} echo (new Child)->value();", "7", 0},
		"missing parent":         {"<?php class Child extends Missing {}", "", 255},
		"cycle":                  {"<?php class Base extends Child {} class Child extends Base {}", "", 255},
	} {
		t.Run(name, func(t *testing.T) {
			var out, diagnostics bytes.Buffer
			code := Run(tc.source, main, &out, &diagnostics)
			if code != tc.exit || out.String() != tc.out || (tc.exit == 0 && diagnostics.Len() != 0) {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
			}
		})
	}
}

func TestBindingUnsetThenRedeclare(t *testing.T) {
	source := `<?php
$x = 10;
function f() {
    static $x = 1;
    global $x;
    unset($x);
    $x = 3;
    echo $x, '|';
    static $x = 1;
    echo $x, '|';
    global $x;
    echo $x, '|';
}
f(); echo $x;
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "binding.php", &out, &diagnostics)
	if code != 0 || out.String() != "3|1|10|10" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestBindingDeclarationOrderAndFrameIsolation(t *testing.T) {
	source := `<?php
$x = 10;
function first() { static $x = 1; global $x; $x++; echo $x, '|'; }
function second() { global $x; static $x = 1; $x++; echo $x, '|'; }
first(); first(); second(); second(); echo $x;
`
	for i := 0; i < 2; i++ {
		var out, diagnostics bytes.Buffer
		code := Run(source, "binding.php", &out, &diagnostics)
		if code != 0 || out.String() != "11|12|2|3|12" || diagnostics.Len() != 0 {
			t.Fatalf("run %d: exit=%d stdout=%q stderr=%q", i, code, out.String(), diagnostics.String())
		}
	}
}

func TestIncludeChangesBindingInCallerFrame(t *testing.T) {
	dir := t.TempDir()
	main := filepath.Join(dir, "main.php")
	child := filepath.Join(dir, "bind.inc")
	if err := os.WriteFile(child, []byte("<?php global $x; $x++;"), 0600); err != nil {
		t.Fatal(err)
	}
	source := "<?php $x = 10; function f() { static $x = 1; include __DIR__ . '/bind.inc'; echo $x; } f(); echo '|', $x;"
	var out, diagnostics bytes.Buffer
	code := Run(source, main, &out, &diagnostics)
	if code != 0 || out.String() != "11|11" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
	if strings.Contains(diagnostics.String(), "panic") {
		t.Fatal("Go panic escaped")
	}
}

func TestFinalAndOverrideChecks(t *testing.T) {
	for name, source := range map[string]string{
		"extend final": `<?php final class A {} class B extends A {}`,
		"override final": `<?php class A { final public function f() {} } class B extends A { public function f() {} }`,
		"narrow method": `<?php class A { public function f() {} } class B extends A { protected function f() {} }`,
		"narrow property": `<?php class A { protected $x; } class B extends A { private $x; }`,
		"unimplemented abstract": `<?php abstract class A { abstract public function f(); } class B extends A {}`,
		"instantiate abstract": `<?php abstract class A {} new A;`,
	} {
		t.Run(name, func(t *testing.T) {
			var out, diagnostics bytes.Buffer
			code := Run(source, "class-check.php", &out, &diagnostics)
			if code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "Fatal error:") {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
			}
		})
	}
}
