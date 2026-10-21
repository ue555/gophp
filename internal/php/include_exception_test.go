package php

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIncludedParseErrorIsCatchableButMainParseErrorIsNot(t *testing.T) {
	dir := t.TempDir()
	child := filepath.Join(dir, "bad.inc")
	if err := os.WriteFile(child, []byte("<?php\n$x = ;\n"), 0600); err != nil {
		t.Fatal(err)
	}
	main := filepath.Join(dir, "main.php")
	canonicalChild, err := filepath.EvalSymlinks(child)
	if err != nil {
		t.Fatal(err)
	}
	var out, diagnostics bytes.Buffer
	source := "<?php echo 'before|'; try { include __DIR__ . '/bad.inc'; } catch (ParseError $e) { echo 'caught:', $e->getFile(), ':', $e->getLine(); } echo '|after';"
	if code := Run(source, main, &out, &diagnostics); code != 0 || out.String() != "before|caught:"+canonicalChild+":2|after" || diagnostics.Len() != 0 {
		t.Fatalf("included exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
	out.Reset()
	diagnostics.Reset()
	if code := Run("<?php try { $x = ; } catch (ParseError $e) { echo 'caught'; }", main, &out, &diagnostics); code != 255 || out.Len() != 0 || !strings.HasPrefix(diagnostics.String(), "Parse error:") {
		t.Fatalf("main exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestIncludeOnceHistoryAfterThrowAndParseFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "throw.inc"), []byte("<?php\n$n++; throw new Exception('once');\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bad.inc"), []byte("<?php\n$x = ;\n"), 0600); err != nil {
		t.Fatal(err)
	}
	source := `<?php
$n = 0;
try { include_once __DIR__ . '/throw.inc'; } catch (Exception $e) { echo $e->getMessage(); }
echo '|', include_once __DIR__ . '/./throw.inc', '|', $n;
try { include_once __DIR__ . '/bad.inc'; } catch (ParseError $e) { echo '|parse'; }
echo '|', include_once __DIR__ . '/./bad.inc';
try { include __DIR__ . '/bad.inc'; } catch (ParseError $e) { echo '|normal'; }
`
	for i := 0; i < 2; i++ {
		var out, diagnostics bytes.Buffer
		if code := Run(source, filepath.Join(dir, "main.php"), &out, &diagnostics); code != 0 || out.String() != "once|1|1|parse|1|normal" || diagnostics.Len() != 0 {
			t.Fatalf("run %d exit=%d stdout=%q stderr=%q", i, code, out.String(), diagnostics.String())
		}
	}
}

func TestParseFailureHistoryDependsOnIncludeForm(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad.inc"), []byte("<?php\n$x = ;\n"), 0600); err != nil {
		t.Fatal(err)
	}
	main := filepath.Join(dir, "main.php")
	cases := []struct {
		name, first, second, want string
	}{
		{"include then include_once", "include", "include_once", "first|second|1"},
		{"require then require_once", "require", "require_once", "first|second|1"},
		{"include_once then include_once", "include_once", "include_once", "first|1|1"},
		{"require_once then include_once", "require_once", "include_once", "first|1|1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source := "<?php try { " + tc.first + " __DIR__ . '/bad.inc'; } catch (ParseError $e) { echo 'first|'; } " +
				"try { echo " + tc.second + " __DIR__ . '/bad.inc', '|'; } catch (ParseError $e) { echo 'second|'; } " +
				"echo " + tc.second + " __DIR__ . '/bad.inc';"
			// Run twice to ensure the include history belongs to one execution.
			for i := 0; i < 2; i++ {
				var out, diagnostics bytes.Buffer
				if code := Run(source, main, &out, &diagnostics); code != 0 || out.String() != tc.want || diagnostics.Len() != 0 {
					t.Fatalf("run %d exit=%d stdout=%q stderr=%q", i, code, out.String(), diagnostics.String())
				}
			}
		})
	}
}

func TestIncludeFailureAndExceptionTrace(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "throw.inc"), []byte("<?php\nthrow new Exception('inside');\n"), 0600); err != nil {
		t.Fatal(err)
	}
	main := filepath.Join(dir, "main.php")
	canonicalThrow, err := filepath.EvalSymlinks(filepath.Join(dir, "throw.inc"))
	if err != nil {
		t.Fatal(err)
	}
	var out, diagnostics bytes.Buffer
	source := "<?php\ntry { require __DIR__ . '/throw.inc'; } catch (Exception $e) { echo $e->getFile(), '|', $e->getLine(), '|', $e->getTraceAsString(); }"
	if code := Run(source, main, &out, &diagnostics); code != 0 || !strings.Contains(out.String(), canonicalThrow+"|2|#0 "+main+"(2): require()") || diagnostics.Len() != 0 {
		t.Fatalf("trace exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
	out.Reset()
	diagnostics.Reset()
	source = "<?php\ntry { require __DIR__ . '/absent.inc'; } catch (Error $e) { echo 'caught'; } finally { echo '|finally'; }"
	if code := Run(source, main, &out, &diagnostics); code != 0 || out.String() != "caught|finally" || !strings.Contains(diagnostics.String(), "Failed to open stream") {
		t.Fatalf("missing exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestIncludedFatalDeclarationCannotBeReplacedByFinallyReturn(t *testing.T) {
	dir := t.TempDir()
	child := filepath.Join(dir, "duplicate.inc")
	if err := os.WriteFile(child, []byte("<?php\nfunction repeated() {}\nfunction repeated() {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	main := filepath.Join(dir, "main.php")
	source := `<?php
function load() {
    echo 'before|';
    try { require __DIR__ . '/duplicate.inc'; }
    finally { echo 'finally|'; return 9; }
}
echo load(), '|after';
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		code := RunWithVersion(source, main, &out, &diagnostics, version)
		if code != 255 || out.String() != "before|" || !strings.Contains(diagnostics.String(), "Cannot redeclare") || strings.Contains(diagnostics.String(), ": require()") {
			t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
		}
		if version == PHP85 && !strings.Contains(diagnostics.String(), ": load()") {
			t.Fatalf("%s: missing caller frame in %q", version, diagnostics.String())
		}
	}
}
