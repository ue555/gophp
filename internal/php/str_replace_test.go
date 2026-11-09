package php

import (
	"bytes"
	"strings"
	"testing"
)

func TestStrReplaceArraySemantics(t *testing.T) {
	source := `<?php
echo str_replace(['-', '_'], ' ', 'todo-item_x'), "\n";
echo str_replace(['a', 'b'], ['x'], 'ab'), "\n";
echo str_replace(['a', 'b'], ['x', 'y', 'z'], 'ab'), "\n";
echo str_replace(['a', 'b'], ' ', 'ab'), "\n";
echo implode('|', str_replace(['a', 'b'], ' ', ['ab', 'ba'])), "\n";
`
	want := "todo item x\nx\nxy\n  \n  |  \n"
	var out, err bytes.Buffer
	if code := Run(source, "input.php", &out, &err); code != 0 || out.String() != want || err.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &err)
	}
}

func TestStrReplaceRejectsArrayReplaceWithStringSearch(t *testing.T) {
	var out, err bytes.Buffer
	code := Run(`<?php str_replace('a', ['x'], 'a');`, "input.php", &out, &err)
	if code != 255 || !strings.Contains(err.String(), "Argument #2 ($replace) must be of type string when argument #1 ($search) is a string") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &out, &err)
	}
}
