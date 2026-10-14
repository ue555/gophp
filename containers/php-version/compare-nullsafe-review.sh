#!/bin/sh
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
mkdir -p "$root/.cache"
fixture_dir=$(mktemp -d "$root/.cache/nullsafe-review.XXXXXX")
trap 'rm -rf "$fixture_dir"' EXIT HUP INT TERM

for relative in \
    success/nullsafe_quiet_short_circuit.php \
    success/nullsafe_quiet_missing.php \
    success/nullsafe_quiet_nested.php \
    success/nullsafe_quiet_side_effects.php \
    success/nullsafe_legal_unset.php \
    runtime/nullsafe_normal_read_warning.php \
    runtime/nullsafe_key_warning.php \
    runtime/nullsafe_argument_warning.php \
    parser/unset_nullsafe_direct.php \
    parser/unset_nullsafe_property_chain.php \
    parser/unset_nullsafe_offset_chain.php; do
    destination="$fixture_dir/${relative%/*}"
    mkdir -p "$destination"
    cp "$root/tests/fixtures/$relative" "$destination/"
done

"$root/containers/php-version/compare-all.sh" "$fixture_dir" "$root/bin/gophp"
