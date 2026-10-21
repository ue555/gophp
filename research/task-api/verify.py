#!/usr/bin/env python3
"""Real TCP / SQLite contract verification, shared by official PHP and Go PHP."""
import argparse
from concurrent.futures import ThreadPoolExecutor
from contextlib import redirect_stdout, redirect_stderr
from datetime import datetime, timezone
import hashlib
import json
import os
from pathlib import Path
import re
import shlex
import shutil
import socket
import sqlite3
import subprocess
import sys
import tempfile
import time
import traceback
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen

APP = Path(__file__).resolve().parent
ROOT = APP.parents[1]
VERSIONS = {'8.3': '8.3.35', '8.4': '8.4.26', '8.5': '8.5.11'}


def check(condition, message):
    if not condition:
        raise RuntimeError(message)


def command(args, env=None, cwd=APP):
    args = list(map(str, args))
    print('command:', shlex.join(args), flush=True)
    result = subprocess.run(args, cwd=cwd, env=env, capture_output=True, text=True)
    print(result.stdout, end='')
    print(result.stderr, end='', file=sys.stderr)
    print('exit_code:', result.returncode, flush=True)
    check(result.returncode == 0, f'command failed: {result.returncode}')
    return result.stdout


def request(port, method, path, body=None):
    headers = {'Accept': 'application/json'}
    data = None if body is None else json.dumps(body, ensure_ascii=False).encode()
    if data is not None:
        headers['Content-Type'] = 'application/json'
    req = Request(f'http://127.0.0.1:{port}{path}', data=data, headers=headers, method=method)
    try:
        response = urlopen(req, timeout=60)
    except HTTPError as error:
        response = error
    with response:
        raw = response.read().decode('utf-8')
        status = response.status
        content_type = response.headers.get('Content-Type')
    # Preserve the raw body; warnings/extra output cause JSON parsing to fail.
    return {'method': method, 'path': path, 'request': body, 'status': status,
            'content_type': content_type, 'raw': raw}


def rows(database):
    with sqlite3.connect(database) as connection:
        connection.row_factory = sqlite3.Row
        return [dict(row) for row in connection.execute('SELECT id,title,completed,created_at,updated_at FROM tasks ORDER BY id')]


def timestamp(value):
    check(isinstance(value, str) and re.fullmatch(r'\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{6}Z', value), f'invalid timestamp: {value!r}')
    parsed = datetime.strptime(value, '%Y-%m-%dT%H:%M:%S.%fZ').replace(tzinfo=timezone.utc)
    check(abs((datetime.now(timezone.utc) - parsed).total_seconds()) < 600, 'timestamp not near current time')
    return parsed


def task(value, title, completed, task_id=None):
    check(isinstance(value, dict) and set(value) == {'id','title','completed','created_at','updated_at'}, f'task fields: {value!r}')
    check(type(value['id']) is int and value['id'] > 0, 'invalid id')
    check(task_id is None or value['id'] == task_id, 'id changed')
    check(value['title'] == title and value['completed'] is completed, f'task differs: {value!r}')
    check(timestamp(value['created_at']) <= timestamp(value['updated_at']), 'timestamps out of order')
    return value


