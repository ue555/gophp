#!/bin/sh
set -eu
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
if [ "$#" -lt 2 ]; then
  echo "Usage: $0 FIXTURE_PATH INTERPRETER [INTERPRETER_ARGS...]" >&2
  echo 'The interpreter must accept --php-version=8.x; the PHP file path is appended last.' >&2
  exit 2
fi
fixtures=$1
shift
cd "$script_dir/../.."
status=0
for spec in '83:8.3.35' '84:8.4.26' '85:8.5.11'; do
  short=${spec%%:*}
  expected=${spec#*:}
  target=${expected%.*}
  printf '\nComparing with PHP %s\n' "$expected"
  printf 'Go target: PHP %s\n' "$target"
  if python3 skills/php-go-interpreter/scripts/compare_php.py \
    --php "$script_dir/bin/php$short" --php-version "$expected" \
    --fixtures "$fixtures" --timeout 15 --interpreter "$@" "--php-version=$target" '{file}'; then
    :
  else
    result=$?
    if [ "$result" -gt "$status" ]; then status=$result; fi
  fi
done
exit "$status"
