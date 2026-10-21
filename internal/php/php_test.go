package php

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The four parser diagnostics use the PHP output supplied by the user.
// User-supplied official PHP comparison logs confirm all 81 fixtures (the
// original 57 and 24 logical-operator fixtures) on PHP 8.3.35, 8.4.26 and 8.5.11.
// A later user-supplied comparison ran all 213 fixtures. Its PHP output corrects
// eight diagnostic expectations (12 version/fixture differences); the other
// results matched. The subsequent user-supplied official PHP comparison confirms
// all 213 fixtures on all three targets after those fixes (639 matches, no differences).
// The next user-supplied comparison ran 337 fixtures: 336/333/336 matched on
// 8.3/8.4/8.5. Its PHP output corrects object modulo and PHP 8.4's three function
// redeclaration expectations. The subsequent user-supplied official PHP log
// confirms all 337 fixtures on all three targets (1011 matches, no differences).
// The 129 fixtures added for items 9-12 have source-derived, hand-written
// expectations. One fixture's helper collided with PHP's key() builtin and
// was renamed without changing its expected result. After starting the stopped
// containers, the user's full official PHP comparison confirms all 466 fixtures
// on PHP 8.3.35, 8.4.26 and 8.5.11 (1398 matches, no differences or errors); see
// docs/compatibility.md and docs/validation/v1-user-comparison-complete.md.
// PHP 8.5.11's additional NaN warning is captured as a version-specific stderr
// expectation. The default expectations retain PHP 8.4.26 behavior.
// The differential runner also checks them against an official CLI when available.
func TestFixtures(t *testing.T) {
	root, err := filepath.Abs("../../tests/fixtures")
	if err != nil {
		t.Fatal(err)
	}
	// The differential runner executes fixtures from the repository root.
	// Relative include_path entries and explicit ./ paths must see that CWD
	// during the saved-expectation test as well.
	priorCWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(filepath.Dir(filepath.Dir(root))); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(priorCWD) })
	data, err := os.ReadFile(filepath.Join(root, "expected.json"))
	if err != nil {
		t.Fatal(err)
	}
	var expected map[string]struct {
		Stdout          string            `json:"stdout"`
		Stderr          string            `json:"stderr"`
		Exit            int               `json:"exit"`
		StderrByVersion map[string]string `json:"stderr_by_version"`
		StdoutByVersion map[string]string `json:"stdout_by_version"`
		ExitByVersion   map[string]int    `json:"exit_by_version"`
	}
	if err := json.Unmarshal(data, &expected); err != nil {
		t.Fatal(err)
	}
	seen := 0
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".php" {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		want, ok := expected[rel]
		if !ok {
			t.Errorf("fixture has no expected result: %s", rel)
			return nil
		}
		seen++
		t.Run(rel, func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if strings.HasPrefix(rel, "success/") && (want.Exit != 0 || want.Stderr != "") {
				t.Fatal("success fixture must exit 0 without diagnostics")
			}
			for version := range want.StderrByVersion {
				if _, err := ParseVersion(version); err != nil {
					t.Fatal(err)
				}
			}
			for _, version := range []Version{"default", PHP83, PHP84, PHP85} {
				t.Run(string(version), func(t *testing.T) {
					var stdout, stderr bytes.Buffer
					var code int
					if version == "default" {
						code = Run(string(source), path, &stdout, &stderr)
					} else {
						code = RunWithVersion(string(source), path, &stdout, &stderr, version)
					}
					wantExit := want.Exit
					if override, ok := want.ExitByVersion[string(version)]; ok {
						wantExit = override
					}
					wantStdout := want.Stdout
					if override, ok := want.StdoutByVersion[string(version)]; ok {
						wantStdout = override
					}
					wantStdout = strings.ReplaceAll(wantStdout, "{file}", path)
					if code != wantExit {
						t.Errorf("exit: got %d, want %d", code, wantExit)
					}
					if stdout.String() != wantStdout {
						t.Errorf("stdout: got %q, want %q", stdout.String(), wantStdout)
					}
					wantStderr := want.Stderr
					if override, ok := want.StderrByVersion[string(version)]; ok {
						wantStderr = override
					}
					wantStderr = strings.ReplaceAll(wantStderr, "{file}", path)
					wantStderr = strings.ReplaceAll(wantStderr, "{fixtures}", root)
					traceFile := path
					if len(traceFile) > 15 {
						traceFile = traceFile[:15] + "..."
					}
					wantStderr = strings.ReplaceAll(wantStderr, "{trace_file}", traceFile)
					if strings.HasPrefix(rel, "success/") && wantExit != 0 {
						t.Fatal("success fixture must exit 0 in every version")
					}
					if stderr.String() != wantStderr {
						t.Errorf("stderr: got %q, want %q", stderr.String(), wantStderr)
					}
				})
			}
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if seen != len(expected) {
		t.Errorf("found %d fixtures for %d expectations", seen, len(expected))
	}
}

func TestRejectUnsupported(t *testing.T) {
	// Keep the original inputs: newly supported forms now assert their result.
	newlySupported := map[string]string{"$a = &$b;": "", "foreach ([] as &$x) {}": "", "echo 2 ** 3;": "8", "echo 1 & 2;": "0", "echo 1 << 2;": "4", "echo true & false;": "0", "echo true | false;": "1", "echo true and false;": "", "echo true or false;": "1", "echo true xor false;": "1", "echo \"$name\";": "", "echo \"{$name}\";": "", "function f(): int {}": "", "function f(...$x) {}": "", "function f(&$x) {}": "", "echo 0x10;": "16", "echo $argc;": "1", "echo 1_000;": "1000", "#[Attribute] class A {}": "", "echo 0b10;": "2", "echo 0o10;": "8", "echo \"\\u{41}\";": "A", "echo 1 ?>trailing HTML": "1trailing HTML", "echo 1; ?><?php echo 2;": "12"}
	cases := []string{
		"$a = &$b;", "foreach ([] as &$x) {}", "function f(&$x) {}", "class A extends B {}",
		"function f(...$x) {}", "function f(): int {}", "echo 2 ** 3;",
		"$a = 1; $a **= 2;", "echo 1 & 2;", "echo 0x10;",
		"echo 1_000;", "echo \"$name\";", "echo \"${name}\";", "echo \"{$name}\";",
		"echo \"\\u{41}\";", "echo \"\\777\";", "#[Attribute] class A {}",
		"echo 1 ?>trailing HTML", "echo 1; ?><?php echo 2;",
		"echo 1e;", "echo 1e+;", "echo 0b10;", "echo 0o10;", "echo `ls`;",
		"echo $argc;",
		"echo 1 << 2;", "echo true and false;", "echo true or false;",
		"echo true xor false;", "echo true & false;", "echo true | false;",
	}
	for _, source := range cases {
		t.Run(source, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run("<?php\necho 'must not execute';\n"+source, "input.php", &stdout, &stderr)
			if want, supported := newlySupported[source]; supported {
				wantStderr := ""
				if strings.Contains(source, "$name") {
					wantStderr = "Warning: Undefined variable $name in input.php on line 3\n"
				}
				if code != 0 || stdout.String() != "must not execute"+want || stderr.String() != wantStderr {
					t.Fatalf("supported form: exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
				}
				return
			}
			if code != 255 || stdout.Len() != 0 || stderr.Len() == 0 {
				t.Fatalf("must reject before evaluation: exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
		})
	}
}

func TestMalformedIfBeforeExecution(t *testing.T) {
	for _, source := range []string{
		"if () {}", "if true {}", "if (true { }", "if (true) {",
		"if (true) {} else {", "if (true) { echo 1 }", "if (true) {} }",
		"else {}", "if (true) {} else {} else {}",
		"if (false) { echo ; }", "if (true) {} else { echo ; }",
		"if (false) { while () {} }", "if (true) {} else { $a = [1, ; }",
		"elseif (true) {}", "if (true) {} elseif () {}",
		"if (true) {} elseif (false) {", "if (true) {} else if () {}",
		"if (true) {} else {} elseif (true) {}",
		"if (true) {} elseif (false) { echo ; }",
		"if (true) {} else if (false) { echo ; }",
	} {
		t.Run(source, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run("<?php echo 'must not execute'; "+source, "input.php", &stdout, &stderr)
			if code != 255 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "Parse error:") {
				t.Fatalf("must parse both branches before executing: exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
		})
	}
}

func TestStrictComparisonGrouping(t *testing.T) {
	for _, first := range []string{"===", "!=="} {
		for _, second := range []string{"===", "!=="} {
			source := "<?php echo 'must not execute'; echo 1 " + first + " 1 " + second + " true;"
			var stdout, stderr bytes.Buffer
			if code := Run(source, "input.php", &stdout, &stderr); code != 255 || stdout.Len() != 0 {
				t.Fatalf("chained comparisons must be rejected: %s; exit=%d stdout=%q stderr=%q", source, code, &stdout, &stderr)
			}
		}
	}
	for _, source := range []string{
		"<?php echo ((1 === 1) === true) + 0, (true === (1 !== 2)) + 0;",
		"<?php echo ((1 !== 1) !== true) + 0, (false !== (1 === 1)) + 0;",
	} {
		var stdout, stderr bytes.Buffer
		if code := Run(source, "input.php", &stdout, &stderr); code != 0 || stdout.String() != "11" || stderr.Len() != 0 {
			t.Fatalf("explicit grouping must succeed: %s; exit=%d stdout=%q stderr=%q", source, code, &stdout, &stderr)
		}
	}
}

// Exercise changed content, grammar context and line numbers, in addition to
// the four fixtures captured in the user's official PHP comparison.
func TestParserDiagnosticVariants(t *testing.T) {
	cases := []struct {
		name, source, message string
		line                  int
	}{
		{"echo comma", "<?php\necho 1,\n;", `syntax error, unexpected token ";"`, 3},
		{"echo list EOF", "<?php\necho 1, 2", `syntax error, unexpected end of file, expecting "," or ";"`, 2},
		{"assignment EOF", "<?php\n$x = 1", `syntax error, unexpected end of file, expecting ";"`, 2},
		{"missing operand", "<?php\necho 1 +", `syntax error, unexpected end of file`, 2},
		{"optional array element", "<?php\necho 'bad';\nif (false) $a = [2, ;", `syntax error, unexpected token ";", expecting "]"`, 3},
		{"long array delimiter", "<?php\necho 'bad';\n$a = array(2, ;", `syntax error, unexpected token ";", expecting ")"`, 3},
		{"required array operand", "<?php\necho 'bad';\n$a = [2 + ;", `syntax error, unexpected token ";"`, 3},
		{"parse before compile error", "<?php\necho 'bad'; $a[]; echo ;", `syntax error, unexpected token ";"`, 2},
		{"comment line", "<?php\r\n\r\n/* open\r\ncomment", `Unterminated comment starting line 3`, 3},
		{"raw string", "<?php\n\necho 'alpha\\nbeta;\nmore", `syntax error, unexpected string content "alpha\nbeta;"`, 3},
		{"33 byte string", "<?php echo '" + strings.Repeat("x", 33), `syntax error, unexpected string content "` + strings.Repeat("x", 33) + `"`, 1},
		{"34 byte string", "<?php echo '" + strings.Repeat("x", 34), `syntax error, unexpected string content "` + strings.Repeat("x", 30) + `..."`, 1},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(test.source, "input.php", &stdout, &stderr)
			want := fmt.Sprintf("Parse error: %s in input.php on line %d\n", test.message, test.line)
			if code != 255 || stdout.Len() != 0 || stderr.String() != want {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want stderr=%q", code, &stdout, &stderr, want)
			}
		})
	}
}

func TestValueTypesAndBoundaries(t *testing.T) {
	cases := []struct {
		source string
		want   value
	}{
		{"6 / 3", integer(2)},
		{"7 / 2", floating(3.5)},
		{"6.0 / 3", floating(2)},
		{"'6' / '3'", integer(2)},
		{"'6.0' / '3'", floating(2)},
		{"9223372036854775806 + 1", integer(math.MaxInt64)},
		{"-9223372036854775807 - 1", integer(math.MinInt64)},
		{"9223372036854775807 + 1", floating(9223372036854775808)},
		{"(-9223372036854775807 - 1) / -1", floating(9223372036854775808)},
		{"(-9223372036854775807 - 1) * -1", floating(9223372036854775808)},
		{"-1 * (-9223372036854775807 - 1)", floating(9223372036854775808)},
		{"(-9223372036854775807 - 1) * 1", integer(math.MinInt64)},
		{"3037000499 * 3037000499", integer(9223372030926249001)},
		{"3037000500 * 3037000500", floating(float64(3037000500) * float64(3037000500))},
		{"9223372036854775808", floating(9223372036854775808)},
		{"+null", integer(0)},
		{"+true", integer(1)},
		{"!0", boolean(true)},
		{"!-1", boolean(false)},
		{"!!'0'", boolean(false)},
		{"true && 'text'", boolean(true)},
		{"false || 'text'", boolean(true)},
		{"true && '0'", boolean(false)},
		{"false && (1 / 0)", boolean(false)},
		{"true || (1 / 0)", boolean(true)},
		{"1 === 1 && 2 === 2", boolean(true)},
	}
	for _, test := range cases {
		t.Run(test.source, func(t *testing.T) {
			tokens, err := lex("<?php " + test.source + ";")
			if err != nil {
				t.Fatal(err)
			}
			statements, err := parse(tokens)
			if err != nil {
				t.Fatal(err)
			}
			e := evaluator{stderr: io.Discard}
			got, err := e.eval(statements[0].expressions[0])
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("got %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestRunHasFreshVariables(t *testing.T) {
	Run("<?php $x = 1;", "input.php", io.Discard, io.Discard)
	var stderr bytes.Buffer
	if code := Run("<?php echo $x;", "input.php", io.Discard, &stderr); code != 0 || !strings.Contains(stderr.String(), "Undefined variable $x") {
		t.Fatalf("variables leaked: exit=%d stderr=%q", code, &stderr)
	}
}

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) { return 0, io.ErrClosedPipe }
func TestOutputError(t *testing.T) {
	var stderr bytes.Buffer
	for _, source := range []string{"<?php echo 'x';", "<?php if (true) { if (false) {} else { echo 'x'; } }"} {
		stderr.Reset()
		if code := Run(source, "input.php", failingWriter{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "write stdout") {
			t.Fatalf("write failure must not succeed: exit=%d stderr=%q", code, &stderr)
		}
	}
}

func FuzzRun(f *testing.F) {
	for _, seed := range []string{"", "<?php ", "<?php echo 'x';", "<?php echo 1/0;", "<?php /*", "<?php $x=9223372036854775807+1; echo $x;", "<?php echo \"\\x41\";"} {
		f.Add(seed)
	}
	for _, seed := range []string{"<?php if (1 === 1) { echo 'yes'; } else { echo 1/0; }", "<?php if ('0') { if ($missing) {} } else { echo 1 !== '1'; }", "<?php echo 1 === 1 !== false;", "<?php if (true) {"} {
		f.Add(seed)
	}
	for _, seed := range []string{
		"<?php echo !!!'0' || (true && !$missing);",
		"<?php echo false && (1 / 0), true || $missing;",
		"<?php $n = (1e308 * 10) - (1e308 * 10); echo !$n, $n && $n, $n || $n;",
		"<?php echo true || ;", "<?php echo 1 === 1 && 2 === 2;",
		"<?php echo 0 == 'word', null < '0', 1 < 2 == true;",
		"<?php if (0 > 1) {} elseif (0 == '0') {} else if (true) {} else {}",
		"<?php echo 1 < 2 < 3;", "<?php if (true) {} elseif (1 < ) {}",
	} {
		f.Add(seed)
	}
	for _, seed := range []string{
		"<?php $a = []; $a[] = $a; echo $a === [[]];",
		"<?php $o = new stdClass; $o->self = $o; echo $o == $o;",
		"<?php $a['x'][] = ($b = 2); if ($b) echo $a['x'][0]; else echo 1/0;",
		"<?php echo false && ($a[] = 1/0); echo [1] <=> ['1'];",
	} {
		f.Add(seed)
	}
	for _, seed := range []string{
		"<?php class C { public $a = []; function f($v = 1) { $this->a[] = $v; return $this; } } $c = new C; echo $c->f(2)->a[0];",
		"<?php class C { function f() { return $this->f(); } } $c = new C; $c->f();",
		"<?php class C { private $x; function f($o) { return $o->x; } }",
		"<?php if (false) { class C { function f() { echo ; } } }",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, source string) {
		if len(source) > 4096 {
			t.Skip()
		}
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			code := run(source, "fuzz.php", io.Discard, io.Discard, version, 1000)
			if code != 0 && code != 255 {
				t.Fatalf("unexpected exit code %d for PHP %s", code, version)
			}
		}
	})
}
