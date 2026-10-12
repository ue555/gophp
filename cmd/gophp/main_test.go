package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.php")
	if err := os.WriteFile(file, []byte("<?php echo 'ok';"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name            string
		args            []string
		exit            int
		out, diagnostic string
	}{
		{"success", []string{file}, 0, "ok", ""},
		{"usage", nil, 2, "", "usage:"},
		{"extra arguments", []string{file, file}, 2, "", "usage:"},
		{"missing file", []string{filepath.Join(dir, "missing.php")}, 1, "", "missing.php"},
		{"version 8.3", []string{"--php-version=8.3", file}, 0, "ok", ""},
		{"version 8.4", []string{"--php-version", "8.4", file}, 0, "ok", ""},
		{"version 8.5", []string{"--php-version=8.5", file}, 0, "ok", ""},
		{"unsupported version", []string{"--php-version=8.6", file}, 2, "", "unsupported PHP version"},
		{"empty version", []string{"--php-version=", file}, 2, "", "unsupported PHP version"},
		{"missing version", []string{"--php-version"}, 2, "", "flag needs an argument"},
		{"unknown option", []string{"--unknown", file}, 2, "", "flag provided but not defined"},
		{"help", []string{"--help"}, 0, "", "usage:"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(test.args, &stdout, &stderr)
			if code != test.exit || stdout.String() != test.out {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
			if test.diagnostic == "" && stderr.Len() != 0 || !strings.Contains(stderr.String(), test.diagnostic) {
				t.Errorf("unexpected stderr %q", &stderr)
			}
		})
	}
}

// The warning text and line number are from the user's PHP 8.5.11 comparison.
// Exercise the same CLI entry point used by compare-all, including its default.
func TestCLINaNVersion(t *testing.T) {
	file, err := filepath.Abs("../../tests/fixtures/success/truth_true.php")
	if err != nil {
		t.Fatal(err)
	}
	for _, version := range []string{"8.5", "8.4", "8.3", ""} {
		t.Run("target="+version, func(t *testing.T) {
			args := []string{file}
			if version != "" {
				args = append([]string{"--php-version=" + version}, args...)
			}
			var stdout, stderr bytes.Buffer
			code := run(args, &stdout, &stderr)
			wantWarning := ""
			if version == "8.5" {
				wantWarning = "Warning: unexpected NAN value was coerced to bool in " + file + " on line 21\n"
			}
			if code != 0 || stdout.String() != strings.Repeat("1", 20)+"\n" || stderr.String() != wantWarning {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want stderr=%q", code, &stdout, &stderr, wantWarning)
			}
		})
	}
}

func TestCLIV1ProgramAndDiagnostics(t *testing.T) {
	source, err := os.ReadFile("../../tests/fixtures/success/v1_report.php")
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "集計 report.php")
	if err = os.WriteFile(file, source, 0600); err != nil {
		t.Fatal(err)
	}
	for _, version := range []string{"8.3", "8.4", "8.5"} {
		var out, diagnostics bytes.Buffer
		code := run([]string{"--php-version=" + version, "--", file}, &out, &diagnostics)
		if code != 0 || out.String() != "a:tea:12:3\nc:水:3:3\n2:2:1:0:integer\n" || diagnostics.Len() != 0 {
			t.Fatalf("%s exit=%d out=%q err=%q", version, code, &out, &diagnostics)
		}
	}
	if err = os.WriteFile(file, []byte("<?php\necho 'before';\ncount(null);\necho 'after';"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostics bytes.Buffer
	code := run([]string{file}, &out, &diagnostics)
	want := "Fatal error: Uncaught TypeError: count(): Argument #1 ($value) must be of type Countable|array, null given in " + file + ":3\nStack trace:\n#0 {main}\n  thrown in " + file + " on line 3\n"
	if code != 255 || out.String() != "before" || diagnostics.String() != want {
		t.Fatalf("exit=%d out=%q err=%q", code, &out, &diagnostics)
	}
}
