#!/usr/bin/env python3
"""Run the isolated Laravel/API Platform HTTP probe on three official PHP images."""

import json
import shutil
from pathlib import Path
import sqlite3
import subprocess
import sys
import tempfile
import time
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen


APP = Path(__file__).resolve().parent
ROOT = APP.parents[1]
VERSIONS = ("8.3.35", "8.4.26", "8.5.11")


def command(args):
    print("command:", " ".join(args), flush=True)
    result = subprocess.run(args, text=True, capture_output=True)
    print("exit:", result.returncode, flush=True)
    if result.stdout:
        print(result.stdout, end="", flush=True)
    if result.stderr:
        print(result.stderr, end="", file=sys.stderr, flush=True)
    if result.returncode:
        raise RuntimeError("command failed")
    return result.stdout.strip()


def request(port, method, path, body=None, accept="application/json", content_type=None, include_content_type=False):
    data = None if body is None else json.dumps(body).encode()
    headers = {"Accept": accept}
    if data is not None:
        headers["Content-Type"] = content_type or accept
    req = Request(f"http://127.0.0.1:{port}{path}", data=data,
                  headers=headers, method=method)
    def decoded(raw):
        text = raw.decode("utf-8", errors="replace")
        try:
            return json.loads(text)
        except json.JSONDecodeError:
            return text

    try:
        with urlopen(req, timeout=5) as response:
            result = response.status, decoded(response.read())
            response_content_type = response.headers.get("Content-Type", "")
    except HTTPError as error:
        result = error.code, decoded(error.read())
        response_content_type = error.headers.get("Content-Type", "")
    print(f"HTTP {method} {path}: {result[0]} {result[1]!r}", flush=True)
    if include_content_type:
        return result[0], result[1], response_content_type
    return result


def check(version):
    image = f"localhost/gophp-dependency-research:{version}"
    with tempfile.TemporaryDirectory(prefix="laravel-api-", dir=ROOT / ".cache") as temporary:
        db = Path(temporary) / "database.sqlite"
        db.touch()
        # Metadata is cached per DB schema, so each run needs its own cache.
        cache = Path(temporary) / "cache"
        cache.mkdir()
        mounts = ["-v", f"{APP}:/app", "-v", f"{temporary}:/probe",
                  "-v", f"{cache}:/app/storage/framework/cache/data"]
        env = ["-e", "DB_DATABASE=/probe/database.sqlite",
               "-e", "DB_CONNECTION=sqlite", "-e", "APP_ENV=testing",
               "-e", "LOG_CHANNEL=stderr",
               "-e", "APP_KEY=base64:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="]
        command(["podman", "run", "--rm", *mounts, *env, image,
                 "php", "artisan", "package:discover", "--no-interaction"])
        command(["podman", "run", "--rm", *mounts, *env, image,
                 "php", "artisan", "migrate", "--force"])
        # Booting before the migration cached metadata for a missing table.
        shutil.rmtree(cache)
        cache.mkdir()
        command(["podman", "run", "--rm", *mounts, *env, image,
                 "php", "artisan", "api-platform:install", "--no-interaction"])
        dump = Path(temporary) / "metadata.dump"
        command(["podman", "run", "--rm", *mounts, *env, image,
                 "php", "artisan", "api-platform:metadata:dump",
                 "--path=/probe/metadata.dump", "--no-interaction"])
        if not dump.is_file() or dump.stat().st_size == 0:
            raise RuntimeError("metadata:dump returned success without writing /probe/metadata.dump")
        print(f"metadata dump: {dump.stat().st_size} bytes", flush=True)
        container = command(["podman", "run", "-d", "--rm",
                             "-p", "127.0.0.1::8080", *mounts, *env, image,
                             "php", "-S", "0.0.0.0:8080", "-t", "public", "public/index.php"])
        try:
            mapping = command(["podman", "port", container, "8080/tcp"])
            port = int(mapping.rsplit(":", 1)[1])
            for _ in range(50):
                try:
                    status, payload = request(port, "GET", "/api/probe")
                    break
                except (URLError, TimeoutError):
                    time.sleep(0.2)
            else:
                raise RuntimeError("HTTP server did not become ready")
            assert (status, payload) == (200, {"ok": True}), (status, payload)
            status, payload = request(port, "POST", "/api/books-manual", {"title": "Alpha"})
            assert status == 201 and payload["title"] == "Alpha", (status, payload)
            manual_id = payload["id"]
            status, payload = request(port, "GET", f"/api/books-manual/{manual_id}")
            assert (status, payload) == (200, {"id": manual_id, "title": "Alpha"}), (status, payload)
            status, payload = request(port, "POST", "/api/books-manual", {})
            assert status == 422 and "title" in payload.get("errors", {}), (status, payload)
            with sqlite3.connect(db) as connection:
                rows = connection.execute("select id, title from books order by id").fetchall()
            assert rows == [(manual_id, "Alpha")], rows
            status, payload = request(port, "POST", "/api/books", {"title": "Beta"},
                                      "application/ld+json")
            assert status == 201 and payload["title"] == "Beta", (status, payload)
            iri = payload["@id"]
            status, payload = request(port, "GET", iri, accept="application/ld+json")
            assert status == 200 and payload["title"] == "Beta", (status, payload)
            with sqlite3.connect(db) as connection:
                rows = connection.execute("select title from books order by id").fetchall()
            assert rows == [("Alpha",), ("Beta",)], rows
            print(f"PASS PHP {version}: autoload, JSON, SQLite, validation, API Platform", flush=True)
        finally:
            try:
                command(["podman", "logs", container])
            finally:
                command(["podman", "stop", container])


def main():
    if not (APP / "composer.lock").is_file() or not (APP / "vendor/autoload.php").is_file():
        print("ERROR: composer.lock and vendor/autoload.php are required; run resolve.sh first", file=sys.stderr)
        return 2
    (ROOT / ".cache").mkdir(exist_ok=True)
    (APP / "bootstrap/cache").mkdir(parents=True, exist_ok=True)
    (APP / "storage/framework/sessions").mkdir(parents=True, exist_ok=True)
    for version in VERSIONS:
        try:
            check(version)
        except (AssertionError, RuntimeError, OSError, ValueError, KeyError) as error:
            print(f"FAIL PHP {version}: {error}", file=sys.stderr)
            return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
