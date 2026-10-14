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
		exit int
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
