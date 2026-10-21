#!/usr/bin/env python3
"""Check upstream skeleton bytes and protected local files without changing them."""
import hashlib
import io
import json
import subprocess
import sys
import tempfile
import urllib.request
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
APP = Path(__file__).resolve().parent / 'app'
UPSTREAM = '945f4e5a9fd3695dc0ee512f497c650fb82cfbb8'

def digest(p):
    return hashlib.sha256(p.read_bytes()).hexdigest()

def main():
    out = Path(tempfile.mkdtemp(prefix='fresh-laravel-audit-', dir=ROOT / 'docs/validation'))
    result = {'command': sys.argv, 'exit_code': 1, 'upstream_commit': UPSTREAM}
    try:
        url = 'https://codeload.github.com/laravel/laravel/zip/' + UPSTREAM
        archive = urllib.request.urlopen(url, timeout=60).read()
        result['archive_url'] = url
        result['archive_sha256'] = hashlib.sha256(archive).hexdigest()
        z = zipfile.ZipFile(io.BytesIO(archive))
        source = {}
        for name in z.namelist():
            if name.endswith('/'):
                continue
            rel = name.split('/', 1)[1]
            if rel.startswith(('.github/', '.gitattributes')):
                continue
            assert (APP / rel).read_bytes() == z.read(name), 'modified skeleton: ' + rel
            source[rel] = digest(APP / rel)
        result['unchanged_upstream_files'] = source
        baseline = ROOT / '.cache/fresh-laravel/protected.json'
        if baseline.exists():
            protected = json.loads(baseline.read_text())
            changed = [name for name, expected in protected.items() if not (ROOT / name).exists() or digest(ROOT / name) != expected]
            result['protected_files_checked'] = len(protected)
            result['protected_files_changed'] = changed
            assert not changed, 'protected lock/vendor/manual DB changed'
        tracked = subprocess.check_output(['git', 'ls-files', '-z', '--cached', '--others', '--exclude-standard'], cwd=ROOT).decode().split('\0')
        fresh = [p for p in tracked if p.startswith('research/fresh-laravel/')]
        forbidden = [p for p in fresh if '/vendor/' in p or p.endswith(('.sqlite', '.sqlite-wal', '.sqlite-shm')) or Path(p).name == '.env' or '/bootstrap/cache/' in p and not p.endswith('.gitignore') or '/storage/' in p and not p.endswith('.gitignore')]
        result['unignored_artifacts'] = forbidden
        assert not forbidden, 'generated artifact is not ignored'
        result['exit_code'] = 0
        print('unchanged skeleton:', len(source), 'protected:', result.get('protected_files_checked'), 'artifact exclusions: OK')
    finally:
        (out / 'result.json').write_text(json.dumps(result, indent=2) + '\n')
        print('evidence:', out)
    return 0

if __name__ == '__main__':
    raise SystemExit(main())
