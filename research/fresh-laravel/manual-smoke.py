#!/usr/bin/env python3
"""Check documented key generation/package discovery without touching a manual DB."""
import base64
import os
import shutil
import subprocess
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
APP = Path(__file__).resolve().parent / 'app'

def main():
    with tempfile.TemporaryDirectory(prefix='manual-smoke-', dir=ROOT / '.cache/fresh-laravel') as temp:
        app = Path(temp) / 'app'
        shutil.copytree(APP, app, ignore=shutil.ignore_patterns('vendor', '.env', '*.sqlite'))
        (app / 'vendor').symlink_to(APP / 'vendor', target_is_directory=True)
        for path in (app / 'bootstrap/cache').glob('*.php'):
            path.unlink()
        shutil.copyfile(app / '.env.example', app / '.env')
        for command in ['key:generate', 'package:discover']:
            cmd = [str(ROOT / 'bin/gophp'), '--php-version=8.4', 'artisan', command, '--no-ansi']
            print('command:', cmd, 'cwd:', app, flush=True)
            run = subprocess.run(cmd, cwd=app, capture_output=True, text=True)
            print(run.stdout, run.stderr, 'exit_code:', run.returncode, flush=True)
            assert run.returncode == 0 and not run.stderr
        key = next(line.split('=', 1)[1] for line in (app / '.env').read_text().splitlines() if line.startswith('APP_KEY='))
        assert key.startswith('base64:') and len(base64.b64decode(key[7:], validate=True)) == 32
        print('PASS generated 32-byte key in disposable .env; package discovery passed')

if __name__ == '__main__':
    main()
