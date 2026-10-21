package php

import (
	"bytes"
	"testing"
)

func TestDynamicPropertyTargetEvaluatedOncePerOperation(t *testing.T) {
	source := `<?php
class Box { public int $value = 1; }
$box = new Box;
$calls = 0;
function member() { global $calls; $calls++; return 'value'; }
echo ($box->{member()} .= '2'), '|', $calls, '|';
echo $box->{member()}++, '|', $calls, '|', $box->value;
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for i := 0; i < 2; i++ {
			var out, diagnostics bytes.Buffer
			code := RunWithVersion(source, "dynamic.php", &out, &diagnostics, version)
			if code != 0 || out.String() != "12|1|12|2|13" || diagnostics.Len() != 0 {
				t.Fatalf("PHP %s run %d: exit=%d stdout=%q stderr=%q", version, i, code, &out, &diagnostics)
			}
		}
	}
}

func TestDynamicPropertyAssignmentUsesReboundVariable(t *testing.T) {
	source := `<?php
class Box { public int $value; }
$first = new Box; $first->value = 1;
$second = new Box; $second->value = 2;
$selected = $first;
function member() { global $selected, $second; $selected = $second; return 'value'; }
$selected->{member()} = 9;
echo $first->value, '|', $second->value, '|';
$selected = $first;
function rhs() { global $selected, $second; $selected = $second; return 10; }
$selected->{'value'} = rhs();
echo $first->value, '|', $second->value;
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		code := RunWithVersion(source, "dynamic.php", &out, &diagnostics, version)
		if code != 0 || out.String() != "1|9|1|10" || diagnostics.Len() != 0 {
			t.Fatalf("PHP %s: exit=%d stdout=%q stderr=%q", version, code, &out, &diagnostics)
		}
	}
}
