#!/usr/bin/env python3
"""Run the existing differential runner, retaining each version's real exit code."""
import json
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / 'scripts'))
from release_evidence import source_identity

def main():
    out = Path(tempfile.mkdtemp(prefix='fresh-laravel-fixture-matrix-', dir=ROOT / 'docs/validation'))
    (out / 'source.json').write_text(json.dumps(source_identity(ROOT), indent=2) + '\n')
    commands = []
    for target, patch in [('8.3', '8.3.35'), ('8.4', '8.4.26'), ('8.5', '8.5.11')]:
        cmd = [sys.executable, str(ROOT / 'skills/php-go-interpreter/scripts/compare_php.py'), '--php', str(ROOT / 'containers/php-version/bin' / ('php' + target.replace('.', ''))), '--php-version', patch, '--fixtures', 'tests/fixtures', '--timeout', '15', '--interpreter', './bin/gophp', '--php-version=' + target, '{file}']
        wrapped = [sys.executable, str(Path(__file__).with_name('record.py')), '--label', 'final-fixtures-' + patch, '--', *cmd]
        result = subprocess.run(wrapped, cwd=ROOT)
        commands.append({'command': cmd, 'version': patch, 'exit_code': result.returncode})
        (out / 'commands.json').write_text(json.dumps(commands, indent=2) + '\n')
    code = max(item['exit_code'] for item in commands)
    (out / 'result.json').write_text(json.dumps({'command': sys.argv, 'exit_code': code}, indent=2) + '\n')
    return code

if __name__ == '__main__':
    raise SystemExit(main())
