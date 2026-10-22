#!/usr/bin/env python3
"""Check the real CLI, database and TCP server. Never fall back from Go to PHP."""
import argparse
import concurrent.futures
import datetime
import hashlib
import json
import os
from pathlib import Path
import socket
import sqlite3
import subprocess
import sys
import tempfile
import time
import urllib.error
import urllib.request
import uuid

ROOT = Path(__file__).resolve().parent.parent


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--official', action='store_true', help='Explicit official-PHP reference run in Podman')
    args = parser.parse_args()
    stamp = datetime.datetime.now().strftime('%Y%m%d-%H%M%S') + '-' + uuid.uuid4().hex[:6]
    kind = 'official' if args.official else 'gophp'
    logs = ROOT / 'docs/validation' / f'{stamp}-{kind}'
    logs.mkdir(parents=True)
    cache = ROOT / '.cache'
    cache.mkdir(exist_ok=True)
    results = []
    opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
    server = None
    container_name = 'gophp-app-verify-' + uuid.uuid4().hex[:12]
    with tempfile.TemporaryDirectory(prefix='verify-', dir=cache) as tmp:
        tmp = Path(tmp)
        db = tmp / 'database.sqlite'
        db.touch()
        views = tmp / 'views'
        views.mkdir()
        env = os.environ.copy()
        overrides = {
            'GOPHP_VERSION': '8.5', 'APP_ENV': 'testing', 'APP_DEBUG': 'false',
            'DB_CONNECTION': 'sqlite', 'DB_DATABASE': str(db),
            'VIEW_COMPILED_PATH': str(views),
            'APP_CONFIG_CACHE': str(tmp / 'config.php'),
            'APP_SERVICES_CACHE': str(tmp / 'services.php'),
            'APP_PACKAGES_CACHE': str(tmp / 'packages.php'),
            'APP_ROUTES_CACHE': str(tmp / 'routes.php'),
            'APP_EVENTS_CACHE': str(tmp / 'events.php'),
            'LOG_CHANNEL': 'stderr',
        }
        env.update(overrides)
        image = env.get('OFFICIAL_PHP_IMAGE', 'localhost/gophp-dependency-research:8.5.11')
        podman_base = ['podman', 'run', '--rm', '--pull=never', '-v', f'{ROOT}:{ROOT}', '-w', str(ROOT)]
        for key, value in overrides.items():
            podman_base += ['-e', f'{key}={value}']
        if args.official:
            engine = [*podman_base, image, 'php']
            identity = {'image': image}
        else:
            binary = Path(env.get('GOPHP_BIN', str(Path(env.get('GOPHPENV_ROOT', str(Path.home()/'.gophpenv'))) / 'shims/gophp'))).expanduser()
            if not binary.is_absolute():
                binary = ROOT / binary
            if not binary.is_file():
                print(f'Missing gophp: {binary}', file=sys.stderr)
                return 2
            engine = [str(binary), '--php-version=8.5']
            identity = {'launcher': str(binary), 'launcher_sha256': hashlib.sha256(binary.read_bytes()).hexdigest()}
        (logs / 'environment.json').write_text(json.dumps({
            'runtime': kind, **identity, 'composer_lock_sha256': hashlib.sha256((ROOT/'composer.lock').read_bytes()).hexdigest()
        }, indent=2) + '\n')

        def run(label, command):
            try:
                r = subprocess.run(command, cwd=ROOT, env=env, capture_output=True, timeout=120)
                code, output = r.returncode, r.stdout + b'\nstderr:\n' + r.stderr
            except subprocess.TimeoutExpired as error:
                code, output = 124, (error.stdout or b'') + b'\nTIMEOUT\n' + (error.stderr or b'')
            (logs / f'{label}.log').write_bytes(output + f'\nexit_code: {code}\n'.encode())
            results.append({'check': label, 'ok': code == 0, 'exit_code': code, 'command': command})

        probe = tmp / 'autoload.php'
        probe.write_text('<?php require ' + json.dumps(str(ROOT/'vendor/autoload.php')) + '; echo "autoload-ok\\n";')
        run('autoload', [*engine, str(probe)])
        for label, words in [
            ('package-discover', ['package:discover']), ('list', ['list']),
            ('about', ['about']), ('routes', ['route:list']),
            ('migrate', ['migrate', '--force']), ('migrate-again', ['migrate', '--force']),
        ]:
            run(label, [*engine, str(ROOT/'artisan'), *words])
        with sqlite3.connect(db) as conn:
            tables = sorted(row[0] for row in conn.execute("SELECT name FROM sqlite_master WHERE type='table'"))
        results.append({'check': 'database-tables', 'ok': {'users', 'cache', 'jobs', 'migrations'}.issubset(tables), 'tables': tables})
        with socket.socket() as sock:
            sock.bind(('127.0.0.1', 0))
            port = sock.getsockname()[1]
        base = f'http://127.0.0.1:{port}'

        def fetch(path):
            try:
                with opener.open(base + path, timeout=15) as response:
                    return response.status, response.read()
            except urllib.error.HTTPError as error:
                return error.code, error.read()

        if args.official:
            command = [*podman_base, '-w', str(ROOT/'public'), '--name', container_name, '-p', f'127.0.0.1:{port}:8000', image,
                       'php', '-S', '0.0.0.0:8000', '-t', str(ROOT/'public'),
                       str(ROOT/'vendor/laravel/framework/src/Illuminate/Foundation/resources/server.php')]
        else:
            command = [*engine, f'--serve=127.0.0.1:{port}', f'--document-root={ROOT / "public"}', str(ROOT/'public/index.php')]
        with (logs/'server.log').open('wb') as output:
            try:
                server = subprocess.Popen(command, cwd=ROOT, env=env, stdout=output, stderr=output)
                ready = False
                for _ in range(100):
                    if server.poll() is not None:
                        break
                    try:
                        fetch('/up')
                        ready = True
                        break
                    except (OSError, urllib.error.URLError):
                        time.sleep(0.1)
                results.append({'check': 'tcp-listening', 'ok': ready})
                if ready:
                    home_ok = False
                    for label, path, expected in [('health', '/up', 200), ('welcome-cold', '/', 200),
                                                  ('welcome-warm', '/', 200), ('not-found', '/not-a-real-route-gophp', 404)]:
                        status, body = fetch(path)
                        ok = status == expected and b'Fatal error' not in body and b'<b>Warning</b>' not in body
                        if path == '/':
                            ok = ok and b'Laravel' in body and b'<html' in body.lower()
                            home_ok = ok
                        (logs/f'{label}.body').write_bytes(body)
                        results.append({'check': label, 'ok': ok, 'status': status})
                    status, body = fetch('/.env')
                    results.append({'check': 'env-not-exposed', 'ok': status in (403, 404) and b'APP_KEY=' not in body, 'status': status})
                    manifest_path = ROOT/'public/build/manifest.json'
                    results.append({'check': 'asset-manifest', 'ok': manifest_path.is_file()})
                    manifest = json.loads(manifest_path.read_text()) if manifest_path.is_file() else {}
                    assets = sorted({entry['file'] for entry in manifest.values()} | {css for entry in manifest.values() for css in entry.get('css', [])})
                    for asset in assets:
                        status, body = fetch('/build/' + asset)
                        expected = (ROOT/'public/build'/asset).read_bytes()
                        results.append({'check': 'asset:'+asset, 'ok': status == 200 and body == expected, 'status': status})
                    if home_ok:
                        checks = [fetch('/') for _ in range(20)]
                        with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
                            checks += list(pool.map(fetch, ['/']*4))
                        results.append({'check': '20-sequential-4-parallel', 'ok': all(s == 200 and b'Laravel' in b for s, b in checks)})
                    else:
                        results.append({'check': '20-sequential-4-parallel', 'ok': False, 'skipped': 'welcome did not pass'})
            finally:
                if args.official:
                    subprocess.run(['podman', 'stop', '--time', '2', container_name], capture_output=True, timeout=20)
                if server is not None:
                    if server.poll() is None:
                        server.terminate()
                    try:
                        server.wait(timeout=10)
                    except subprocess.TimeoutExpired:
                        server.kill()
                        server.wait()
        code = 0 if all(x['ok'] for x in results) else 1
        (logs/'results.json').write_text(json.dumps({'runtime': kind, 'checks': results, 'exit_code': code}, indent=2)+'\n')
        for result in results:
            print(('PASS' if result['ok'] else 'FAIL'), result['check'])
        print(f'logs: {logs}\nexit_code: {code}')
        return code


if __name__ == '__main__':
    sys.exit(main())
