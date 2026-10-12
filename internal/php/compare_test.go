package php

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
)

func TestComparisonTypesAndNaN(t *testing.T) {
	// The columns are == != < <= > >=. These expectations include PHP's
	// intentionally non-total order for NaN, and null's string exception.
	s := func(text string) value { return value{kind: stringKind, s: text} }
	cases := []struct {
		name string
		a, b value
		bits string
	}{
		{"large integers stay distinct", integer(9007199254740993), integer(9007199254740992), "010011"},
		{"int float rounding", integer(9007199254740993), floating(9007199254740992), "100101"},
		{"numeric string overflow", s("9223372036854775807"), s("9223372036854775808"), "011100"},
		{"overflowing strings stay distinct", s("9223372036854775808"), s("9223372036854775809"), "011100"},
		{"overflow strings use original bytes", s("-9223372036854775809"), s("-9223372036854775810"), "011100"},
		{"number string rounding", integer(math.MaxInt64), s("9223372036854775808"), "100101"},
		{"non-numeric string", integer(0), s("word"), "011100"},
		{"prefix is not numeric", integer(100), s("2x"), "011100"},
		{"null empty", value{}, s(""), "100101"},
		{"null zero string", value{}, s("0"), "011100"},
		{"null negative", value{}, integer(-1), "011100"},
		{"NaN number", floating(math.NaN()), integer(0), "010000"},
		{"number NaN", integer(0), floating(math.NaN()), "010000"},
		{"NaN string", floating(math.NaN()), s("NAN"), "010000"},
		{"NaN true", floating(math.NaN()), boolean(true), "100101"},
		{"NaN false", floating(math.NaN()), boolean(false), "010011"},
		{"NaN null", floating(math.NaN()), value{}, "010011"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for i, op := range []string{"==", "!=", "<", "<=", ">", ">="} {
				for _, version := range []Version{PHP83, PHP84, PHP85} {
					var stderr bytes.Buffer
					e := evaluator{version: version, stderr: &stderr, stdout: io.Discard}
					node := binary{token{kind: op, text: op, line: 1}, literal{tc.a, 1}, literal{tc.b, 1}}
					got, err := e.eval(node)
					want := boolean(tc.bits[i] == '1')
					if err != nil || got != want || stderr.Len() != 0 {
						t.Fatalf("PHP %s %s: got %#v, err=%v stderr=%q; want %#v", version, op, got, err, &stderr, want)
					}
				}
			}
		})
	}
}

func TestComparisonAssociativity(t *testing.T) {
	groups := [][]string{{"<", "<=", ">", ">="}, {"==", "!=", "===", "!==", "<>", "<=>"}}
	for _, group := range groups {
		for _, first := range group {
			for _, second := range group {
				source := fmt.Sprintf("<?php echo 'must not execute'; echo 1 %s 2 %s 3;", first, second)
				var stdout, stderr bytes.Buffer
				if code := Run(source, "input.php", &stdout, &stderr); code != 255 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "unexpected token") {
					t.Fatalf("unparenthesized chain accepted: %s; exit=%d stdout=%q stderr=%q", source, code, &stdout, &stderr)
				}
			}
		}
	}
	for _, ordering := range groups[0] {
		for _, equality := range groups[1] {
			// Mixing the two precedence groups is legal. Ordering binds first
			// on either side of equality; explicit grouping remains legal.
			source := fmt.Sprintf(`<?php
echo (1 %s 2 %s true) === ((1 %s 2) %s true);
echo (true %s 1 %s 2) === (true %s (1 %s 2));
echo ((1 %s 2) %s 3) === ((1 %s 2) %s true);
`, ordering, equality, ordering, equality, equality, ordering, equality, ordering, ordering, ordering, ordering, ordering)
			var stdout, stderr bytes.Buffer
			if code := Run(source, "input.php", &stdout, &stderr); code != 0 || stdout.String() != "111" || stderr.Len() != 0 {
				t.Fatalf("grouping/precedence: %s; exit=%d stdout=%q stderr=%q", source, code, &stdout, &stderr)
			}
		}
	}
}