def scenario(port, database, evidence):
    records = []
    expected = {}

    def db_check():
        actual = rows(database)
        wanted = []
        for value in sorted(expected.values(), key=lambda v: v['id']):
            row = dict(value)
            row['completed'] = int(row['completed'])
            for key in ('created_at', 'updated_at'):
                row[key] = timestamp(row[key]).strftime('%Y-%m-%d %H:%M:%S')
            wanted.append(row)
        check(actual == wanted, f'DB mismatch: actual={actual!r} expected={wanted!r}')
        print('DB verified:', json.dumps(actual, ensure_ascii=False), flush=True)
        return actual

    def accept(record, status):
        records.append(record)
        # Write each result before checking, so failed responses remain evidence.
        evidence.write(json.dumps(record, ensure_ascii=False)+'\n'); evidence.flush()
        print(f"HTTP {record['method']} {record['path']}: {record['status']} {record['content_type']} {record['raw']}", flush=True)
        check(record['status'] == status, f'expected {status}: {record!r}')
        if status == 204:
            check(record['raw'] == '' and record['content_type'] is None, '204 must have no body or Content-Type')
            value = None
        else:
            check(record['content_type'] and record['content_type'].split(';')[0] == 'application/json', 'non-JSON Content-Type')
            value = json.loads(record['raw'])
            if status in (200, 201):
                for item in value if isinstance(value, list) else [value]:
                    check(set(item) == {'id','title','completed','created_at','updated_at'}, 'Task fields differ')
                    check(type(item['id']) is int and type(item['title']) is str and type(item['completed']) is bool,
                          'Task JSON types differ')
        record['json'] = value
        return value

    def call(method, path, body=None, status=200):
        value = accept(request(port, method, path, body), status)
        return value

    def unchanged(method, path, body, status, field=None):
        before = rows(database)
        value = call(method, path, body, status)
        if field:
            check(set(value) == {'message','errors'} and isinstance(value['message'], str), 'validation error structure')
            check(field in value['errors'] and all(isinstance(m, str) for m in value['errors'][field]), 'missing validation error')
        check(rows(database) == before, 'failed/no-op request changed DB')
        print('DB unchanged:', json.dumps(before, ensure_ascii=False), flush=True)
        return value

    check(call('GET', '/api/tasks') == [], 'initial list must be []'); db_check()
    first = task(call('POST', '/api/tasks', {'title':'日本語のタスク'}, 201), '日本語のタスク', False)
    check(first['created_at'] == first['updated_at'], 'create timestamps differ')
    expected[first['id']] = first; db_check()
    second = task(call('POST', '/api/tasks', {'title':'二つ目', 'completed':True}, 201), '二つ目', True)
    check(second['id'] > first['id'], 'ids do not increase')
    expected[second['id']] = second; db_check()
    check(call('GET', '/api/tasks') == list(expected.values()), 'list/order differs')
    path = f"/api/tasks/{first['id']}"
    check(call('GET', path) == first, 'item differs')
    time.sleep(1.1)  # SQLite stores seconds; ensure an actual update changes updated_at.
    updated = task(call('PATCH', path, {'title':'更新した日本語'}), '更新した日本語', False, first['id'])
    check(updated['created_at'] == first['created_at'] and updated['updated_at'] > first['updated_at'], 'update timestamp did not advance')
    expected[first['id']] = updated; db_check()
    updated2 = task(call('PATCH', path, {'completed':True}), '更新した日本語', True, first['id'])
    check(updated2['created_at'] == first['created_at'] and updated2['updated_at'] >= updated['updated_at'], 'completed update timestamps')
    expected[first['id']] = updated2; db_check()
    check(unchanged('PATCH', path, {}, 200) == updated2, 'empty PATCH differs')
    check(unchanged('PATCH', path, None, 200) == updated2, 'bodyless PATCH differs')
    check(unchanged('PATCH', path, {'ignored':'value','id':999,'created_at':'bad'}, 200) == updated2, 'unknown fields not ignored')
    check(unchanged('PATCH', path, {'title':updated2['title'], 'completed':True}, 200) == updated2, 'unchanged PATCH differs')

    # Both create and partial update cover each invalid field; no write on failure.
    unchanged('POST', '/api/tasks', {}, 422, 'title')
    for title in ['', '   ', 'あ'*121, 'a'*121, None, 123, True, [], {'x':1}]:
        unchanged('POST', '/api/tasks', {'title':title}, 422, 'title')
        unchanged('PATCH', path, {'title':title,'completed':False}, 422, 'title')
    for completed in [None, '', 'true', 'false', 'yes', 2, -1, 0.5, [], {}]:
        unchanged('POST', '/api/tasks', {'title':'invalid','completed':completed}, 422, 'completed')
        unchanged('PATCH', path, {'title':'must not persist','completed':completed}, 422, 'completed')
    # Accepted boolean representations and the Unicode character-count boundary.
    for index, completed in enumerate([True, False, 0, 1, '0', '1']):
        title = 'あ'*120 if index == 0 else f'boolean-{index}'
        value = task(call('POST', '/api/tasks', {'title':title, 'completed':completed, 'id':999}, 201), title, bool(int(completed)))
        check(value['id'] != 999, 'unvalidated id assigned')
        expected[value['id']] = value; db_check()
        value = task(call('PATCH', f"/api/tasks/{value['id']}", {'completed':completed}), title, bool(int(completed)), value['id'])
        expected[value['id']] = value; db_check()
    # Default middleware trims strings and turns empty strings into null.
    trimmed = task(call('POST', '/api/tasks', {'title':'  trim 日本語  '}, 201), 'trim 日本語', False)
    expected[trimmed['id']] = trimmed; db_check()
    for method, body in [('GET',None), ('PATCH',{'title':'missing'}), ('PATCH',{}), ('DELETE',None)]:
        value = unchanged(method, '/api/tasks/999999', body, 404)
        check(isinstance(value.get('message'), str), '404 structure')
    db_check()
    # Alternate statuses/bodies/methods/query parameters in one long-lived server.
    check(call('GET', path+'?title=leak&completed=0') == updated2, 'request input leaked')
    unchanged('POST', '/api/tasks', {'completed':True}, 422, 'title')
    check(call('GET', '/api/tasks') == list(expected.values()), 'validation state leaked')
    call('DELETE', path, status=204); del expected[first['id']]; db_check()
    unchanged('GET', path, None, 404)
    unchanged('DELETE', path, None, 404)
    check(call('GET', '/api/tasks') == list(expected.values()), 'list after delete differs')
    for index in range(100):
        check(call('GET', f'/api/tasks?round={index}') == list(expected.values()), f'sequential GET {index}')
    with ThreadPoolExecutor(max_workers=8) as pool:
        responses = list(pool.map(lambda index: request(port, 'GET', f'/api/tasks?parallel={index}'), range(8)))
    for record in responses:
        check(accept(record, 200) == list(expected.values()), 'parallel GET differs')
    db_check()
    for task_id in list(expected):
        call('DELETE', f'/api/tasks/{task_id}', status=204); del expected[task_id]; db_check()
    check(call('GET', '/api/tasks') == [], 'final empty list differs')
    for record in records:
        # JSON key order and transport whitespace are not meaningful; only these
        # two timestamp fields are replaced after all format/consistency checks.
        record.pop('raw')
        record['json'] = normalize(record['json'])
    print(f'PASS requests={len(records)} sequential_GET=100 parallel_GET=8 SQLite verified', flush=True)
    return records


