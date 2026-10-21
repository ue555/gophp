package php

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestLoopControlValidationBeforeExecution(t *testing.T) {
	for _, source := range []string{
		"if (false) { break; }",
		"for (;;) { continue 0; }",
		"while (false) { break +1; }",
		"while (false) { break -1; }",
		"while (false) { continue true; }",
		"while (false) { continue null; }",
		"while (false) { break 9223372036854775808; }",
		"do { continue 2; } while (false);",
		"while (true) { function nested() { continue; } break; }",
	} {
		t.Run(source, func(t *testing.T) {
			var out, err bytes.Buffer
			code := Run("<?php echo 'must not execute'; "+source, "input.php", &out, &err)
			if code != 255 || out.Len() != 0 || !strings.Contains(err.String(), "Fatal error:") {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &err)
			}
		})
	}
}

func TestExecutionBudgetForLoopFuzzing(t *testing.T) {
	for _, source := range []string{
		"while (true) {}", "for (;;) continue;", "do {} while (true);",
		"while (true) { for (;;) { break; } }",
	} {
		var err bytes.Buffer
		code := run("<?php "+source, "input.php", io.Discard, &err, PHP84, 10)
		if code != 255 || !strings.Contains(err.String(), "execution budget exhausted") {
			t.Fatalf("%s: exit=%d stderr=%s", source, code, &err)
		}
	}
}

func TestFunctionScopesAndRegistryAreFresh(t *testing.T) {
	source := "<?php function fresh($n = 0) { return ++$n; } echo fresh(), fresh();"
	for i := 0; i < 2; i++ {
		var out, err bytes.Buffer
		if code := Run(source, "input.php", &out, &err); code != 0 || out.String() != "11" || err.Len() != 0 {
			t.Fatalf("run %d: exit=%d stdout=%q stderr=%q", i, code, &out, &err)
		}
	}
	var err bytes.Buffer
	if code := Run("<?php fresh();", "input.php", io.Discard, &err); code != 255 || !strings.Contains(err.String(), "Call to undefined function fresh()") {
		t.Fatalf("registry leaked: exit=%d stderr=%s", code, &err)
	}
}

func TestFunctionRecursionLimit(t *testing.T) {
	var err bytes.Buffer
	code := Run("<?php function recur() { return recur(); } recur();", "input.php", io.Discard, &err)
	if code != 255 || !strings.Contains(err.String(), "depth exceeded (256)") {
		t.Fatalf("exit=%d stderr=%s", code, &err)
	}
}

func TestUnsupportedFunctionAndLoopSyntax(t *testing.T) {
	for _, source := range []string{
		"function f(...$args, $tail) {}", "function f(int&float $x) {}",
		"function f(): int&float {}",
		"while (false): endwhile;", "for (;;): endfor;",
	} {
		var out, err bytes.Buffer
		code := Run("<?php echo 'must not execute'; "+source, "input.php", &out, &err)
		if code != 255 || out.Len() != 0 || err.Len() == 0 {
			t.Errorf("%s: exit=%d stdout=%q stderr=%q", source, code, &out, &err)
		}
	}
}
