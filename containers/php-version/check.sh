#!/bin/sh
set -eu
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$script_dir/../.."

for spec in '83:8.3.35' '84:8.4.26' '85:8.5.11'; do
  short=${spec%%:*}
  expected=${spec#*:}
  actual=$("$script_dir/bin/php$short" -n -r 'echo PHP_VERSION;')
  if [ "$actual" != "$expected" ]; then
    echo "ERROR: php$short expected $expected, got $actual" >&2
    exit 1
  fi
  result=$("$script_dir/bin/php$short" -n "$script_dir/fixtures/smoke.php")
  if [ "$result" != 'php-reference-ok:14' ]; then
    echo "ERROR: php$short smoke result: $result" >&2
    exit 1
  fi
  printf 'PASS: PHP %s, shared fixture executed\n' "$actual"
done
