#!/usr/bin/env python3
"""Exercise the fixed Laravel/API Platform contract through gophp --serve."""

import argparse
from concurrent.futures import ThreadPoolExecutor
import json
import os
from pathlib import Path
import shutil
import socket
import sqlite3
import subprocess
import sys
import tempfile
import time
from urllib.error import URLError

from verify import APP, ROOT, request


def check(condition, description):
    if not condition:
        raise RuntimeError(description)


def free_port():
    with socket.socket() as listener:
        listener.bind(("127.0.0.1", 0))
        return listener.getsockname()[1]


def run_command(args, env):
    print("command:", " ".join(map(str, args)), flush=True)
    result = subprocess.run(args, cwd=APP, env=env, text=True, capture_output=True)
    print("exit_code:", result.returncode, flush=True)
    if result.stdout:
        print(result.stdout, end="", flush=True)
    if result.stderr:
        print(result.stderr, end="", file=sys.stderr, flush=True)
    check(result.returncode == 0, "CLI command failed")
    return result.stdout


def verify_http(port, database, stress_requests):
    status, payload, content_type = request(port, "GET", "/api/probe", include_content_type=True)
    check((status, payload) == (200, {"ok": True}) and content_type.split(";", 1)[0].strip() == "application/json", "probe response differs")
    status, payload = request(port, "GET", "/api/probe?round=0")
    check((status, payload) == (200, {"ok": True}), "probe query response differs")

    status, payload = request(port, "POST", "/api/books-manual", {"title": "Alpha"})
    check(status == 201 and payload.get("title") == "Alpha" and isinstance(payload.get("id"), int) and payload["id"] > 0, "manual POST differs")
    manual_id = payload["id"]
    status, payload = request(port, "GET", f"/api/books-manual/{manual_id}")
    check((status, payload) == (200, {"id": manual_id, "title": "Alpha"}), "manual GET differs")
    status, payload = request(port, "POST", "/api/books-manual", {})
    check(status == 422 and "title" in payload.get("errors", {}), "validation response differs")
    with sqlite3.connect(database) as connection:
        rows = connection.execute("SELECT id, title FROM books ORDER BY id").fetchall()
    check(rows == [(manual_id, "Alpha")], f"rows after manual validation: {rows!r}")

    status, payload, content_type = request(port, "POST", "/api/books", {"title": "Beta"}, "application/ld+json", include_content_type=True)
    check(status == 201 and payload.get("title") == "Beta" and payload.get("@id"), "API Platform POST differs")
    check(content_type.split(";", 1)[0].strip() == "application/ld+json", "API Platform POST content type differs")
    iri = payload["@id"]
    status, payload, content_type = request(port, "GET", iri, accept="application/ld+json", include_content_type=True)
    check(status == 200 and payload.get("title") == "Beta", "API Platform GET differs")
    check(content_type.split(";", 1)[0].strip() == "application/ld+json", "API Platform GET content type differs")
    with sqlite3.connect(database) as connection:
        rows = connection.execute("SELECT title FROM books ORDER BY id").fetchall()
    check(rows == [("Alpha",), ("Beta",)], f"rows after API Platform POST: {rows!r}")

    status, payload = request(port, "PATCH", iri, {"title": "Gamma"}, "application/ld+json", "application/merge-patch+json")
    check(status == 200 and payload.get("title") == "Gamma", "API Platform PATCH differs")
    status, payload = request(port, "GET", "/api/docs.jsonopenapi", accept="application/vnd.openapi+json")
    check(status == 200 and isinstance(payload, dict) and "openapi" in payload, "OpenAPI response differs")
    book_path = payload.get("paths", {}).get("/api/books", {})
    check("get" in book_path and "post" in book_path, "OpenAPI lacks Book collection GET/POST")
    status, _ = request(port, "DELETE", iri, accept="application/ld+json")
    check(status == 204, "API Platform DELETE differs")
    status, _ = request(port, "GET", iri, accept="application/ld+json")
    check(status == 404, "deleted resource remains visible")
    with sqlite3.connect(database) as connection:
        rows = connection.execute("SELECT title FROM books ORDER BY id").fetchall()
    check(rows == [("Alpha",)], f"rows after API Platform DELETE: {rows!r}")

    for index in range(stress_requests):
        path = f"/api/books-manual/{manual_id}" if index % 8 == 0 else f"/api/probe?round={index}"
        status, payload = request(port, "GET", path)
        expected = {"id": manual_id, "title": "Alpha"} if index % 8 == 0 else {"ok": True}
        check((status, payload) == (200, expected), f"sustained request {index} differs")
    with ThreadPoolExecutor(max_workers=8) as pool:
        results = list(pool.map(lambda index: request(port, "GET", f"/api/probe?parallel={index}"), range(8)))
    check(all(result == (200, {"ok": True}) for result in results), "parallel probe responses differ")
    with sqlite3.connect(database) as connection:
        rows = connection.execute("SELECT title FROM books ORDER BY id").fetchall()
    check(rows == [("Alpha",)], f"rows after sustained requests: {rows!r}")
    print(f"PASS sustained_requests={stress_requests} parallel_requests=8", flush=True)


