#!/usr/bin/env python3
"""Reinstall the unchanged skeleton's fixed lock, including development packages.

The output is a new isolated app, not an update of either existing application.
Composer runs with official PHP; scripts are disabled so they cannot count as
Go application execution. The printed --app path is accepted by verify.py.
"""
import argparse
import hashlib
import json
import os
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
APP = Path(__file__).resolve().parent / 'app'

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--engine', choices=['podman', 'docker'], default='podman')
    args = parser.parse_args()
    (ROOT / '.cache/fresh-laravel').mkdir(parents=True, exist_ok=True)
    work = Path(tempfile.mkdtemp(prefix='rebuild-', dir=ROOT / '.cache/fresh-laravel'))
    app = work / 'app'
    def ignore(path, names):
        rel = Path(path).relative_to(APP)
        if rel == Path('.'):
            return [n for n in names if n in {'vendor', '.env', '.git', 'node_modules'}]
        if rel == Path('bootstrap/cache') or str(rel).startswith('storage/'):
            return [n for n in names if n != '.gitignore']
        return [n for n in names if n.endswith('.sqlite') or n.endswith('.log')]
    shutil.copytree(APP, app, ignore=ignore)
    out = Path(tempfile.mkdtemp(prefix='fresh-laravel-rebuild-', dir=ROOT / 'docs/validation'))
    commands = []
    result = 0
    try:
        for version in ['8.3.35', '8.4.26', '8.5.11']:
            operations = [['composer', 'check-platform-reqs', '--lock', '--format=json'], ['php', '-m'], ['php', '-r', 'echo PHP_VERSION,"\\n";']]
            if version == '8.3.35':
                operations.insert(0, ['composer', 'install', '--no-scripts', '--no-interaction', '--prefer-dist'])
            for i, operation in enumerate(operations):
                user_args = ['--user', f'{os.getuid()}:{os.getgid()}'] if args.engine == 'docker' else []
                cmd = [args.engine, 'run', '--rm', *user_args, '-e', 'COMPOSER_HOME=/tmp/composer', '-v', f'{ROOT}:{ROOT}', '-w', str(app), 'localhost/gophp-dependency-research:' + version, *operation]
                r = subprocess.run(cmd, capture_output=True)
                label = version + '-' + str(i)
                (out / (label + '.stdout')).write_bytes(r.stdout)
                (out / (label + '.stderr')).write_bytes(r.stderr)
                commands.append({'command': cmd, 'version': version, 'exit_code': r.returncode})
                print(version, operation, 'exit', r.returncode, flush=True)
                if r.returncode:
                    raise RuntimeError('dependency/platform command failed')
        assert (app / 'composer.lock').read_bytes() == (APP / 'composer.lock').read_bytes()
        lock = json.loads((app / 'composer.lock').read_text())
        installed = json.loads((app / 'vendor/composer/installed.json').read_text())['packages']
        expected = {(p['name'], p['version']) for p in lock['packages'] + lock['packages-dev']}
        assert expected == {(p['name'], p['version']) for p in installed}
        (out / 'identity.json').write_text(json.dumps({'app': str(app), 'lock_sha256': hashlib.sha256((app / 'composer.lock').read_bytes()).hexdigest(), 'packages': len(installed), 'dev_packages': len(lock['packages-dev']), 'framework': dict(expected)['laravel/framework']}, indent=2) + '\n')
        (ROOT / '.cache/fresh-laravel/rebuilt-app.txt').write_text(str(app) + '\n')
        print('Rebuilt app:', app, '\nEvidence:', out)
    except Exception as error:
        print(repr(error), file=sys.stderr)
        result = 1
    finally:
        (out / 'commands.json').write_text(json.dumps(commands, indent=2) + '\n')
        (out / 'result.json').write_text(json.dumps({'command': sys.argv, 'exit_code': result}, indent=2) + '\n')
    return result

if __name__ == '__main__':
    raise SystemExit(main())
