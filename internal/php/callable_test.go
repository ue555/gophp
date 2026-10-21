package php

import (
	"bytes"
	"testing"
)

func TestCallableTargetsAndClosureState(t *testing.T) {
	source := `<?php
function combine(int $a, int $b = 1): int { return $a + $b; }
class Worker {
    private function hidden(int $n): int { return $n + 4; }
    public function saved(): Closure { return $this->hidden(...); }
}
$worker = new Worker;
$callback = $worker->saved();
$sum = combine(...);
echo $callback(2), '|', $sum(b: 3, a: 5), '|';
echo (int)is_callable([$worker, 'hidden']), (int)is_callable($callback);
`
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for repeat := 0; repeat < 2; repeat++ {
			var stdout, stderr bytes.Buffer
			code := RunWithVersion(source, "callable.php", &stdout, &stderr, version)
			if code != 0 || stdout.String() != "6|8|01" || stderr.Len() != 0 {
				t.Fatalf("%s run %d: exit=%d stdout=%q stderr=%q", version, repeat, code, &stdout, &stderr)
			}
		}
	}
}

func TestCallablePropertyAndPromotionAreRejected(t *testing.T) {
	for _, source := range []string{
		"<?php class C { public callable $callback; }",
		"<?php class C { public function __construct(public callable $callback) {} }",
	} {
		var stdout, stderr bytes.Buffer
		code := RunWithVersion(source, "callable-property.php", &stdout, &stderr, PHP84)
		if code != 255 || stdout.Len() != 0 || !bytes.Contains(stderr.Bytes(), []byte("cannot have type callable")) {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
		}
	}
}

func TestCallableBuiltinsBindNamedArgumentsBeforeInvocation(t *testing.T) {
	source := `<?php
function add(int $a, int $b): int { return $a + $b; }
echo call_user_func_array(args: [2, 3], callback: 'add'), '|';
echo call_user_func(b: 3, callback: 'add', a: 2), '|';
try { call_user_func_array(unknown: 1, callback: 'add', args: [2, 3]); }
catch (Error $e) { echo $e->getMessage(), '|'; }
try { call_user_func_array(...['callback' => 'add'], callback: 'add', args: [2, 3]); }
catch (Error $e) { echo $e->getMessage(); }
`
	const want = "5|5|Unknown named parameter $unknown|Named parameter $callback overwrites previous argument"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var stdout, stderr bytes.Buffer
		code := RunWithVersion(source, "callable-named.php", &stdout, &stderr, version)
		if code != 0 || stdout.String() != want || stderr.Len() != 0 {
			t.Fatalf("%s: exit=%d stdout=%q stderr=%q", version, code, &stdout, &stderr)
		}
	}
}
