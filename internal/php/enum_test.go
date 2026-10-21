package php

import (
	"bytes"
	"testing"
)

func TestEnumCasesHaveStableIdentityAndFreshRuns(t *testing.T) {
	const source = `<?php
enum E { case A; case B; }
$a = E::A;
$array = [E::A];
echo $a === $array[0] ? 'same|' : 'bad|';
echo E::A === E::B ? 'bad' : 'different';
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for i := 0; i < 2; i++ {
			var out, diagnostic bytes.Buffer
			code := RunWithVersion(source, "enum-isolation.php", &out, &diagnostic, version)
			if code != 0 || out.String() != "same|different" || diagnostic.Len() != 0 {
				t.Fatalf("%s run %d: exit=%d stdout=%q stderr=%q", version, i, code, out.String(), diagnostic.String())
			}
		}
	}
}
