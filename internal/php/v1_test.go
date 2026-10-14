package php

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// These cases intentionally describe the Go subset, not differential fixtures:
// official PHP supports references, object iteration and additional casts.
func TestV1UnsupportedForms(t *testing.T) {
	for _, source := range []string{
		"foreach ([] as &$v) {}", "foreach ([] as &$k => $v) {}",
		"foreach ([[1]] as [$v]) {}", "foreach ([] as list($v)) {}",
		"foreach ([] as $v): endforeach;",
		"echo (integer)1;", "echo (boolean)1;", "echo (double)1;", "echo (real)1;",
		"echo (binary)1;", "echo (array)1;", "echo (object)1;", "echo (unset)1;",
		"echo count_normal;", "echo PHP_VERSION;", "echo isset();", "empty(1, 2);", "unset();",
	} {
		t.Run(source, func(t *testing.T) {
			var out, diagnostics bytes.Buffer
			code := Run("<?php\necho 'unreachable';\nif (false) { "+source+" }", "input.php", &out, &diagnostics)
			if code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), "Parse error:") {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &diagnostics)
			}
		})
	}
	for _, object := range []string{"new stdClass", "new C"} {
		var out, diagnostics bytes.Buffer
		code := Run("<?php\nclass C {} foreach ("+object+" as $v) {} echo 'unreachable';", "input.php", &out, &diagnostics)
		want := "Fatal error: Uncaught Error: gophp: foreach over objects is not supported in input.php:2\nStack trace:\n#0 {main}\n  thrown in input.php on line 2\n"
		if code != 255 || out.Len() != 0 || diagnostics.String() != want {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &diagnostics)
		}
	}
}

func TestV1MalformedBeforeExecution(t *testing.T) {
	for _, source := range []string{
		"foreach () {}", "foreach ([] $v) {}", "foreach ([] as) {}", "foreach ([] as $k =>) {}",
		"foreach ([1] as $v) {", "isset($x,); echo ;", "empty();", "unset(1 + 2);",
		"echo (int);", "isset($x + 1);", "empty($x[]);",
	} {
		t.Run(source, func(t *testing.T) {
			var out, diagnostics bytes.Buffer
			code := Run("<?php echo 'unreachable'; "+source, "input.php", &out, &diagnostics)
			if code != 255 || out.Len() != 0 || diagnostics.Len() == 0 {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &diagnostics)
			}
		})
	}
}

func TestBuiltinNamesReserved(t *testing.T) {
	for _, name := range []string{"strlen", "count", "gettype", "is_int", "is_float", "is_string", "is_bool", "is_null", "is_array", "is_object", "is_numeric"} {
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			t.Run(name+"/"+string(version), func(t *testing.T) {
				var out, diagnostics bytes.Buffer
				code := RunWithVersion("<?php\necho 'unreachable'; function "+strings.ToUpper(name)+"() {}", "input.php", &out, &diagnostics, version)
				word := ""
				if version != PHP83 {
					word = "function "
				}
				want := fmt.Sprintf("Fatal error: Cannot redeclare %s%s() in input.php on line 2\n", word, strings.ToUpper(name))
				if version == PHP85 {
					want += "Stack trace:\n#0 {main}\n"
				}
				if code != 255 || out.Len() != 0 || diagnostics.String() != want {
					t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &diagnostics)
				}
			})
		}
	}
}
