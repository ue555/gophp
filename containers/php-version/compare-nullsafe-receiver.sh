#!/bin/sh
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
mkdir -p "$root/.cache"
fixture_dir=$(mktemp -d "$root/.cache/nullsafe-receiver.XXXXXX")
trap 'rm -rf "$fixture_dir"' EXIT HUP INT TERM

for relative in \
    runtime/nullsafe_receiver_error_order.php \
    runtime/nullsafe_receiver_missing_warning.php \
    runtime/nullsafe_receiver_private.php \
    runtime/nullsafe_receiver_protected.php \
    success/nullsafe_receiver_inside.php \
    success/nullsafe_receiver_offset_once.php \
    success/nullsafe_receiver_prior_stop.php \
    success/nullsafe_receiver_public.php \
    success/nullsafe_receiver_undefined.php \
    success/nullsafe_quiet_missing.php \
    success/nullsafe_quiet_short_circuit.php; do
    destination="$fixture_dir/${relative%/*}"
    mkdir -p "$destination"
    cp "$root/tests/fixtures/$relative" "$destination/"
done

"$root/containers/php-version/compare-all.sh" "$fixture_dir" "$root/bin/gophp"
