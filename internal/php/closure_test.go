package php

import (
	"bytes"
	"testing"
)

func TestClosureAndArrowCaptureByValue(t *testing.T) {
	source := `<?php
$x = 2;
$f = function ($n) use ($x) { return $n + $x; };
$x = 10;
$arrow = fn ($n) => $n + $x;
$x = 100;
echo $f(3), '|', $arrow(1), '|', $x;
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "closure.php", &out, &diagnostics)
	if code != 0 || out.String() != "5|11|100" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestClosureKeepsCapturedStateAndMethodScope(t *testing.T) {
	source := `<?php
class Vault {
    private $secret = 4;
    public function reader() { return function () { return $this->secret; }; }
}
$f = (new Vault)->reader();
$x = 1;
$counter = function () use ($x) { return ++$x; };
echo $f(), '|', $counter(), $counter(), '|', $x;
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "closure.php", &out, &diagnostics)
	if code != 0 || out.String() != "4|22|1" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}

func TestClosureStaticLocalBelongsToInstance(t *testing.T) {
	source := `<?php
function makeCounter() { return function () { static $n = 0; return ++$n; }; }
$a = makeCounter(); $b = makeCounter();
echo $a(), $a(), $b();
`
	var out, diagnostics bytes.Buffer
	code := Run(source, "closure.php", &out, &diagnostics)
	if code != 0 || out.String() != "121" || diagnostics.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostics.String())
	}
}
