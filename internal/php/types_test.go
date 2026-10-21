package php

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSimpleParameterAndReturnTypes(t *testing.T) {
	for _, tc := range []struct {
		name, source, output, diagnostic string
		exit                             int
	}{
		{"weak coercion", `<?php function total(int $a, float $b): float { return $a + $b; } echo total("2", 3);`, "5", "", 0},
		{"strict scalars", `<?php declare(strict_types=1); function total(int $a, float $b): float { return $a + $b; } echo total(2, 3);`, "5", "", 0},
		{"strict argument failure", `<?php declare(strict_types=1); function id(int $x): int { return $x; } echo id("2");`, "", "Uncaught TypeError: id(): Argument #1 ($x) must be of type int, string given", 255},
		{"strict return failure", `<?php declare(strict_types=1); function id(): int { return "2"; } echo id();`, "", "Uncaught TypeError: id(): Return value must be of type int, string returned", 255},
		{"nullable class", `<?php class Box {} function inspect(?Box $x): bool { return $x === null; } echo inspect(null), inspect(new Box);`, "1", "", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, version := range []Version{PHP83, PHP84, PHP85} {
				var out, diagnostics bytes.Buffer
				code := RunWithVersion(tc.source, "typed.php", &out, &diagnostics, version)
				if code != tc.exit || out.String() != tc.output || (tc.diagnostic == "" && diagnostics.Len() != 0) || (tc.diagnostic != "" && !strings.Contains(diagnostics.String(), tc.diagnostic)) {
					t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
				}
			}
		})
	}
}

func TestUnionArgumentsReturnsAndDeclarationChecks(t *testing.T) {
	for _, tc := range []struct {
		name, source, output, diagnostic string
		exit                             int
	}{
		{"exact scalar wins", `<?php function f(int|string $x): int|string { return $x; } echo gettype(f("42")), "|", gettype(f(42));`, "string|integer", "", 0},
		{"numeric string preference", `<?php function f(float|int $x): float|int { return $x; } echo gettype(f("42")), "|", gettype(f("42.5"));`, "integer|double", "", 0},
		{"strict no string coercion", `<?php declare(strict_types=1); function f(int|float $x) {} f("42");`, "", "must be of type int|float, string given", 255},
		{"order independent contract", `<?php interface I { public function f(int|string $x): string|int; } class C implements I { public function f(string|int $x): int|string { return $x; } } echo (new C)->f(3);`, "3", "", 0},
		{"nullable spelling contract", `<?php interface I { public function f(?int $x): ?int; } class C implements I { public function f(int|null $x): int|null { return $x; } } echo (new C)->f(3);`, "3", "", 0},
		{"incompatible union contract", `<?php interface I { public function f(int|string $x): int|string; } class C implements I { public function f(int|float $x): int|string { return $x; } }`, "", "must be compatible", 255},
		{"duplicate aliases", `<?php namespace N; class A {} use N\A as B; function f(A|B $x) {} echo 'unreachable';`, "", "Duplicate type", 255},
		{"true false redundant", `<?php function f(true|false $x) {} echo 'unreachable';`, "", "Type contains both true and false, bool must be used instead", 255},
		{"object class redundant", `<?php class A {} function f(object|A $x) {} echo 'unreachable';`, "", "contains both object and a class type, which is redundant", 255},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, version := range []Version{PHP83, PHP84, PHP85} {
				var out, diagnostics bytes.Buffer
				code := RunWithVersion(tc.source, "union.php", &out, &diagnostics, version)
				if code != tc.exit || out.String() != tc.output || (tc.diagnostic == "" && diagnostics.Len() != 0) || (tc.diagnostic != "" && !strings.Contains(diagnostics.String(), tc.diagnostic)) {
					t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
				}
			}
		})
	}
}

func TestDeclaredUnionDiagnosticOrder(t *testing.T) {
	for _, tc := range []struct{ source, want string }{
		{"int|string", "string|int"},
		{"string|int", "string|int"},
		{"int|string|null", "string|int|null"},
		{"null|int", "?int"},
		{"?int", "?int"},
		{"Foo|int|string", "Foo|string|int"},
	} {
		if got := declaredTypeText(tc.source, false); got != tc.want {
			t.Errorf("declaredTypeText(%q) = %q, want %q", tc.source, got, tc.want)
		}
	}
	for _, source := range []string{"int|string", "string|int"} {
		program := `<?php function take(` + source + ` $x) {} take([]);`
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			var out, diagnostics bytes.Buffer
			if code := RunWithVersion(program, "union.php", &out, &diagnostics, version); code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "must be of type string|int, array given") {
				t.Errorf("%s %s: exit=%d stdout=%q stderr=%q", version, source, code, out.String(), diagnostics.String())
			}
		}
	}
}

