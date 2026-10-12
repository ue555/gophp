package php

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
)

func TestArrayKeyConversions(t *testing.T) {
	for _, tc := range []struct {
		key  value
		want arrayKey
	}{
		{integer(math.MaxInt64), arrayKey{index: math.MaxInt64}},
		{integer(math.MinInt64), arrayKey{index: math.MinInt64}},
		{floating(-1), arrayKey{index: -1}},
		{floating(math.Copysign(0, -1)), arrayKey{index: 0}},
		{value{kind: stringKind, s: "-9223372036854775808"}, arrayKey{index: math.MinInt64}},
		{value{kind: stringKind, s: "9223372036854775808"}, arrayKey{isString: true, text: "9223372036854775808"}},
		{value{kind: stringKind, s: "-0"}, arrayKey{isString: true, text: "-0"}},
		{value{kind: stringKind, s: "01"}, arrayKey{isString: true, text: "01"}},
		{value{kind: stringKind, s: "1\x00"}, arrayKey{isString: true, text: "1\x00"}},
	} {
		e := evaluator{version: PHP84, stderr: io.Discard}
		got, err := e.arrayKey(tc.key, 1)
		if err != nil || got != tc.want {
			t.Fatalf("key %#v: got %#v (%v), want %#v", tc.key, got, err, tc.want)
		}
	}
}

func TestAssignmentAndAggregateTypes(t *testing.T) {
	for _, tc := range []struct {
		source string
		want   value
	}{
		{"1 <=> 2", integer(-1)},
		{"2 <=> 2", integer(0)},
		{"3 <=> 2", integer(1)},
		{"1 <> '1'", boolean(false)},
		{"[1] == ['1']", boolean(true)},
		{"[1] === ['1']", boolean(false)},
		{"$a = $b = 1.5", floating(1.5)},
		{"false || $a = 7", boolean(true)},
		{"$a = false || 7", boolean(true)},
		{"$a['x'][] = 8", integer(8)},
		{"![]", boolean(true)},
		{"!new stdClass", boolean(false)},
	} {
		t.Run(tc.source, func(t *testing.T) {
			tokens, err := lex("<?php " + tc.source + ";")
			if err != nil {
				t.Fatal(err)
			}
			statements, err := parse(tokens)
			if err != nil {
				t.Fatal(err)
			}
			e := evaluator{version: PHP84, stdout: io.Discard, stderr: io.Discard}
			got, err := e.eval(statements[0].expressions[0])
			if err != nil || got != tc.want {
				t.Fatalf("got %#v (%v), want %#v", got, err, tc.want)
			}
		})
	}
}

func TestAppendReadRejectedBeforeExecution(t *testing.T) {
	for _, code := range []string{
		"if (false) echo $a[];", "true || $a[];", "$b = [$a[]];",
		"if (false) $a[$b[]] = 1;", "if (true) ; else echo $a[][0];",
	} {
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			var out, err bytes.Buffer
			exit := RunWithVersion("<?php echo 'must not execute'; "+code, "input.php", &out, &err, version)
			if exit != 255 || out.Len() != 0 || !strings.Contains(err.String(), "Cannot use [] for reading") {
				t.Fatalf("PHP %s %s: exit=%d out=%q err=%q", version, code, exit, &out, &err)
			}
		}
	}
}

func TestNewDereferenceGrammar(t *testing.T) {
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for _, form := range []string{"(new stdClass)", "(new stdClass())", "new stdClass", "new stdClass()"} {
			var out, err bytes.Buffer
			exit := RunWithVersion("<?php echo "+form+"->x;", "input.php", &out, &err, version)
			invalid := form == "new stdClass" || version == PHP83 && form == "new stdClass()"
			if invalid && (exit != 255 || !strings.Contains(err.String(), "Parse error:")) || !invalid && (exit != 0 || !strings.Contains(err.String(), "Undefined property:")) {
				t.Fatalf("PHP %s %s: exit=%d stderr=%q", version, form, exit, &err)
			}
		}
	}
}

func TestScalarOffsetDiagnostics(t *testing.T) {
	// Zend uses the value name (true/false), rather than the bool type name.
	for _, tc := range []struct{ source, name string }{
		{"null", "null"}, {"true", "true"}, {"false", "false"}, {"42", "int"}, {"1.5", "float"},
	} {
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			var out, err bytes.Buffer
			source := "<?php\n$a = " + tc.source + ";\necho ($a[0] === null) + 0;"
			code := RunWithVersion(source, "scalar.php", &out, &err, version)
			want := "Warning: Trying to access array offset on " + tc.name + " in scalar.php on line 3\n"
			if code != 0 || out.String() != "1" || err.String() != want {
				t.Fatalf("PHP %s %s: exit=%d out=%q err=%q", version, tc.source, code, &out, &err)
			}
		}
	}
}

func TestNewDereferenceErrorContext(t *testing.T) {
	for _, tc := range []struct{ source, expected string }{
		{"echo new stdClass()->x;", `"," or ";"`},
		{"$a = new stdClass()->x;", `";"`},
		{"echo (new stdClass()->x);", `")"`},
		{"if (new stdClass()->x) echo 1;", `")"`},
	} {
		var out, err bytes.Buffer
		code := RunWithVersion("<?php\necho 'bad';\n"+tc.source, "context.php", &out, &err, PHP83)
		want := "Parse error: syntax error, unexpected token \"->\", expecting " + tc.expected + " in context.php on line 3\n"
		if code != 255 || out.Len() != 0 || err.String() != want {
			t.Fatalf("%s: exit=%d out=%q err=%q", tc.source, code, &out, &err)
		}
	}
}

func TestRecursiveObjectComparisonTerminates(t *testing.T) {
	for _, op := range []string{"==", "!=", "<>", "<=>", "<", "<=", ">", ">="} {
		var out, err bytes.Buffer
		source := fmt.Sprintf("<?php $a = new stdClass; $a->x = [$a]; $b = new stdClass; $b->x = [$b]; echo $a %s $b;", op)
		if exit := Run(source, "cycle.php", &out, &err); exit != 255 || !strings.Contains(err.String(), "recursive dependency") {
			t.Fatalf("%s: exit=%d out=%q err=%q", op, exit, &out, &err)
		}
	}
}
