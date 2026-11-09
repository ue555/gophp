package php

import (
	"bytes"
	"testing"
)

func TestFilterConstantsAndSerializable(t *testing.T) {
	source := `<?php echo FILTER_FORCE_ARRAY, " ", FILTER_NULL_ON_FAILURE, " ", FILTER_FLAG_GLOBAL_RANGE, " ", interface_exists('Serializable') ? 'yes' : 'no';`
	want := "67108864 134217728 536870912 yes"
	var out, err bytes.Buffer
	if code := Run(source, "input.php", &out, &err); code != 0 || out.String() != want || err.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &err)
	}
}