def verify_version(version, executable, cli_only=False, stress_requests=96):
    cache = ROOT / ".cache"
    cache.mkdir(exist_ok=True)
    with tempfile.TemporaryDirectory(prefix=f"gophp-app-{version}-", dir=cache) as temporary:
        temporary = Path(temporary)
        database = temporary / "database.sqlite"
        database.touch()
        storage = temporary / "storage"
        for directory in ("framework/cache/data", "framework/views", "framework/sessions", "logs"):
            (storage / directory).mkdir(parents=True)
        env = os.environ.copy()
        env.update({
            "APP_ENV": "testing",
            "APP_KEY": "base64:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
            "DB_CONNECTION": "sqlite",
            "DB_DATABASE": str(database),
            "LARAVEL_STORAGE_PATH": str(storage),
            "LOG_CHANNEL": "stderr",
        })
        run_command([str(executable), f"--php-version={version}", "artisan", "migrate", "--force"], env)
        # migrate boots the app before the books table exists, so API Platform
        # caches empty property metadata and generates no item routes.
        shutil.rmtree(storage / "framework/cache/data")
        (storage / "framework/cache/data").mkdir(parents=True)
        if cli_only:
            output = run_command([str(executable), f"--php-version={version}", "artisan", "gophp:probe"], env)
            check("probe-ok" in output, "custom Artisan command did not produce probe-ok")
            output = run_command([str(executable), f"--php-version={version}", "artisan", "api:openapi:export"], env)
            specification = json.loads(output)
            book_path = specification.get("paths", {}).get("/api/books", {})
            check(specification.get("openapi") and "get" in book_path and "post" in book_path,
                  "CLI OpenAPI lacks version or Book collection GET/POST")
            with sqlite3.connect(database) as connection:
                rows = connection.execute("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'books'").fetchall()
            check(rows == [("books",)], f"migration did not create books: {rows!r}")
            print(f"PASS Go PHP {version}: CLI migration, custom command, OpenAPI, SQLite table", flush=True)
            return
        port = free_port()
        command = [str(executable), f"--php-version={version}", f"--serve=127.0.0.1:{port}", "public/index.php"]
        print("command:", " ".join(command), flush=True)
        log_path = temporary / "server.log"
        with log_path.open("w+") as server_log:
            process = subprocess.Popen(command, cwd=APP, env=env, stdout=server_log, stderr=subprocess.STDOUT)
            try:
                for _ in range(150):
                    if process.poll() is not None:
                        raise RuntimeError(f"server exited before ready: {process.returncode}")
                    try:
                        request(port, "GET", "/api/probe")
                        break
                    except (URLError, TimeoutError):
                        time.sleep(0.2)
                else:
                    raise RuntimeError("server did not become ready")
                verify_http(port, database, stress_requests)
                print(f"PASS Go PHP {version}: CLI, HTTP, SQLite, validation, API Platform CRUD, OpenAPI", flush=True)
            finally:
                process.terminate()
                try:
                    process.wait(timeout=5)
                except subprocess.TimeoutExpired:
                    process.kill()
                    process.wait()
                print("server_exit_code:", process.returncode, flush=True)
                server_log.flush()
                server_log.seek(0)
                print(f"server_log_begin Go PHP {version}", flush=True)
                print(server_log.read(), end="", flush=True)
                print(f"server_log_end Go PHP {version}", flush=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--version", choices=("8.3", "8.4", "8.5"), action="append", dest="versions")
    parser.add_argument("--interpreter", type=Path, default=ROOT / "bin/gophp")
    parser.add_argument("--cli-only", action="store_true", help="check migration, custom Artisan command and OpenAPI without binding an HTTP port")
    parser.add_argument("--stress-requests", type=int, default=96, help="number of sustained TCP GET requests (96-2000)")
    args = parser.parse_args()
    if not 96 <= args.stress_requests <= 2000:
        parser.error("--stress-requests must be between 96 and 2000")
    executable = args.interpreter.resolve()
    if not executable.is_file() or not (APP / "vendor/autoload.php").is_file():
        print("gophp binary and fixed vendor/autoload.php are required", file=sys.stderr)
        return 2
    for version in args.versions or ("8.3", "8.4", "8.5"):
        try:
            verify_version(version, executable, args.cli_only, args.stress_requests)
        except PermissionError as error:
            print(f"ENVIRONMENT ERROR Go PHP {version}: {error}", file=sys.stderr)
            return 2
        except (OSError, RuntimeError, ValueError, KeyError, TypeError, json.JSONDecodeError) as error:
            print(f"FAIL Go PHP {version}: {error}", file=sys.stderr)
            return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
