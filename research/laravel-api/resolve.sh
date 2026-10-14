#!/bin/sh
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
app="$root/research/laravel-api"
image=localhost/gophp-dependency-research:8.3
podman build --build-arg PHP_VERSION=8.3.35 -f "$root/containers/dependency-research/Containerfile" -t "$image" "$root"
podman run --rm "$image" php -v
podman run --rm "$image" php -m
podman run --rm "$image" composer --version
podman run --rm -v "$app:/app" "$image" composer validate --no-check-publish
podman run --rm -v "$app:/app" "$image" composer update --no-dev --no-interaction --prefer-dist --no-scripts --no-plugins
for version in 8.3.35 8.4.26 8.5.11; do
    tag="localhost/gophp-dependency-research:$version"
    podman build --build-arg "PHP_VERSION=$version" -f "$root/containers/dependency-research/Containerfile" -t "$tag" "$root"
    podman run --rm "$tag" php -v
    podman run --rm "$tag" php -m
    podman run --rm "$tag" composer --version
    podman run --rm -v "$app:/app" "$tag" composer install --no-dev --no-interaction --prefer-dist --no-scripts --no-plugins
    podman run --rm -v "$app:/app" "$tag" composer check-platform-reqs --no-dev
done