func TestIntersectionAndDNFTypes(t *testing.T) {
	for _, tc := range []struct{ source, want string }{
		{"A&B", "A&B"},
		{"(A&B)|null", "(A&B)|null"},
		{"int|(B&A)|string", "(B&A)|string|int"},
	} {
		if got := declaredTypeText(tc.source, false); got != tc.want {
			t.Errorf("declaredTypeText(%q) = %q, want %q", tc.source, got, tc.want)
		}
	}
	for _, tc := range []struct {
		name, source, output, diagnostic string
		exit                             int
	}{
		{"intersection", `<?php interface A {} interface B {} class Both implements A, B {} function f(A&B $x): B&A { return $x; } echo gettype(f(new Both));`, "object", "", 0},
		{"dnf scalar and null", `<?php interface A {} interface B {} class Both implements A, B {} function f((A&B)|string|null $x): (B&A)|string|null { return $x; } echo gettype(f(new Both)), "|", gettype(f("x")), "|", gettype(f(null));`, "object|string|NULL", "", 0},
		{"intersection rejects partial", `<?php interface A {} interface B {} class One implements A {} function f(A&B $x) {} f(new One);`, "", "must be of type A&B", 255},
		{"dnf rejects array", `<?php interface A {} interface B {} function f((A&B)|int $x) {} f([]);`, "", "must be of type (A&B)|int", 255},
		{"same intersection contract", `<?php interface A {} interface B {} interface I { public function f(A&B $x): A&B; } class C implements A, B, I { public function f(B&A $x): B&A { return $x; } } echo gettype((new C)->f(new C));`, "object", "", 0},
		{"builtin intersection invalid", `<?php function f(int&A $x) {} echo "bad";`, "", "cannot be part of an intersection type", 255},
		{"duplicate intersection invalid", `<?php function f(A&A $x) {} echo "bad";`, "", "Duplicate type", 255},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, version := range []Version{PHP83, PHP84, PHP85} {
				var out, diagnostics bytes.Buffer
				code := RunWithVersion(tc.source, "intersection.php", &out, &diagnostics, version)
				if code != tc.exit || out.String() != tc.output || (tc.diagnostic == "" && diagnostics.Len() != 0) || (tc.diagnostic != "" && !strings.Contains(diagnostics.String(), tc.diagnostic)) {
					t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
				}
			}
		})
	}
}

func TestNullableShorthandCannotJoinUnion(t *testing.T) {
	source := "<?php\necho \"before\\n\";\nfunction wrong(?int|string $value) {}\n"
	want := "Parse error: syntax error, unexpected token \"|\", expecting variable in union.php on line 3\n"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		code := RunWithVersion(source, "union.php", &out, &diagnostics, version)
		if code != 255 || out.Len() != 0 || diagnostics.String() != want {
			t.Errorf("%s: exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
		}
	}
}

func TestWeakIntParameterReportsPrecisionLossAtDeclaration(t *testing.T) {
	source := "<?php\nfunction take(int $value): int { return $value; }\necho take(3.5);\n"
	want := "Deprecated: Implicit conversion from float 3.5 to int loses precision in typed.php on line 2\n"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		if code := RunWithVersion(source, "typed.php", &out, &diagnostics, version); code != 0 || out.String() != "3" || diagnostics.String() != want {
			t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
		}
	}
}

