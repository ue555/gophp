package php

import (
	"bytes"
	"strings"
	"testing"
)

func TestRelativeTypesAndDeferredCompatibility(t *testing.T) {
	tests := []struct {
		name, source, output, diagnostic string
		exit                             int
	}{
		{"inherited self uses declaring class", `<?php class A { public function same(): self { return $this; } } class B extends A {} echo (new B)->same() instanceof B ? 'yes' : 'no';`, "yes", "", 0},
		{"static return uses called class", `<?php class A { public function make(): static { return new A; } } class B extends A {} try { (new B)->make(); } catch (TypeError $e) { echo 'caught'; }`, "caught", "", 0},
		{"parent property uses declaring class", `<?php class A {} class B extends A { public parent $value; } $b = new B; $b->value = new A; echo $b->value instanceof A ? 'yes' : 'no';`, "yes", "", 0},
		{"later type in same file cannot prove variance", `<?php class Base { public function make(): LaterBase { return new LaterBase; } } class Child extends Base { public function make(): LaterChild { return new LaterChild; } } class LaterBase {} class LaterChild extends LaterBase {} echo 'unreachable';`, "", "because class LaterChild is not available", 255},
		{"type known when child is declared", `<?php class Base { public function make(): LaterBase { return new LaterBase; } } class LaterBase {} class LaterChild extends LaterBase {} class Child extends Base { public function make(): LaterChild { return new LaterChild; } } echo (new Child)->make() instanceof LaterBase ? 'yes' : 'no';`, "yes", "", 0},
		{"same unknown name", `<?php interface Contract { public function make(): UnknownType; } class Implementation implements Contract { public function make(): UnknownType { return new UnknownType; } } echo 'declared';`, "declared", "", 0},
		{"unknown hierarchy cannot be assumed", `<?php echo 'before'; interface Contract { public function make(): MissingBase; } class Implementation implements Contract { public function make(): MissingChild { return new MissingChild; } }`, "before", "Could not check compatibility", 255},
		{"unknown return type names offending child", `<?php class Base { public function make(): MissingBase {} } class Child extends Base { public function make(): MissingChild {} }`, "", "because class MissingChild is not available", 255},
		{"self in compatibility diagnostic resolves declaration", `<?php class Base { public function make(): static {} } class Child extends Base { public function make(): self {} }`, "", "Child::make(): Child must be compatible with Base::make(): static", 255},
		{"parent without parent", `<?php class Lone { public function make(): parent {} }`, "", "Cannot use \"parent\"", 255},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, version := range []Version{PHP83, PHP84, PHP85} {
				var out, diagnostics bytes.Buffer
				exit := RunWithVersion(tc.source, "special.php", &out, &diagnostics, version)
				if exit != tc.exit || out.String() != tc.output || !strings.Contains(diagnostics.String(), tc.diagnostic) || tc.diagnostic == "" && diagnostics.Len() != 0 {
					t.Errorf("%s: exit=%d stdout=%q stderr=%q", version, exit, out.String(), diagnostics.String())
				}
			}
		})
	}
}
