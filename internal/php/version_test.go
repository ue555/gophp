package php

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestNaNConditionVersions(t *testing.T) {
	const init = "<?php\n$n = (1e308 * 10) - (1e308 * 10);\n"
	for _, test := range []struct {
		name, source, out string
		lines             []int
	}{
		{"condition and unchanged value", init + "if ($n) { echo 'yes'; } else { echo 'no'; }\necho $n !== $n;", "yes1", []int{3}},
		{"each conversion warns", init + "if ($n) { echo 'a'; }\nif ($n) { echo 'b'; }", "ab", []int{3, 4}},
		{"skip then", init + "if (false) { if ($n) { echo 'bad'; } } else { echo 'ok'; }", "ok", nil},
		{"skip else", init + "if (true) { echo 'ok'; } else { if ($n) { echo 'bad'; } }", "ok", nil},
		{"selected nested condition", init + "if ($n) {\nif ($n) { echo 'ok'; }\n}", "ok", []int{3, 4}},
		{"strict comparison does not convert NaN", init + "if ($n !== $n) { echo 'ok'; }", "ok", nil},
		{"other truthy types", "<?php if (1e309) { echo 'i'; } if ('NAN') { echo 's'; } if (-0.0) { echo 'bad'; }", "is", nil},
		{"multiline CRLF condition", strings.ReplaceAll(init+"if (\n($n\n)\n) { echo 'ok'; }", "\n", "\r\n"), "ok", []int{4}},
		{"multiline arithmetic condition", "<?php\nif (\n(1e308 * 10) -\n(1e308 * 10)\n) { echo 'ok'; }", "ok", []int{4}},
	} {
		for _, version := range []Version{PHP83, PHP84, PHP85} {
			t.Run(test.name+"/"+string(version), func(t *testing.T) {
				var stdout, stderr bytes.Buffer
				code := RunWithVersion(test.source, "input.php", &stdout, &stderr, version)
				want := ""
				if version == PHP85 {
					for _, line := range test.lines {
						want += fmt.Sprintf("Warning: unexpected NAN value was coerced to bool in input.php on line %d\n", line)
					}
				}
				if code != 0 || stdout.String() != test.out || stderr.String() != want {
					t.Fatalf("exit=%d stdout=%q stderr=%q; want stdout=%q stderr=%q", code, &stdout, &stderr, test.out, want)
				}
			})
		}
	}
}

func TestRejectUnknownVersionBeforeExecution(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := RunWithVersion("<?php echo 'must not run';", "input.php", &stdout, &stderr, "8.6")
	if code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "unsupported PHP version") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
}
