package php

import (
	"bytes"
	"strings"
	"testing"
)

func TestErrorHandlerWarningsAndState(t *testing.T) {
	cases := []struct {
		name, source, out, errText string
		exit                       int
	}{
		{"handled variable", `<?php
function handle($level, $message, $file, $line) { echo $level, ':', $message, '|'; return true; }
set_error_handler('handle'); echo $missing; echo 'done';`, "2:Undefined variable $missing|done", "", 0},
		{"false falls through", `<?php
function handle($level, $message) { echo 'called|'; return false; }
set_error_handler('handle'); echo $missing;`, "called|", "Warning: Undefined variable $missing in fixture.php on line 3\n", 0},
		{"mask and restore", `<?php
function handle($level) { echo 'handled|'; return true; }
set_error_handler('handle', E_USER_NOTICE); trigger_error('one');
restore_error_handler(); trigger_error('two', E_USER_WARNING);`, "handled|", "Warning: two in fixture.php on line 4\n", 0},
		{"throw interrupts expression", `<?php
function handle($level) { throw new Exception('stop'); }
set_error_handler('handle');
try { echo $missing, 'bad'; } catch (Exception $e) { echo $e->getMessage(); }`, "stop", "", 0},
		{"independent execution", `<?php echo $missing;`, "", "Warning: Undefined variable $missing in fixture.php on line 1\n", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var out, diagnostics bytes.Buffer
			exit := RunWithVersion(tc.source, "fixture.php", &out, &diagnostics, PHP84)
			if exit != tc.exit || out.String() != tc.out || diagnostics.String() != tc.errText {
				t.Errorf("exit %d, stdout %q, stderr %q", exit, out.String(), diagnostics.String())
			}
		})
	}
}

func TestDiagnosticExceptionKeepsDynamicPropertyWrite(t *testing.T) {
	source := `<?php
class Box {}
function stop($level, $message) { throw new Exception('halt'); }
set_error_handler('stop');
$box = new Box;
try { $box->extra = 7; } catch (Exception $e) { echo isset($box->extra) ? 'set' : 'absent'; }
restore_error_handler();
echo '|done';`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		exit := RunWithVersion(source, "fixture.php", &out, &diagnostics, version)
		if exit != 0 || out.String() != "set|done" || diagnostics.Len() != 0 {
			t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, exit, out.String(), diagnostics.String())
		}
	}
}

func TestErrorHandlerNestedWarning(t *testing.T) {
	source := `<?php
function handle($level) { echo $inside; return true; }
set_error_handler('handle'); echo $outside;`
	var out, diagnostics bytes.Buffer
	if exit := RunWithVersion(source, "fixture.php", &out, &diagnostics, PHP84); exit != 0 {
		t.Fatalf("exit %d: %s", exit, diagnostics.String())
	}
	if !strings.Contains(diagnostics.String(), "Undefined variable $inside") || strings.Contains(diagnostics.String(), "$outside") {
		t.Fatalf("nested diagnostic = %q", diagnostics.String())
	}
}

func TestErrorHandlerReplacementDuringCall(t *testing.T) {
	source := `<?php
function replacement($level, $message) { echo "new:", $message, "|"; return true; }
function original($level, $message) {
    echo set_error_handler('replacement') === null ? "null|" : "wrong|";
    trigger_error('nested', E_USER_WARNING);
    return true;
}
set_error_handler('original');
trigger_error('outer', E_USER_WARNING);
trigger_error('later', E_USER_WARNING);`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		if exit := RunWithVersion(source, "fixture.php", &out, &diagnostics, version); exit != 0 || out.String() != "null|new:nested|new:later|" || diagnostics.Len() != 0 {
			t.Errorf("%s: exit=%d, stdout=%q, stderr=%q", version, exit, out.String(), diagnostics.String())
		}
	}
}

func TestErrorHandlerNullRegistrationDuringCall(t *testing.T) {
	source := `<?php
function original($level, $message) {
    echo $message, "|";
    echo set_error_handler(null) === null ? "null|" : "wrong|";
    return true;
}
set_error_handler('original');
trigger_error('first', E_USER_WARNING);
trigger_error('second', E_USER_WARNING);`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		if exit := RunWithVersion(source, "fixture.php", &out, &diagnostics, version); exit != 0 || out.String() != "first|null|second|null|" || diagnostics.Len() != 0 {
			t.Errorf("%s: exit=%d, stdout=%q, stderr=%q", version, exit, out.String(), diagnostics.String())
		}
	}
}

func TestErrorHandlerReceivesFullMessage(t *testing.T) {
	source := `<?php
function handle($level, $message) { echo strlen($message); return true; }
set_error_handler('handle');
trigger_error($message, E_USER_WARNING);`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		input := strings.Replace(source, "trigger_error", "$message = '"+strings.Repeat("x", 1100)+"'; trigger_error", 1)
		if exit := RunWithVersion(input, "fixture.php", &out, &diagnostics, version); exit != 0 || out.String() != "1100" || diagnostics.Len() != 0 {
			t.Errorf("%s: exit=%d, stdout=%q, stderr=%q", version, exit, out.String(), diagnostics.String())
		}
	}
}

func TestErrorHandlerVersionedConstantsAndFatal(t *testing.T) {
	for _, tc := range []struct {
		version Version
		all     string
		fatal   string
	}{
		{PHP83, "32767", "256|done"},
		{PHP84, "30719", "8192|256|done"},
		{PHP85, "30719", "8192|256|done"},
	} {
		var out, diagnostics bytes.Buffer
		source := `<?php
function handled($level) { echo $level, '|'; return true; }
echo E_ALL;`
		if exit := RunWithVersion(source, "fixture.php", &out, &diagnostics, tc.version); exit != 0 || out.String() != tc.all || diagnostics.Len() != 0 {
			t.Errorf("%s E_ALL: exit=%d out=%q err=%q", tc.version, exit, out.String(), diagnostics.String())
		}
		out.Reset()
		diagnostics.Reset()
		source = `<?php
function handled($level) { echo $level, '|'; return true; }
set_error_handler('handled');
trigger_error('fatal', E_USER_ERROR);
echo 'done';`
		if exit := RunWithVersion(source, "fixture.php", &out, &diagnostics, tc.version); exit != 0 || out.String() != tc.fatal || diagnostics.Len() != 0 {
			t.Errorf("%s user error: exit=%d out=%q err=%q", tc.version, exit, out.String(), diagnostics.String())
		}
	}
}
