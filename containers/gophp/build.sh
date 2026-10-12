#!/bin/sh
set -eu
project_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
image=${GOPHP_IMAGE:-localhost/gophp:v1}
exec podman build --file "$project_root/containers/gophp/Containerfile" \
  --tag "$image" "$@" "$project_root"
