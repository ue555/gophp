package php

import (
	"bytes"
	"strings"
	"testing"
)

func TestNamedArgumentsAndUnpacking(t *testing.T) {
	cases := []struct{ source, stdout string }{
		{`function f($a, $b = 2, $c = 3) { echo $a, $b, $c; } f(c: 7, a: 1);`, "127"},
		{`function f($a, $b, ...$rest) { echo $a, $b, '|'; foreach ($rest as $key => $v) echo $key, '=', $v, ';'; } f(1, b: 2, extra: 3);`, "12|extra=3;"},
		{`function f($a, $b, $c) { echo $a, $b, $c; } f(...['b' => 2, 'a' => 1], c: 3);`, "123"},
		{`function f($a, $b, $c) { echo $a, $b, $c; } f(...[1], ...['c' => 3], b: 2);`, "123"},
		{`class C { public function __construct($a, $b) { echo $a, $b; } public static function s($a, $b) { echo $a, $b; } public function m($a, $b) { echo $a, $b; } } $c = new C(b: 2, a: 1); $c->m(b: 4, a: 3); C::s(b: 6, a: 5);`, "123456"},
		{`$f = function($a, $b) { echo $a, $b; }; $f(b: 2, a: 1);`, "12"},
		{`echo strlen(string: 'abc'), '|', count(mode: COUNT_NORMAL, value: [1, 2]);`, "3|2"},
		{`$i = 0; function nextValue() { global $i; return ++$i; } function f($a, $b, $c) { echo $a, $b, $c; } f(c: nextValue(), a: nextValue(), b: nextValue()); echo '|', $i;`, "231|3"},
	}
	for _, tc := range cases {
		var out, diagnostics bytes.Buffer
		code := Run("<?php "+tc.source, "named.php", &out, &diagnostics)
		if code != 0 || out.String() != tc.stdout || diagnostics.Len() != 0 {
			t.Errorf("source=%s exit=%d stdout=%q stderr=%q", tc.source, code, out.String(), diagnostics.String())
		}
	}
}

func TestNamedArgumentFailures(t *testing.T) {
	cases := []struct{ source, diagnostic, stdout string }{
		{`function f($a) {} f(a: 1, a: 2);`, "Named parameter $a overwrites previous argument", ""},
		{`function f($a) {} f(1, a: 2);`, "Named parameter $a overwrites previous argument", ""},
		{`function f($a) {} f(b: 1);`, "Unknown named parameter $b", ""},
		{`function f($a, $b) {} f(b: 2);`, "Too few arguments", ""},
		{`echo 'before'; function f($a, $b) {} if (false) { f(a: 1, 2); }`, "Cannot use positional argument after named argument", ""},
		{`function f($a, $b) {} f(...['a' => 1, 2]);`, "Cannot use positional argument after named argument", ""},
		{`function f($a) {} f(...['a' => 1], a: 2);`, "Named parameter $a overwrites previous argument", ""},
	}
	for _, tc := range cases {
		var out, diagnostics bytes.Buffer
		code := Run("<?php "+tc.source, "named.php", &out, &diagnostics)
		if code != 255 || out.String() != tc.stdout || !strings.Contains(diagnostics.String(), tc.diagnostic) {
			t.Errorf("source=%s exit=%d stdout=%q stderr=%q", tc.source, code, out.String(), diagnostics.String())
		}
	}
}

func TestUnpackAfterNamedRejectedBeforeExecution(t *testing.T) {
	for _, source := range []string{
		`function f($a, $b) {} f(a: 1, ...['b' => 2]);`,
		`function f($a, $b) {} echo 'before'; if (false) { f(a: 1, ...['b' => 2]); }`,
		`class C { function f($a, $b) {} } (new C)->f(a: 1, ...['b' => 2]);`,
		`$f = function($a, $b) {}; $f(a: 1, ...['b' => 2]);`,
	} {
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			var out, diagnostics bytes.Buffer
			code := RunWithVersion("<?php "+source, "named.php", &out, &diagnostics, version)
			if code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "Cannot use argument unpacking after named arguments") {
				t.Errorf("PHP %s source=%s exit=%d stdout=%q stderr=%q", version, source, code, out.String(), diagnostics.String())
			}
		}
	}
}

func TestNamedBindingFailureDoesNotEnterCalleeFrame(t *testing.T) {
	for _, source := range []string{
		`function target($value) {} target(unknown: 1);`,
		`function target($value) {} target(value: 1, value: 2);`,
	} {
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			var out, diagnostics bytes.Buffer
			code := RunWithVersion("<?php "+source, "named.php", &out, &diagnostics, version)
			if code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "Stack trace:\n#0 {main}\n") || strings.Contains(diagnostics.String(), ": target(") {
				t.Errorf("PHP %s source=%s exit=%d stdout=%q stderr=%q", version, source, code, out.String(), diagnostics.String())
			}
		}
	}
	var out, diagnostics bytes.Buffer
	source := `<?php function target($value) {} function outer() { target(unknown: 1); } outer();`
	code := Run(source, "named.php", &out, &diagnostics)
	if code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "outer()") || strings.Contains(diagnostics.String(), "target(1)") {
		t.Errorf("nested exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}
