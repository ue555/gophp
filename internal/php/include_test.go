package php

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIncludeUsesCallerScopeAndIndependentHistory(t *testing.T) {
	dir := t.TempDir()
	main := filepath.Join(dir, "main.php")
	child := filepath.Join(dir, "child.inc")
	if err := os.WriteFile(child, []byte("<?php\n$n++; return __DIR__;\n"), 0600); err != nil {
		t.Fatal(err)
	}
	canonicalDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	source := "<?php\n$n = 0; $a = include_once __DIR__ . '/child.inc'; $b = include_once __DIR__ . '/./child.inc'; echo $a, '|', $b, '|', $n;\n"
	for i := 0; i < 2; i++ {
		var out, diagnostics bytes.Buffer
		code := RunWithVersion(source, main, &out, &diagnostics, PHP84)
		if code != 0 || out.String() != canonicalDir+"|1|1" || diagnostics.Len() != 0 {
			t.Fatalf("run %d: exit=%d stdout=%q stderr=%q", i, code, out.String(), diagnostics.String())
		}
	}
}

func TestIncludeCompilesBeforeExecutingChild(t *testing.T) {
	dir := t.TempDir()
	main := filepath.Join(dir, "main.php")
	child := filepath.Join(dir, "bad.inc")
	if err := os.WriteFile(child, []byte("<?php\necho 'child';\n$x = ;\n"), 0600); err != nil {
		t.Fatal(err)
	}
	canonicalChild, err := filepath.EvalSymlinks(child)
	if err != nil {
		t.Fatal(err)
	}
	source := "<?php\necho 'parent';\ninclude __DIR__ . '/bad.inc';\necho 'after';\n"
	var out, diagnostics bytes.Buffer
	code := RunWithVersion(source, main, &out, &diagnostics, PHP84)
	if code != 255 || out.String() != "parent" || !strings.Contains(diagnostics.String(), "in "+canonicalChild+" on line 3") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestIncludeInFunctionKeepsLocalScope(t *testing.T) {
	dir := t.TempDir()
	main := filepath.Join(dir, "main.php")
	child := filepath.Join(dir, "local.inc")
	if err := os.WriteFile(child, []byte("<?php\n$x++; return $x;\n"), 0600); err != nil {
		t.Fatal(err)
	}
	source := "<?php\n$x=10; function f() { $x=1; return include __DIR__ . '/local.inc'; } echo f(), '|', $x;\n"
	var out, diagnostics bytes.Buffer
	code := RunWithVersion(source, main, &out, &diagnostics, PHP84)
	if code != 0 || out.String() != "2|10" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}
