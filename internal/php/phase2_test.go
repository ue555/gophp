package php

import (
	"bytes"
	"strings"
	"testing"
)

func TestPhase2DeadBranchValidation(t *testing.T) {
	cases := []struct{ source, diagnostic string }{
		{`<?php echo "before"; if (false) { $a?->x = 1; }`, "Can't use nullsafe operator in write context"},
		{`<?php echo "before"; if (false) { echo match (1) { default => 1, default => 2 }; }`, "Match expressions may only contain one default arm"},
		{`<?php echo "before"; if (false) { switch (1) { default: break; default: break; } }`, "Switch statements may only contain one default clause"},
		{`<?php echo "before"; if (false) { echo "${old}"; }`, "unsupported syntax: deprecated ${} interpolation"},
		{`<?php echo "before"; if (false) { const X = 1; }`, "unexpected token"},
		{`<?php echo "before"; namespace N;`, "Namespace declaration statement has to be the very first statement"},
		{`<?php namespace A; namespace B {} `, "Cannot mix bracketed namespace declarations"},
		{`<?php namespace A { const X = strlen("x"); }`, "Constant expression contains invalid operations"},
		{`<?php namespace A { const X = 1; } namespace A { const X = 2; }`, "Cannot redeclare constant A\\X"},
		{`<?php namespace N; use Other\Thing as Alias; use Another\Thing as Alias;`, "because the name is already in use"},
	}
	for _, tc := range cases {
		t.Run(tc.diagnostic+tc.source, func(t *testing.T) {
			for _, version := range []Version{PHP83, PHP84, PHP85} {
				var out, diagnostics bytes.Buffer
				code := RunWithVersion(tc.source, "input.php", &out, &diagnostics, version)
				if code == 0 || out.Len() != 0 || !strings.Contains(diagnostics.String(), tc.diagnostic) {
					t.Errorf("PHP %s: exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
				}
			}
		})
	}
}

func TestPhase2RuntimeErrorType(t *testing.T) {
	for _, tc := range []struct{ source, diagnostic string }{
		{`<?php echo 1 << -1;`, "Uncaught ArithmeticError: Bit shift by negative number"},
		{`<?php echo match (1) { 2 => 3 };`, "Uncaught UnhandledMatchError: Unhandled match case 1"},
		{`<?php echo match ('abc') { 'other' => 3 };`, "Uncaught UnhandledMatchError: Unhandled match case 'abc'"},
		{`<?php namespace A; echo MISSING;`, `Uncaught Error: Undefined constant "A\\MISSING"`},
	} {
		var out, diagnostics bytes.Buffer
		code := Run(tc.source, "input.php", &out, &diagnostics)
		if code != 255 || out.Len() != 0 || !strings.Contains(diagnostics.String(), tc.diagnostic) {
			t.Errorf("%s: exit=%d stdout=%q stderr=%q", tc.source, code, out.String(), diagnostics.String())
		}
	}
}
