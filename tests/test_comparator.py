"""The fake processes below test the comparison harness, not PHP compatibility."""
from pathlib import Path
from contextlib import redirect_stderr, redirect_stdout
import importlib.util
import io
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[1]
COMPARATOR = ROOT / 'skills/php-go-interpreter/scripts/compare_php.py'
SPEC = importlib.util.spec_from_file_location('compare_php', COMPARATOR)
COMPARISON = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(COMPARISON)
PODMAN_ERROR = (b'Cannot connect to Podman. Please verify your connection.\n'
                b'Error: unable to connect to Podman socket: failed to connect: '
                b'ssh: handshake failed: EOF\n')
PODMAN_STOPPED_ERROR = (b'Error: can only create exec sessions on running containers: '
                        b'container state improper\n')


class ComparatorTests(unittest.TestCase):
    def compare_results(self, results, fixture_count=3):
        """Inject process outcomes to test outages without real timeouts or PHP."""
        with tempfile.TemporaryDirectory() as path:
            directory = Path(path) / 'runtime'
            directory.mkdir()
            for number in range(fixture_count):
                (directory / f'{number}.php').write_text('<?php echo 1;\n')
            argv = [str(COMPARATOR), '--php', 'php-ref', '--php-version', '8.4.26',
                    '--fixtures', str(directory), '--interpreter', 'gophp', '{file}']
            stdout, stderr = io.StringIO(), io.StringIO()
            with patch.object(sys, 'argv', argv), patch.object(COMPARISON, 'run', side_effect=results) as runner:
                with redirect_stdout(stdout), redirect_stderr(stderr):
                    code = COMPARISON.main()
            return code, stdout.getvalue(), stderr.getvalue(), runner.call_args_list

    def compare(self, directory, group, php_out, go_out, code=0):
        fixture = directory / group / 'case.php'
        fixture.parent.mkdir()
        fixture.write_text('<?php echo 1;\n')
        executables = []
        for label, out in [('php', php_out), ('go', go_out)]:
            path = directory / label
            path.write_text('#!' + sys.executable + '\n'
                            'import sys\n'
                            'if "-r" in sys.argv:\n'
                            '    sys.stdout.write("8.4.26"); sys.exit(0)\n'
                            'sys.stdout.buffer.write(' + repr(out) + ')\n'
                            'sys.exit(' + str(code) + ')\n')
            path.chmod(0o755)
            executables.append(path)
        return subprocess.run(
            [sys.executable, str(COMPARATOR), '--php', str(executables[0]),
             '--php-version', '8.4.26', '--fixtures', str(fixture),
             '--interpreter', str(executables[1]), '{file}'],
            capture_output=True, timeout=10, cwd=ROOT,
        )

    def test_equal_failures_are_not_success_fixtures(self):
        with tempfile.TemporaryDirectory() as path:
            result = self.compare(Path(path), 'success', b'', b'', 255)
            self.assertEqual(result.returncode, 2, result.stderr)
            self.assertIn(b'0 matched, 0 different, 1 errors; 1 fixtures', result.stdout)
            self.assertNotIn(b'MATCH ', result.stdout)

    def test_parser_failures_can_match(self):
        with tempfile.TemporaryDirectory() as path:
            result = self.compare(Path(path), 'parser', b'', b'', 255)
            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertIn(b'1 matched, 0 different, 0 errors; 1 fixtures', result.stdout)

    def test_whitespace_is_not_trimmed(self):
        with tempfile.TemporaryDirectory() as path:
            result = self.compare(Path(path), 'success', b'1\n', b'1')
            self.assertEqual(result.returncode, 1, result.stderr)
            self.assertIn(b'0 matched, 1 different, 0 errors; 1 fixtures', result.stdout)

    def test_podman_outage_is_error_and_remaining_fixtures_are_not_run(self):
        code, stdout, stderr, calls = self.compare_results([
            (0, b'8.4.26', b''),
            (0, b'1\n', b''), (0, b'1\n', b''),
            (125, b'', PODMAN_ERROR),
        ])
        self.assertEqual(code, 2, stderr)
        self.assertIn('1 matched, 0 different, 1 errors; 3 fixtures; 1 not run', stdout)
        self.assertIn('Podman connection failed (exit 125)', stdout)
        self.assertIn(repr(PODMAN_ERROR), stdout)
        self.assertNotIn('DIFF ', stdout)
        self.assertEqual(len(calls), 4)

    def test_stopped_podman_container_is_error_and_stops_version(self):
        code, stdout, stderr, calls = self.compare_results([
            (0, b'8.4.26', b''),
            (0, b'1\n', b''), (0, b'1\n', b''),
            (125, b'', PODMAN_STOPPED_ERROR),
            # These outcomes must never be consumed after the container stops.
            (0, b'1\n', b''),
            (0, b'1\n', b''), (0, b'1\n', b''),
            (0, b'1\n', b''), (0, b'1\n', b''),
        ], fixture_count=4)
        self.assertEqual(code, 2, stderr)
        self.assertIn('1 matched, 0 different, 1 errors; 4 fixtures; 2 not run', stdout)
        self.assertIn('Podman container is not running (exit 125)', stdout)
        self.assertIn(repr(PODMAN_STOPPED_ERROR), stdout)
        self.assertNotIn('DIFF ', stdout)
        self.assertEqual([call.args[0][0] for call in calls],
                         ['php-ref', 'php-ref', 'gophp', 'php-ref'])
        self.assertEqual(Path(calls[-1].args[0][-1]).name, '1.php')

    def test_similar_stopped_container_results_do_not_stop_comparison(self):
        for php, go in [
            ((125, b'', b'program error\n'), (0, b'', b'')),
            ((0, b'', PODMAN_STOPPED_ERROR), (0, b'', b'')),
            ((255, b'', PODMAN_STOPPED_ERROR), (0, b'', b'')),
            ((125, PODMAN_STOPPED_ERROR, b''), (0, b'', b'')),
            ((125, b'', b'Application: ' + PODMAN_STOPPED_ERROR), (0, b'', b'')),
            ((0, b'', b''), (125, b'', PODMAN_STOPPED_ERROR)),
        ]:
            with self.subTest(php=php, go=go):
                code, stdout, stderr, calls = self.compare_results([
                    (0, b'8.4.26', b''), php, go,
                    (0, b'1\n', b''), (0, b'1\n', b''),
                ], fixture_count=2)
                self.assertEqual(code, 1, stderr)
                self.assertIn('1 matched, 1 different, 0 errors; 2 fixtures', stdout)
                self.assertNotIn('not run', stdout)
                self.assertEqual(len(calls), 5)

    def test_php_timeout_with_failed_probe_stops_version(self):
        for probe in [(125, b'', PODMAN_ERROR),
                      subprocess.TimeoutExpired(['php-ref', '-r', 'echo PHP_VERSION;'], 5)]:
            with self.subTest(probe=probe):
                code, stdout, stderr, calls = self.compare_results([
                    (0, b'8.4.26', b''),
                    subprocess.TimeoutExpired(['php-ref', 'fixture.php'], 5),
                    probe,
                ])
                self.assertEqual(code, 2, stderr)
                self.assertIn('availability check failed', stdout)
                self.assertIn('0 matched, 0 different, 1 errors; 3 fixtures; 2 not run', stdout)
                self.assertEqual(len(calls), 3)
                self.assertEqual(calls[-1].args[0][-2:], ['-r', 'echo PHP_VERSION;'])

    def test_php_timeout_with_healthy_probe_continues(self):
        code, stdout, stderr, calls = self.compare_results([
            (0, b'8.4.26', b''),
            subprocess.TimeoutExpired(['php-ref', 'fixture.php'], 5),
            (0, b'8.4.26', b''),
            (0, b'1\n', b''), (0, b'1\n', b''),
            (0, b'1\n', b''), (0, b'1\n', b''),
        ])
        self.assertEqual(code, 2, stderr)
        self.assertIn('2 matched, 0 different, 1 errors; 3 fixtures', stdout)
        self.assertNotIn('not run', stdout)
        self.assertEqual(len(calls), 7)

    def test_similar_program_results_still_compare_exactly(self):
        for php, go in [
            ((125, b'', b'program error\n'), (0, b'', b'')),
            ((0, b'', PODMAN_ERROR), (0, b'', b'')),
            ((255, b'', PODMAN_ERROR), (0, b'', b'')),
            ((0, b'', b''), (125, b'', PODMAN_ERROR)),
        ]:
            with self.subTest(php=php, go=go):
                code, stdout, stderr, calls = self.compare_results(
                    [(0, b'8.4.26', b''), php, go], fixture_count=1)
                self.assertEqual(code, 1, stderr)
                self.assertIn('0 matched, 1 different, 0 errors; 1 fixtures', stdout)
                self.assertEqual(len(calls), 3)

    def test_program_exit_125_can_match(self):
        code, stdout, stderr, _ = self.compare_results([
            (0, b'8.4.26', b''),
            (125, b'', b''), (125, b'', b''),
        ], fixture_count=1)
        self.assertEqual(code, 0, stderr)
        self.assertIn('1 matched, 0 different, 0 errors; 1 fixtures', stdout)


if __name__ == '__main__':
    unittest.main()
