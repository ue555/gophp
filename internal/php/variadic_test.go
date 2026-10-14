package php

import (
	"bytes"
	"strings"
	"testing"
)

func TestVariadicParameterEvaluatesArgumentsOnce(t *testing.T) {
	source := `<?php
$counter = 0;
function nextValue() { global $counter; return ++$counter; }
function collect($first, ...$rest) {
    echo $first, ':', count($rest), ':';
    foreach ($rest as $value) echo $value;
}
collect(nextValue(), nextValue(), nextValue());
echo '|', $counter;
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "variadic.php", &out, &diagnostics)
	if code != 0 || out.String() != "1:2:23|3" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestVariadicParameterInClosureAndMethod(t *testing.T) {
	source := `<?php
class C { public function size(...$items) { return count($items); } }
$f = function (...$items) { return count($items); };
echo $f(), '|', $f(1, 2), '|', (new C)->size(1, 2, 3);
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "variadic.php", &out, &diagnostics)
	if code != 0 || out.String() != "0|2|3" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestTypedVariadicParameter(t *testing.T) {
	source := `<?php
function total(int ...$items): int {
    $sum = 0;
    foreach ($items as $item) $sum += $item;
    return $sum;
}
echo total(1, 2, 3);
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "variadic.php", &out, &diagnostics)
	if code != 0 || out.String() != "6" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestInvalidVariadicDeclarationRejectedBeforeExecution(t *testing.T) {
	for _, body := range []string{"function f(...$items, $tail) {}", "function f(...$items = []) {}"} {
		var out, diagnostics bytes.Buffer
		code := Run("<?php\necho 'before';\n"+body+"\n", "variadic.php", &out, &diagnostics)
		if code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "variadic parameter") {
			t.Fatalf("source=%q exit=%d stdout=%q stderr=%q", body, code, out.String(), diagnostics.String())
		}
	}
}

func TestArgumentUnpackOrderAndSingleEvaluation(t *testing.T) {
	source := `<?php
$counter = 0;
function values() { global $counter; $counter++; return [2, 3]; }
function combine($a, $b, $c) { return "$a$b$c"; }
echo combine(1, ...values()), '|', $counter;
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "unpack.php", &out, &diagnostics)
	if code != 0 || out.String() != "123|1" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestStringKeyArgumentUnpackBindsByName(t *testing.T) {
	var out, diagnostics bytes.Buffer
	code := Run("<?php\nfunction f($x) { echo $x; }\nf(...['x' => 1]);\n", "unpack.php", &out, &diagnostics)
	if code != 0 || out.String() != "1" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestPositionalAfterUnpackRejectedBeforeExecution(t *testing.T) {
	cases := []string{
		"function f(...$xs) {} f(...[1], 2);",
		"function f(...$xs) {} if (false) { f(...[1], 2); }",
		"function f(...$xs) {} f(...[1], ...[2], 3);",
		"class C { public function f(...$xs) {} } (new C)->f(...[1], 2);",
		"class C { public static function f(...$xs) {} } C::f(...[1], 2);",
		"class C { public function __construct(...$xs) {} } new C(...[1], 2);",
	}
	for _, source := range cases {
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			var out, diagnostics bytes.Buffer
			code := RunWithVersion("<?php\necho 'before';\n"+source, "unpack.php", &out, &diagnostics, version)
			if code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "Cannot use positional argument after argument unpacking") {
				t.Fatalf("version=%s source=%q exit=%d stdout=%q stderr=%q", version, source, code, out.String(), diagnostics.String())
			}
		}
	}
}

func TestMultipleUnpacksAndPositionalsBeforeThem(t *testing.T) {
	source := `<?php
function f(...$xs) { echo count($xs), '|'; foreach ($xs as $x) echo $x; }
f(0, ...[1], ...[2]);
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "unpack.php", &out, &diagnostics)
	if code != 0 || out.String() != "3|012" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}
