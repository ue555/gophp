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

func TestExceptionChainPreservesIdentityAndAvoidsCycles(t *testing.T) {
	first := value{kind: objectKind, o: &phpObject{exception: &exceptionData{message: "first"}}}
	second := value{kind: objectKind, o: &phpObject{exception: &exceptionData{message: "second", previous: first}}}
	third := value{kind: objectKind, o: &phpObject{exception: &exceptionData{message: "third"}}}
	chainException(third, second)
	if third.o.exception.previous.o != second.o || second.o.exception.previous.o != first.o {
		t.Fatal("chain did not retain the original throwable objects")
	}
	chainException(first, third)
	if first.o.exception.previous.kind != nullKind {
		t.Fatal("cycle was accepted")
	}
}

func TestVariadicTraceKeepsCallArguments(t *testing.T) {
	cases := []string{
		`$xs[0] = 99;`,
		`$xs[] = 30;`,
		`unset($xs[0]);`,
	}
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for _, change := range cases {
			source := `<?php function makeE(...$xs) { ` + change + ` return new Exception; } ` +
				`$args = makeE(10, 20)->getTrace()[0]['args']; ` +
				`echo count($args), '|', $args[0], '|', $args[1];`
			var out, diagnostic bytes.Buffer
			code := RunWithVersion(source, "variadic_trace.php", &out, &diagnostic, version)
			if code != 0 || out.String() != "2|10|20" || diagnostic.Len() != 0 {
				t.Errorf("version=%s change=%s exit=%d stdout=%q stderr=%q", version, change, code, out.String(), diagnostic.String())
			}
		}
	}
}

func TestClosureTraceNameByVersion(t *testing.T) {
	source := "<?php\n$make = function () { return new Exception; };\n$e = $make();\necho $e->getTrace()[0]['function'];\n"
	for _, tc := range []struct {
		version Version
		want    string
	}{
		{PHP83, "{closure}"},
		{PHP84, "{closure:closure_trace.php:2}"},
		{PHP85, "{closure:closure_trace.php:2}"},
	} {
		var out, diagnostic bytes.Buffer
		code := RunWithVersion(source, "closure_trace.php", &out, &diagnostic, tc.version)
		if code != 0 || out.String() != tc.want || diagnostic.Len() != 0 {
			t.Errorf("version=%s exit=%d stdout=%q stderr=%q; want %q", tc.version, code, out.String(), diagnostic.String(), tc.want)
		}
	}
}
