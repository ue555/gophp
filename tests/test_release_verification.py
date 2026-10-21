"""Evidence and portability checks that do not need a container daemon."""
import importlib.util
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0,str(ROOT/'scripts'))
from release_evidence import source_identity
spec = importlib.util.spec_from_file_location('release_runner',ROOT/'scripts/verify-release.py')
release = importlib.util.module_from_spec(spec)
spec.loader.exec_module(release)


class ReleaseEvidenceTests(unittest.TestCase):
    def test_failure_keeps_real_exit_and_never_overwrites(self):
        with tempfile.TemporaryDirectory() as directory:
            runner = release.Runner(Path(directory))
            with self.assertRaises(RuntimeError):
                runner.run('failed',[sys.executable,'-c','print("original failure"); raise SystemExit(7)'])
            runner.run('passed',[sys.executable,'-c','print("later success")'])
            commands=json.loads((Path(directory)/'commands.json').read_text())
            self.assertEqual([c['exit_code'] for c in commands],[7,0])
            self.assertIn('original failure\nexit_code: 7', (Path(directory)/commands[0]['log']).read_text())

    def test_identity_includes_untracked_source_and_executable_mode(self):
        with tempfile.TemporaryDirectory(prefix='gophp release space ') as directory:
            root=Path(directory)
            subprocess.run(['git','init','-q',root],check=True)
            subprocess.run(['git','-C',root,'-c','user.name=Test','-c','user.email=test@example.invalid',
                            'commit','--allow-empty','-qm','base'],check=True)
            for app in ('laravel-api','task-api'):
                (root/'research'/app).mkdir(parents=True)
                (root/'research'/app/'composer.lock').write_text('{}')
            script=root/'new-script.py';script.write_text('print(1)\n')
            first=source_identity(root)
            self.assertIn('new-script.py',[f['path'] for f in first['manifest']])
            script.chmod(0o755)
            self.assertNotEqual(first['source_sha256'],source_identity(root)['source_sha256'])
            script.chmod(0o644)
            link=root/'dangling-link'; link.symlink_to('missing-target')
            linked=source_identity(root)
            self.assertIn('dangling-link',[f['path'] for f in linked['manifest']])
            link.unlink()
            (root/'docs/validation').mkdir(parents=True)
            (root/'docs/validation/log.txt').write_text('changing run evidence')
            self.assertEqual(first['source_sha256'],source_identity(root)['source_sha256'])

    def test_fixture_inventory_matches_all_expectations(self):
        self.assertGreater(release.fixture_inventory(),0)


if __name__ == '__main__':
    unittest.main()
