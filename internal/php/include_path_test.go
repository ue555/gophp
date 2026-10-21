package php

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIncludePathSearchAndRunIsolation(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")
	for _, dir := range []string{first, second} {
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	for path, source := range map[string]string{
		filepath.Join(root, "choice.inc"):   "<?php return 'local';",
		filepath.Join(first, "choice.inc"):  "<?php return 'first';",
		filepath.Join(second, "choice.inc"): "<?php return 'second';",
		filepath.Join(second, "outer.inc"):  "<?php return include 'inner.inc';",
		filepath.Join(second, "inner.inc"):  "<?php return 'inner';",
		filepath.Join(root, "fallback.inc"): "<?php return 'fallback';",
	} {
		if err := os.WriteFile(path, []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	main := filepath.Join(root, "main.php")
	source := "<?php $original = get_include_path(); " +
		"echo set_include_path('" + second + ":" + first + "') === $original ? 'old|' : 'bad|'; " +
		"echo include 'choice.inc', '|'; " +
		"echo include 'fallback.inc', '|'; " +
		"set_include_path('" + first + ":" + second + "'); " +
		"echo include 'choice.inc', '|'; " +
		"echo include '" + filepath.Join(second, "outer.inc") + "', '|'; " +
		"ini_restore('include_path'); echo get_include_path() === $original ? 'restored' : 'bad';"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for i := 0; i < 2; i++ {
			var stdout, stderr bytes.Buffer
			code := RunWithVersion(source, main, &stdout, &stderr, version)
			if code != 0 || stdout.String() != "old|second|fallback|first|inner|restored" || stderr.Len() != 0 {
				t.Fatalf("version %s run %d: exit=%d stdout=%q stderr=%q", version, i, code, stdout.String(), stderr.String())
			}
		}
	}
}

func TestIncludePathCanonicalOnceAndMissingRetry(t *testing.T) {
	root := t.TempDir()
	first, second := filepath.Join(root, "first"), filepath.Join(root, "second")
	for _, dir := range []string{first, second} {
		if err := os.Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(second, "one.inc"), []byte("<?php $visits++; return $visits;"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("one.inc", filepath.Join(second, "alias.inc")); err != nil {
		t.Fatal(err)
	}
	source := "<?php $visits=0; set_error_handler(function($n,$m){return true;}); " +
		"set_include_path('" + first + "'); $missing=include_once 'one.inc'; " +
		"echo ($missing ? 'found' : 'missing'), '|'; " +
		"set_include_path('" + second + "'); echo include_once 'one.inc', '|', " +
		"include_once 'alias.inc', '|', include_once '" + filepath.Join(second, "./one.inc") + "', '|', $visits;"
	var stdout, stderr bytes.Buffer
	if code := Run(source, filepath.Join(root, "main.php"), &stdout, &stderr); code != 0 || stdout.String() != "missing|1|1|1|1" || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestIncludePathHandlerThrowStopsEvaluation(t *testing.T) {
	root := t.TempDir()
	source := "<?php $calls=0; set_include_path('" + root + "'); " +
		"set_error_handler(function($n,$m){ throw new Exception('stop'); }); " +
		"try { include 'missing.inc'; $calls++; } catch(Exception $e) { echo $e->getMessage(),'|',$calls; }"
	var stdout, stderr bytes.Buffer
	if code := Run(source, filepath.Join(root, "main.php"), &stdout, &stderr); code != 0 || stdout.String() != "stop|0" || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := Run("<?php echo get_include_path();", filepath.Join(root, "again.php"), &stdout, &stderr); code != 0 || !strings.HasPrefix(stdout.String(), ".:") || stderr.Len() != 0 {
		t.Fatalf("second run exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestIncludePathStrictNullDoesNotChangeState(t *testing.T) {
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		for _, call := range []string{"set_include_path(null)", "ini_restore(null)"} {
			source := "<?php declare(strict_types=1); $before=get_include_path(); " +
				"try { " + call + "; echo 'accepted'; } catch(TypeError $e) { echo 'TypeError'; } " +
				"echo '|', get_include_path()===$before ? 'unchanged' : 'changed';"
			var stdout, stderr bytes.Buffer
			code := RunWithVersion(source, "strict-null.php", &stdout, &stderr, version)
			if code != 0 || stdout.String() != "TypeError|unchanged" || stderr.Len() != 0 {
				t.Fatalf("version %s %s: exit=%d stdout=%q stderr=%q", version, call, code, stdout.String(), stderr.String())
			}
		}
	}
}

func TestSetIncludePathEmptyReturnsFalseAndKeepsPath(t *testing.T) {
	const source = "<?php $before=get_include_path(); " +
		"echo set_include_path('') === false ? 'false|' : 'bad|'; " +
		"echo get_include_path() === $before ? 'unchanged' : 'changed';"
	for _, version := range []Version{PHP83, PHP84, PHP85} {
		var stdout, stderr bytes.Buffer
		code := RunWithVersion(source, "empty-include-path.php", &stdout, &stderr, version)
		if code != 0 || stdout.String() != "false|unchanged" || stderr.Len() != 0 {
			t.Fatalf("version %s: exit=%d stdout=%q stderr=%q", version, code, stdout.String(), stderr.String())
		}
	}
}
