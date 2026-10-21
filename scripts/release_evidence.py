"""Reproducible source and executable identities for local and CI evidence."""
import hashlib
import json
import os
from pathlib import Path
import subprocess


def digest(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()


def source_identity(root):
    root = Path(root)
    def git(*args):
        return subprocess.check_output(['git', '-C', str(root), *args])
    names = sorted(set(git('ls-files', '-z', '--cached', '--others', '--exclude-standard').decode().split('\0')) - {''})
    manifest = []
    for name in names:
        path = root/name
        # Evidence and prose may be updated after a successful run. Record the
        # executable/test/config source independently so this is explicit.
        if name.startswith('docs/') or name.endswith('.md') or (not path.exists() and not path.is_symlink()):
            continue
        manifest.append({'path':name, 'sha256':digest(path) if not path.is_symlink() else hashlib.sha256(os.readlink(path).encode()).hexdigest(),
                         'mode':'120000' if path.is_symlink() else ('100755' if path.stat().st_mode & 0o111 else '100644'), 'symlink':path.is_symlink()})
    encoded = json.dumps(manifest, sort_keys=True, separators=(',', ':')).encode()
    binary = root/'bin/gophp'
    return {'head':git('rev-parse','HEAD').decode().strip(),
            'source_sha256':hashlib.sha256(encoded).hexdigest(),
            'diff_sha256':hashlib.sha256(git('diff','HEAD','--binary')).hexdigest(),
            'binary_sha256':digest(binary) if binary.exists() else None,
            'locks':{app:digest(root/'research'/app/'composer.lock') for app in ('laravel-api','task-api')},
            'manifest':manifest}
