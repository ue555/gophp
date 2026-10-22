import ast
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[3]
LOGS = Path(__file__).resolve().parent
APP = Path('/Users/kouji/dev/gophp-app')
sys.path.insert(0, str(ROOT/'scripts'))
from release_evidence import source_identity


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


records = []
required = ['full-fixtures-8.3.35', 'full-fixtures-8.4.26', 'full-fixtures-8.5.11',
            'final-make-test', 'final-make-vet', 'final-make-build', 'final-diff-check',
            'fixed-laravel-api-platform-three', 'task-api-three', 'fresh-laravel12-85',
            'go-final-http', 'official-final-http', 'paired-http-and-db', 'daily-launchers', 'original-repro-fixed']
for label in required:
    paths = sorted(LOGS.glob('*-'+label+'-*.json'))
    assert len(paths) == 1, (label, paths)
    path = paths[0]
    record = json.loads(path.read_text())
    assert record['exit_code'] == 0, path
    record.update(label=label, metadata=path.name, stdout=path.with_suffix('.stdout').name,
                  stderr=path.with_suffix('.stderr').name)
    if label.startswith('full-fixtures-'):
        stdout = path.with_suffix('.stdout').read_text()
        assert stdout.count('\nMATCH ') == 1201
        assert stdout.endswith('1201 matched, 0 different, 0 errors; 1201 fixtures\n')
    records.append(record)
targets = sorted(LOGS.glob('*-target-*.json'))
assert len(targets) == 24
for path in targets:
    record = json.loads(path.read_text())
    assert record['exit_code'] == 0
    assert path.with_suffix('.stdout').read_text().endswith('1 matched, 0 different, 0 errors; 1 fixtures\n')
    records.append(dict(record, metadata=path.name))

protected = json.loads((ROOT/'.cache/laravel13/protected.json').read_text())
changed = [name for name, digest in protected.items() if not Path(name).is_file() or sha(Path(name)) != digest]
assert not changed, changed
old = json.loads(subprocess.check_output(['git', 'show', 'HEAD:tests/fixtures/expected.json'], cwd=ROOT))
new = json.loads((ROOT/'tests/fixtures/expected.json').read_text())
assert len(old) == 1193 and len(new) == 1201
assert all(new.get(k) == v for k, v in old.items())
assert len(list((ROOT/'tests/fixtures').rglob('*.php'))) == 1201

tested = json.loads((LOGS/'tested-source.json').read_text())
current = source_identity(ROOT)
assert current['source_sha256'] == tested['source_sha256']
assert sha(Path(tested['candidate_path'])) == tested['candidate_sha256'] == current['binary_sha256']
assert current['head'] == tested['head']
for name, digest in tested['app_scripts'].items():
    assert sha(APP/name) == digest, name
for name in ('verify', 'compare-validation.py'):
    ast.parse((APP/'scripts'/name).read_text())

for dirname in ('20261022-121925-5367ae79-official', '20261022-122544-776c16d7-gophp'):
    folder = APP/'docs/validation'/dirname
    result = json.loads((folder/'results.json').read_text())
    assert result['exit_code'] == 0 and len(result['checks']) == 47
    assert all(check['ok'] for check in result['checks'])
    assert len(list(folder.glob('*.http.json'))) == 46
    exit_info = json.loads((folder/'server-exit.json').read_text())
    assert exit_info['terminated_by_verifier']
    # Temporary paths are recoverable from the saved actual migration argv.
    commands = json.loads((folder/'results.json').read_text())['checks']
    autoload = next(x for x in commands if x['check'] == 'autoload')['command']
    temp = Path(autoload[-1]).parent
    assert not temp.exists(), temp

index = {'result': 'PASS', 'fixtures_per_version': 1201, 'new_fixtures': 8,
         'target_comparisons': len(targets), 'protected_files_unchanged': len(protected),
         'old_expectations_unchanged': len(old), 'head': current['head'],
         'source_sha256': current['source_sha256'], 'candidate_sha256': tested['candidate_sha256'],
         'commands': records}
# A fresh audit log is retained on every run; the index is a derived document.
with tempfile.NamedTemporaryFile(mode='w', prefix='final-audit-', suffix='.json', dir=LOGS, delete=False) as out:
    json.dump(index, out, ensure_ascii=False, indent=2)
    out.write('\n')
    print('audit:', out.name)
(LOGS/'verification-index.json').write_text(json.dumps(index, ensure_ascii=False, indent=2)+'\n')
rows = ['| 検証 | 実測結果 | 生ログ |', '|---|---|---|']
for record in records[:len(required)]:
    label = record['label']
    result = '1201一致・0差異・0エラー、終了0' if label.startswith('full-fixtures-') else '終了0'
    rows.append(f"| {label} | {result} | [stdout]({record['stdout']})・[stderr]({record['stderr']})・[コマンド/終了コード]({record['metadata']}) |")
rows.append('\n追加8 fixtureの対象比較も3版すべて一致（計24比較、全て終了0）。既存Laravel/API PlatformはGo3設定で96連続/8並列、タスクAPIは公式/Go3組で各191要求、Laravel12は8.5の公式/Go各111要求を確認した。')
report = LOGS/'REPORT.md'
text = report.read_text()
text = re.sub(r'<!-- FINAL_RESULTS -->.*?<!-- /FINAL_RESULTS -->', '\n'.join(rows), text, flags=re.S)
report.write_text(text)
print('PASS source, binary, 3 x 1201 fixtures, 24 targeted comparisons, all required regressions, protection and cleanup')