def normalize(value):
    if isinstance(value, list):
        return [normalize(v) for v in value]
    if isinstance(value, dict):
        return {k: '<validated UTC timestamp>' if k in ('created_at','updated_at') else normalize(v) for k,v in value.items()}
    return value


def verify(runtime, version, executable, log_dir, engine='podman', image_template='localhost/gophp-dependency-research:{version}'):
    target = f'{runtime}-{version}'
    print('target:', target, 'official_patch:', VERSIONS[version], flush=True)
    with tempfile.TemporaryDirectory(prefix=f'task-api-{target}-', dir=APP / '.cache') as temporary:
        temporary = Path(temporary)
        database = temporary / 'database.sqlite'; database.touch()
        for directory in ('storage/framework/cache/data','storage/framework/views','storage/framework/sessions','storage/logs','bootstrap-cache'):
            (temporary / directory).mkdir(parents=True)
        base = '/probe' if runtime == 'official' else str(temporary)
        settings = {'APP_ENV':'testing', 'APP_DEBUG':'false', 'DB_CONNECTION':'sqlite',
                    'DB_DATABASE':base+'/database.sqlite', 'LARAVEL_STORAGE_PATH':base+'/storage'}
        for key, name in [('PACKAGES','packages.php'),('SERVICES','services.php'),('CONFIG','config.php'),('ROUTES','routes.php'),('EVENTS','events.php')]:
            settings[f'APP_{key}_CACHE'] = base+'/bootstrap-cache/'+name
        print('environment:', json.dumps(settings, sort_keys=True), flush=True)
        env = os.environ.copy(); env.update(settings)
        image = image_template.format(version=VERSIONS[version])
        container_args = ['-v',f'{APP}:/app:ro','-v',f'{temporary}:/probe']
        if engine == 'docker':
            container_args += ['--user', f'{os.getuid()}:{os.getgid()}', '-e', 'COMPOSER_HOME=/tmp/composer']
        for key,value in settings.items(): container_args += ['-e', f'{key}={value}']
        prefix = ([engine,'run','--rm',*container_args,image,'php'] if runtime == 'official'
                  else [str(executable), f'--php-version={version}'])
        if runtime == 'official':
            actual = command(prefix+['-r','echo PHP_VERSION, PHP_EOL;'])
            check(actual == VERSIONS[version]+'\n', f'wrong official version: {actual!r}')
            command([engine,'run','--rm',*container_args,image,'composer','check-platform-reqs','--no-dev'])
        command(prefix+['artisan','package:discover','--no-interaction'], env)
        command(prefix+['artisan','migrate','--force','--no-interaction'], env)
        check(rows(database) == [], 'migration did not create empty tasks table')
        # Clear only this run's schema metadata/cache after migration.
        shutil.rmtree(temporary/'storage/framework/cache/data')
        (temporary/'storage/framework/cache/data').mkdir()
        command(prefix+['artisan','migrate:status','--no-interaction'], env)
        container = None; process = None
        server_path = log_dir/f'{target}-server.log'
        server_log = server_path.open('x+')
        try:
            if runtime == 'official':
                container = command([engine,'run','-d','--rm','-p','127.0.0.1::8080',*container_args,image,
                                     'php','-S','0.0.0.0:8080','-t','public','public/index.php']).strip()
                mapping = command([engine,'port',container,'8080/tcp'])
                port = int(mapping.rsplit(':',1)[1])
            else:
                with socket.socket() as listener:
                    listener.bind(('127.0.0.1',0)); port = listener.getsockname()[1]
                args = prefix+[f'--serve=127.0.0.1:{port}', 'public/index.php']
                print('command:', shlex.join(args), flush=True)
                process = subprocess.Popen(args, cwd=APP, env=env, stdout=server_log, stderr=subprocess.STDOUT)
            # TCP readiness only: no extra PHP request outside the shared scenario.
            for _ in range(150):
                if process and process.poll() is not None:
                    raise RuntimeError(f'server exited: {process.returncode}')
                try:
                    with socket.create_connection(('127.0.0.1',port), timeout=1): pass
                    break
                except OSError: time.sleep(0.2)
            else: raise RuntimeError('TCP server did not become ready')
            with (log_dir/f'{target}-http.jsonl').open('x') as evidence:
                result = scenario(port, database, evidence)
            (log_dir/f'{target}-normalized.json').write_text(json.dumps(result, ensure_ascii=False, indent=2)+'\n')
            print('final_DB:', json.dumps(rows(database)), flush=True)
            return result
        finally:
            if process:
                process.terminate()
                try: process.wait(timeout=10)
                except subprocess.TimeoutExpired: process.kill(); process.wait()
                print('server_exit_code:', process.returncode, '(terminated by verifier)', flush=True)
            if container:
                try:
                    # Keep the exact server output as a separate artifact.
                    print('command:', shlex.join([engine,'logs',container]), flush=True)
                    r = subprocess.run([engine,'logs',container], stdout=server_log, stderr=subprocess.STDOUT)
                    print('logs_exit_code:', r.returncode, flush=True)
                    check(r.returncode == 0, 'could not collect server logs')
                    inspection = command([engine,'inspect','--format','{{.State.Running}}',container])
                    check(inspection == 'true\n', 'official server stopped unexpectedly')
                finally:
                    command([engine,'stop',container])
            server_log.close()
            print('server_log:', server_path, flush=True)
            output = server_path.read_text()
            print(output, end='', flush=True)
            check(not re.search(r'PHP (Warning|Fatal error|Deprecated|Notice)|Uncaught |panic:|testing.ERROR:', output), 'server emitted PHP diagnostics')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--runtime', choices=('all','go','official'), default='all')
    parser.add_argument('--version', choices=VERSIONS, action='append')
    parser.add_argument('--interpreter', type=Path, default=ROOT/'bin/gophp')
    parser.add_argument('--engine', choices=('podman','docker'), default='podman')
    parser.add_argument('--image-template', default='localhost/gophp-dependency-research:{version}')
    parser.add_argument('--output-root', type=Path, default=ROOT/'docs/validation')
    args = parser.parse_args()
    (APP/'.cache').mkdir(exist_ok=True)
    args.output_root.mkdir(parents=True, exist_ok=True)
    log_dir = Path(tempfile.mkdtemp(prefix='task-api-'+datetime.now().strftime('%Y%m%d-%H%M%S')+'-', dir=args.output_root))
    print('logs:', log_dir, flush=True)
    result = 0
    with (log_dir/'run.log').open('x', buffering=1) as log, redirect_stdout(log), redirect_stderr(log):
        print('command:', shlex.join([sys.executable,*sys.argv]), flush=True)
        sys.path.insert(0, str(ROOT/'scripts'))
        from release_evidence import source_identity
        (log_dir/'source.json').write_text(json.dumps(source_identity(ROOT), indent=2)+'\n')
        print('interpreter_sha256:', hashlib.sha256(args.interpreter.read_bytes()).hexdigest() if args.interpreter.is_file() else 'missing')
        try:
            check((APP/'vendor/autoload.php').is_file(), 'install dependencies first')
            for version in args.version or VERSIONS:
                official = verify('official',version,args.interpreter.resolve(),log_dir,args.engine,args.image_template) if args.runtime in ('all','official') else None
                go = verify('go',version,args.interpreter.resolve(),log_dir,args.engine,args.image_template) if args.runtime in ('all','go') else None
                if args.runtime == 'all':
                    check(json.dumps(official, sort_keys=True) == json.dumps(go, sort_keys=True),
                          f'HTTP semantic comparison differs: {version}; inspect normalized JSON')
                    print(f'PASS comparison PHP {VERSIONS[version]} / Go {version}', flush=True)
        except PermissionError:
            traceback.print_exc(); result = 2
        except Exception:
            traceback.print_exc(); result = 1
        print('exit_code:', result, flush=True)
    print('result:', result, 'log:', log_dir/'run.log', flush=True)
    return result


if __name__ == '__main__':
    raise SystemExit(main())
