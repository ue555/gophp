package php

import (
	"bytes"
	"fmt"
	"testing"
)

// The supplied official log covers an object on the left. Exercise the shared
// conversion on either side, including compound assignments and warning order.
func TestModuloRejectsObjects(t *testing.T) {
	const precisionWarning = "Deprecated: Implicit conversion from float 2.5 to int loses precision in input.php on line 2\n"
	cases := []struct {
		name, source, types, prefix string
	}{
		{"left", "$o = new stdClass; echo $o % 3;", "stdClass % int", ""},
		{"right", "$o = new stdClass; echo 3 % $o;", "int % stdClass", ""},
		{"both", "$o = new stdClass; echo $o % $o;", "stdClass % stdClass", ""},
		{"zero", "$o = new stdClass; echo $o % 0;", "stdClass % int", ""},
		{"compound_left", "$o = new stdClass; $o %= 3;", "stdClass % int", ""},
		{"compound_right", "$o = new stdClass; $n = 3; $n %= $o;", "int % stdClass", ""},
		{"class", "class C {} $o = new C; echo $o % 3;", "C % int", ""},
		{"array_element", "$a = [new stdClass]; $a[0] %= 3;", "stdClass % int", ""},
		{"skip_rhs_conversion", "$o = new stdClass; echo $o % 2.5;", "stdClass % float", ""},
		{"lhs_conversion_warning", "$o = new stdClass; echo 2.5 % $o;", "float % stdClass", precisionWarning},
	}
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for _, tc := range cases {
			t.Run(string(version)+"/"+tc.name, func(t *testing.T) {
				var stdout, stderr bytes.Buffer
				code := RunWithVersion("<?php\n"+tc.source+" echo 'must not execute';", "input.php", &stdout, &stderr, version)
				want := tc.prefix + fmt.Sprintf("Fatal error: Uncaught TypeError: Unsupported operand types: %s in input.php:2\nStack trace:\n#0 {main}\n  thrown in input.php on line 2\n", tc.types)
				if code != 255 || stdout.Len() != 0 || stderr.String() != want {
					t.Fatalf("exit=%d stdout=%q stderr=%q; want stderr=%q", code, &stdout, &stderr, want)
				}
			})
		}
	}
}

func TestModuloHandlerExceptionPreservesTargets(t *testing.T) {
	const source = `<?php
class Box { public float $value = 5.5; }
$number = 5.5; $array = [5.5]; $box = new Box;
set_error_handler(function ($level, $message) { throw new Exception('stop'); });
try { $number %= 2; } catch (Exception $e) { echo 'caught|'; }
try { $array[0] %= 2; } catch (Exception $e) { echo 'caught|'; }
try { $box->value %= 2; } catch (Exception $e) { echo 'caught|'; }
restore_error_handler();
echo gettype($number), ':', $number, '|', gettype($array[0]), ':', $array[0], '|', gettype($box->value), ':', $box->value;`
	const want = "caught|caught|caught|double:5.5|double:5.5|double:0"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		t.Run(string(version), func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := RunWithVersion(source, "input.php", &stdout, &stderr, version)
			if code != 0 || stdout.String() != want || stderr.Len() != 0 {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
		})
	}
}
