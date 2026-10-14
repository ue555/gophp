package php

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExceptionControlFlow(t *testing.T) {
	cases := []struct{ source, want string }{
		{`try { throw new Exception('bad'); } catch (Exception $e) { echo $e->getMessage(); } finally { echo '!'; }`, "bad!"},
		{`try { throw new TypeError('bad'); } catch (Error|Exception $e) { echo $e->getMessage(); }`, "bad"},
		{`function f() { try { return 1; } finally { echo 'f'; } } echo f();`, "f1"},
		{`function f() { try { return 1; } finally { return 2; } } echo f();`, "2"},
		{`try { strlen([]); } catch (TypeError $e) { echo $e->getMessage(); }`, "strlen(): Argument #1 ($string) must be of type string, array given"},
		{`function required($a) {} try { required(); } catch (ArgumentCountError $e) { echo 'caught'; }`, "caught"},
		{`if (false) { throw new Exception('bad'); } echo 'ok';`, "ok"},
		{`$x = null; try { $x ?? throw new Exception('missing'); } catch (Exception $e) { echo $e->getMessage(); }`, "missing"},
		{`$i = 0; while ($i < 2) { try { $i++; continue; } finally { echo $i; } }`, "12"},
		{`try { echo 'ok'; } finally {}`, "ok"},
		{`$e = null; function g() { global $e; try { throw new Exception('x'); } catch (Exception $e) {} } g(); echo gettype($e);`, "object"},
		{`function s() { static $e = null; try { throw new Exception('x'); } catch (Exception $e) {} echo gettype($e); } s(); s();`, "objectobject"},
	}
	for _, tc := range cases {
		var out, diagnostic bytes.Buffer
		code := Run("<?php "+tc.source, "exception.php", &out, &diagnostic)
		if code != 0 || out.String() != tc.want || diagnostic.Len() != 0 {
			t.Errorf("source=%s exit=%d stdout=%q stderr=%q", tc.source, code, out.String(), diagnostic.String())
		}
	}
}

func TestExceptionAcrossInclude(t *testing.T) {
	directory := t.TempDir()
	child := filepath.Join(directory, "child.inc")
	if err := os.WriteFile(child, []byte("<?php throw new Exception('inside');"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostic bytes.Buffer
	main := filepath.Join(directory, "main.php")
	source := "<?php try { require __DIR__ . '/child.inc'; } catch (Exception $e) { echo $e->getMessage(); }"
	code := Run(source, main, &out, &diagnostic)
	if code != 0 || out.String() != "inside" || diagnostic.Len() != 0 {
		t.Errorf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostic.String())
	}
}

func TestUncaughtException(t *testing.T) {
	var out, diagnostic bytes.Buffer
	code := Run("<?php echo 'before'; throw new Exception('bad');", "exception.php", &out, &diagnostic)
	if code != 255 || out.String() != "before" || !strings.Contains(diagnostic.String(), "Uncaught Exception: bad") {
		t.Errorf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostic.String())
	}
}

func TestUserClassCannotImplementThrowableDirectly(t *testing.T) {
	var out, diagnostic bytes.Buffer
	code := Run("<?php class Direct implements Throwable {}", "throwable.php", &out, &diagnostic)
	if code != 255 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "cannot implement interface Throwable") {
		t.Errorf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostic.String())
	}
}
