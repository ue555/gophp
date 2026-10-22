"""Check the actual daily launchers with the app's normal .env, isolated DB/cache."""
import hashlib
import json
import os
from pathlib import Path
import socket
import subprocess
import tempfile
import time
import urllib.request
from record import run, ROOT as LOGS

LOGS = Path(tempfile.mkdtemp(prefix='daily-launcher-', dir=LOGS))
app = Path('/Users/kouji/dev/gophp-app')
env = dict(os.environ, GOPHP_BIN='/Users/kouji/dev/php/.cache/laravel13/gophp')
with tempfile.TemporaryDirectory(prefix='launcher-', dir=app/'.cache') as tmp:
    tmp = Path(tmp)
    (tmp/'database.sqlite').touch()
    (tmp/'views').mkdir()
    env.update(DB_CONNECTION='sqlite', DB_DATABASE=str(tmp/'database.sqlite'),
               VIEW_COMPILED_PATH=str(tmp/'views'), LOG_CHANNEL='stderr')
    for key, name in [('APP_CONFIG_CACHE', 'config'), ('APP_SERVICES_CACHE', 'services'),
                      ('APP_PACKAGES_CACHE', 'packages'), ('APP_ROUTES_CACHE', 'routes'),
                      ('APP_EVENTS_CACHE', 'events')]:
        env[key] = str(tmp/(name+'.php'))
    for command in [('about',), ('migrate', '--force')]:
        result = run('launcher-'+command[0], ['scripts/artisan', *command], cwd=app, env=env)
        assert result.returncode == 0 and not result.stderr
        if command[0] == 'about':
            assert b'local' in result.stdout and b'13.38.0' in result.stdout
    with socket.socket() as sock:
        sock.bind(('127.0.0.1', 0))
        port = sock.getsockname()[1]
    env['PORT'] = str(port)
    with tempfile.NamedTemporaryFile(prefix='launcher-server-', suffix='.log', dir=LOGS) as output:
        server = subprocess.Popen(['scripts/serve'], cwd=app, env=env, stdout=output, stderr=output)
        evidence = {'command': ['scripts/serve'], 'cwd': str(app), 'port': port, 'http': []}
        try:
            for _ in range(200):
                assert server.poll() is None
                try:
                    with socket.create_connection(('127.0.0.1', port), timeout=.2):
                        break
                except OSError:
                    time.sleep(.1)
            opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))
            for i, path in enumerate(['/up', '/', '/', '/build/assets/app-6jSCucsl.css']):
                with opener.open(f'http://127.0.0.1:{port}'+path, timeout=30) as response:
                    body = response.read()
                    assert response.status == 200 and body
                    if path == '/':
                        assert b'Laravel' in body and b'/build/assets/' in body
                    (LOGS/f'launcher-http-{i}.body').write_bytes(body)
                    evidence['http'].append({'path': path, 'status': response.status,
                                             'headers': list(response.headers.items()),
                                             'sha256': hashlib.sha256(body).hexdigest()})
        finally:
            server.terminate()
            server.wait(timeout=15)
            output.seek(0)
            raw = output.read()
            (LOGS/'launcher-server.log').write_bytes(raw)
            evidence['server_exit_code'] = server.returncode
            evidence['terminated_by_verifier'] = True
            (LOGS/'launcher-result.json').write_text(json.dumps(evidence, indent=2)+'\n')
            assert raw == b'', raw
assert not tmp.exists()
print('PASS actual scripts/artisan and scripts/serve; app local environment; temporary data removed')