func TestStrictTypesFollowCallingFile(t *testing.T) {
	dir := t.TempDir()
	main := filepath.Join(dir, "main.php")
	lib := filepath.Join(dir, "lib.inc")
	if err := os.WriteFile(lib, []byte("<?php declare(strict_types=1); function take(int $x): int { return $x; }"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		source, output string
		exit           int
	}{
		{`<?php require __DIR__ . '/lib.inc'; echo take("2");`, "2", 0},
		{`<?php declare(strict_types=1); require __DIR__ . '/lib.inc'; echo take("2");`, "", 255},
	} {
		var out, diagnostics bytes.Buffer
		code := Run(tc.source, main, &out, &diagnostics)
		if code != tc.exit || out.String() != tc.output || (tc.exit == 0 && diagnostics.Len() != 0) || (tc.exit != 0 && !strings.Contains(diagnostics.String(), "TypeError")) {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
		}
	}
}

func TestStrictTypesDoNotLeakBetweenRuns(t *testing.T) {
	for _, tc := range []struct {
		source string
		exit   int
		output string
	}{
		{`<?php declare(strict_types=1); function f(int $x) { return $x; } f("4");`, 255, ""},
		{`<?php function f(int $x) { return $x; } echo f("4");`, 0, "4"},
	} {
		var out, diagnostics bytes.Buffer
		code := Run(tc.source, "fresh.php", &out, &diagnostics)
		if code != tc.exit || out.String() != tc.output {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
		}
	}
}

func TestSimpleTypedMethodVariance(t *testing.T) {
	for _, tc := range []struct {
		source, output string
		exit           int
	}{
		{`<?php class A { public function take(int $x): int { return $x; } } class B extends A { public function take($x): int { return 7; } } echo (new B)->take("3");`, "7", 0},
		{`<?php class A { public function take($x): int { return 1; } } class B extends A { public function take(int $x): int { return $x; } }`, "", 255},
		{`<?php interface I { public function take(int $x): int; } class C implements I { public function take(string $x): int { return 1; } }`, "", 255},
	} {
		var out, diagnostics bytes.Buffer
		code := Run(tc.source, "variance.php", &out, &diagnostics)
		if code != tc.exit || out.String() != tc.output || (tc.exit == 0 && diagnostics.Len() != 0) || (tc.exit != 0 && !strings.Contains(diagnostics.String(), "must be compatible")) {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
		}
	}
}

func TestCompositeTypeDeclarationContainment(t *testing.T) {
	e := &evaluator{}
	e.initClasses()
	a := &classDecl{name: "A", kind: "interface"}
	b := &classDecl{name: "B", kind: "interface"}
	c := &classDecl{name: "C", kind: "class", interfaces: []*classDecl{a, b}}
	d := &classDecl{name: "D", kind: "class", parent: c}
	e.classes["a"], e.classes["b"], e.classes["c"], e.classes["d"] = a, b, c, d
	for _, tc := range []struct {
		receiver, source string
		accept           bool
	}{
		{"A", "A&B", true}, {"A&B", "A&B&C", true},
		{"A&B", "C", true}, {"A&B", "D", true}, {"C", "A&B", false},
		{"A|B", "A", true}, {"A", "A|B", false},
		{"(A&B)|int", "C|int", true}, {"(A&B)|int", "A|int", false},
		{"A|int", "(A&B)|int", true}, {"(A&B)|int", "(B&A)|int", true},
		{"object", "A", true}, {"object", "C", true}, {"A", "object", false},
		{"bool", "true|false", true}, {"true|false", "bool", false},
		{"true", "bool", false}, {"float", "int", false},
		{"float|string", "int|string", false},
		{"mixed", "A|null", true}, {"A|null", "mixed", false},
		{"A|null", "null", true}, {"A", "null", false},
		{"A", "A", true}, {"Missing", "Missing", true},
		{"A", "Missing", false},
	} {
		if got := e.typeAccepts(tc.receiver, false, tc.source, false); got != tc.accept {
			t.Errorf("%s accepts %s = %v, want %v", tc.receiver, tc.source, got, tc.accept)
		}
	}
	if !e.typeAccepts("A|null", false, "A", true) || !e.typeAccepts("A", true, "A|null", false) {
		t.Error("nullable shorthand and explicit null must be equivalent")
	}
}

func TestInheritedInterfaceContractsAreAllChecked(t *testing.T) {
	base := `<?php
interface Two { public function load($first, $second = null); }
interface One { public function load($first); }
interface Both extends Two, One {}
`
	for _, tc := range []struct {
		body string
		exit int
	}{
		{`class Good implements Both { public function load($first, $second = null) { echo $first, ":", $second; } } (new Good)->load("x", "y");`, 0},
		{`class Bad implements Both { public function load($first, $second, $third) {} }`, 255},
	} {
		var out, diagnostics bytes.Buffer
		code := Run(base+tc.body, "diamond.php", &out, &diagnostics)
		if code != tc.exit || (tc.exit == 0 && (out.String() != "x:y" || diagnostics.Len() != 0)) || (tc.exit != 0 && !strings.Contains(diagnostics.String(), "must be compatible with Two::load")) {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
		}
	}
}

func TestConflictingInterfaceContractsFailAtDeclaration(t *testing.T) {
	source := `<?php
interface Two { public function load($first, $second); }
interface One { public function load($first); }
echo "before\n";
interface Both extends Two, One {}
echo "after\n";
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		code := RunWithVersion(source, "diamond.php", &out, &diagnostics, version)
		if code != 255 || out.String() != "before\n" || !strings.Contains(diagnostics.String(), "Declaration of Two::load($first, $second) must be compatible with One::load($first)") || strings.Contains(diagnostics.String(), "Both::load") {
			t.Fatalf("version=%s exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
		}
	}
}

func TestVoidNeverReturnVariance(t *testing.T) {
	e := &evaluator{}
	for _, test := range []struct {
		parent, child string
		compatible   bool
	}{
		{"void", "void", true}, {"void", "never", true},
		{"mixed", "never", true}, {"int|string", "never", true},
		{"", "void", true}, {"", "never", true},
		{"mixed", "void", false}, {"void", "mixed", false},
		{"int", "void", false}, {"never", "void", false},
	} {
		if got := e.typeAccepts(test.parent, false, test.child, false); got != test.compatible {
			t.Errorf("parent %q child %q: compatible=%v, want %v", test.parent, test.child, got, test.compatible)
		}
	}
}
