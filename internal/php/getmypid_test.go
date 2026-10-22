package php

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestGetmypidReportsHostProcess(t *testing.T) {
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var out, diagnostics bytes.Buffer
		code := RunWithArguments("<?php namespace Probe; echo getmypid();", "pid.php", []string{"pid.php"}, &out, &diagnostics, version)
		if code != 0 || out.String() != fmt.Sprint(os.Getpid()) || diagnostics.Len() != 0 {
			t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, code, out.String(), diagnostics.String())
		}
	}
}
