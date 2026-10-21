package php

import (
	"bytes"
	"strings"
	"testing"
)

func TestVariableReferenceIsolatedRuns(t *testing.T) {
	source := "<?php\n$a = 1; $b =& $a; $b = 2; echo $a;\n"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for run := 0; run < 2; run++ {
			var out, diagnostic bytes.Buffer
			if status := RunWithVersion(source, "reference.php", &out, &diagnostic, version); status != 0 || out.String() != "2" || diagnostic.Len() != 0 {
				t.Fatalf("%s run %d: exit=%d stdout=%q stderr=%q", version, run, status, out.String(), diagnostic.String())
			}
		}
	}
}

func TestArrayReferenceInUnselectedBranch(t *testing.T) {
	source := "<?php\necho 'before'; if (false) { $alias =& $array[0]; }\n"
	var out, diagnostic bytes.Buffer
	status := Run(source, "reference.php", &out, &diagnostic)
	if status != 0 || out.String() != "before" || diagnostic.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", status, out.String(), diagnostic.String())
	}
}

func TestUnsupportedReferenceTargetRejectedBeforeExecution(t *testing.T) {
	source := "<?php\necho 'before'; if (false) { $alias =& 1; }\n"
	var out, diagnostic bytes.Buffer
	status := Run(source, "reference.php", &out, &diagnostic)
	if status != 255 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "unsupported reference assignment target") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", status, out.String(), diagnostic.String())
	}
}

func TestCyclicArrayReferenceFailsExplicitly(t *testing.T) {
	source := "<?php\n$a = []; $a[0] =& $a;\n"
	var out, diagnostic bytes.Buffer
	status := Run(source, "reference.php", &out, &diagnostic)
	if status != 255 || !strings.Contains(diagnostic.String(), "cyclic array references are not supported") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", status, out.String(), diagnostic.String())
	}
}

func TestTypedPropertyReferenceEnforcesType(t *testing.T) {
	source := "<?php\nclass Guard { public int $value = 1; public static int $saved = 2; }\n$box = new Guard;\n$alias =& $box->value; $alias = '3';\ntry { $alias = []; } catch (TypeError $e) { echo 'instance|'; }\n$alias =& Guard::$saved; $alias = '4';\ntry { $alias = []; } catch (TypeError $e) { echo 'static|'; }\necho $box->value, '|', Guard::$saved;\n"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostic bytes.Buffer
		status := RunWithVersion(source, "reference.php", &out, &diagnostic, version)
		if status != 0 || out.String() != "instance|static|3|4" || diagnostic.Len() != 0 {
			t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, status, out.String(), diagnostic.String())
		}
	}
}

func TestNullsafeReferenceRejectedBeforeExecution(t *testing.T) {
	source := "<?php\necho 'before'; if (false) { $alias =& $box?->value; }\n"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostic bytes.Buffer
		status := RunWithVersion(source, "reference.php", &out, &diagnostic, version)
		if status != 255 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "Can't use nullsafe operator in write context") {
			t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, status, out.String(), diagnostic.String())
		}
	}
}
