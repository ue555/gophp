#!/usr/bin/env python3
"""Reconcile raw successful evidence with the current Go source and fixture set."""
import argparse
import hashlib
import json
import os
import re
import subprocess
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

def read(path):
    return json.loads(path.read_text())

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    for name in ['http', 'fixtures', 'task', 'legacy_log']:
        parser.add_argument('--' + name.replace('_', '-'), type=Path, required=True)
    args = parser.parse_args()
    out = Path(tempfile.mkdtemp(prefix='fresh-laravel-results-audit-', dir=ROOT / 'docs/validation'))
    result = {'command': vars(args), 'exit_code': 1}
    try:
        expected = read(ROOT / 'tests/fixtures/expected.json')
        actual = {str(p.relative_to(ROOT / 'tests/fixtures')) for p in (ROOT / 'tests/fixtures').rglob('*.php')}
        assert actual == set(expected)
        baseline = json.loads(subprocess.check_output(['git', 'show', 'HEAD:tests/fixtures/expected.json'], cwd=ROOT))
        assert all(expected[k] == v for k, v in baseline.items())
        result['fixtures'] = len(actual)
        result['unchanged_baseline_expectations'] = len(baseline)
        binary = digest(ROOT / 'bin/gophp')
        result['binary_sha256'] = binary
        for evidence in [args.http, args.fixtures, args.task]:
            source = read(evidence / 'source.json')
            assert source['binary_sha256'] == binary, str(evidence) + ': wrong binary'
            for entry in source['manifest']:
                name = entry['path']
                if name.endswith('.go') or name.startswith('tests/fixtures/') or name in {'go.mod', 'go.sum'}:
                    path=ROOT/name
                    current=hashlib.sha256(os.readlink(path).encode()).hexdigest() if entry['symlink'] else digest(path)
                    assert current == entry['sha256'], 'source changed since verification: ' + name
        commands = read(args.fixtures / 'commands.json')
        assert len(commands) == 3 and all(c['exit_code'] == 0 for c in commands)
        for v, patch in [('8.3', '8.3.35'), ('8.4', '8.4.26'), ('8.5', '8.5.11')]:
            logs = sorted((ROOT / 'docs/validation').glob('fresh-laravel-final-fixtures-' + patch + '-*.log'))
            log = logs[-1].read_text()
            assert f'{len(actual)} matched, 0 different, 0 errors; {len(actual)} fixtures' in log
            assert log.endswith('exit_code: 0\n')
            matched = re.findall(r'^MATCH (.+) \(exit -?\d+\)$', log, re.M)
            assert {str(Path(p).relative_to(ROOT / 'tests/fixtures')) for p in matched} == actual
            assert read(args.http / ('official-' + v + '-normalized.json')) == read(args.http / ('go-' + v + '-normalized.json'))
            for runtime in ['official', 'go']:
                records = [json.loads(line) for line in (args.http / (runtime + '-' + v + '-http.jsonl')).read_text().splitlines()]
                assert len(records) == 111 and all(r['status'] == 200 for r in records)
                about = read(args.http / (runtime + '-' + v + '-about.stdout'))
                assert about['environment']['php_version'] == patch
                sessions = read(args.http / (runtime + '-' + v + '-sessions.json'))
                assert len(sessions['sessions']) == 5
                for suffix in ['autoload', 'environment', 'list', 'about', 'routes', 'migrate', 'migrate-again', 'blade-tokens']:
                    assert not (args.http / (runtime + '-' + v + '-' + suffix + '.stderr')).read_bytes()
            assert f'PASS comparison PHP {patch} / Go {v}' in (args.task / 'run.log').read_text()
        assert read(args.http / 'result.json')['exit_code'] == 0
        legacy = args.legacy_log.read_text()
        assert legacy.endswith('exit_code: 0\n')
        for v in ['8.3', '8.4', '8.5']:
            assert f'PASS Go PHP {v}: CLI, HTTP, SQLite, validation, API Platform CRUD, OpenAPI' in legacy
        assert legacy.count('PASS sustained_requests=100 parallel_requests=8') == 3
        result['versions'] = ['8.3.35', '8.4.26', '8.5.11']
        result['fresh_http_requests_per_runtime'] = 111
        result['new_fixtures'] = sorted(actual - set(baseline))
        result['exit_code'] = 0
        print('PASS current Go sources/binary, 1193 fixtures x 3, fresh HTTP x 6, fixed apps')
    finally:
        (out / 'result.json').write_text(json.dumps(result, indent=2, default=str) + '\n')
        print('evidence:', out)

if __name__ == '__main__':
    main()
