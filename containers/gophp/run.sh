#!/bin/sh
set -eu
project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
image=${GOPHP_IMAGE:-localhost/gophp:v1}
# Relative file arguments always refer to the project root. Mount at the same
# absolute path so diagnostics can also be compared with the PHP references.
exec podman run --rm --pull=never --network=none --read-only \
  --cap-drop=all --security-opt=no-new-privileges \
  --volume "$project_root:$project_root:ro" --workdir "$project_root" \
  "$image" "$@"
