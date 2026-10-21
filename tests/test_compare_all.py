"""Check version routing and failure propagation without requiring Podman.

The stand-in only records comparator invocations; it is not a PHP reference.
"""

import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest


ROOT = Path(__file__).resolve().parents[1]


class CompareAllTests(unittest.TestCase):
    def test_version_routing_and_exit_status(self):
        for codes, expected_exit in (([0, 0, 0], 0), ([1, 0, 0], 1), ([1, 2, 0], 2)):
            with self.subTest(codes=codes), tempfile.TemporaryDirectory() as directory:
                directory = Path(directory).resolve()
                checkout = directory / "relocated checkout"
                script = checkout / "containers/php-version/compare-all.sh"
                script.parent.mkdir(parents=True)
                shutil.copy2(ROOT / "containers/php-version/compare-all.sh", script)
                calls = directory / "calls.jsonl"
                stub = directory / "python3"
                stub.write_text(
                    "#!" + sys.executable + "\n"
                    "import json, os, sys\n"
                    "with open(os.environ['COMPARATOR_CALLS'], 'a') as log:\n"
                    "    log.write(json.dumps({'args': sys.argv[1:], 'cwd': os.getcwd()}) + '\\n')\n"
                    "version = sys.argv[sys.argv.index('--php-version') + 1]\n"
                    "codes = json.loads(os.environ['COMPARATOR_CODES'])\n"
                    "sys.exit(codes[version])\n"
                )
                stub.chmod(0o755)
                versions = ("8.3.35", "8.4.26", "8.5.11")
                env = dict(os.environ,
                           PATH=str(directory) + os.pathsep + os.environ.get("PATH", ""),
                           COMPARATOR_CALLS=str(calls),
                           COMPARATOR_CODES=json.dumps(dict(zip(versions, codes))))
                result = subprocess.run(
                    ["/bin/sh", str(script),
                     "tests/fixtures", "./bin/gophp", "--extra=preserved argument"],
                    cwd=directory, env=env, capture_output=True, timeout=10,
                )
                self.assertEqual(result.returncode, expected_exit, result.stderr)
                invocations = [json.loads(line) for line in calls.read_text().splitlines()]
                self.assertEqual(len(invocations), 3)
                for version, invocation in zip(versions, invocations):
                    self.assertEqual(Path(invocation["cwd"]), checkout.resolve())
                    args = invocation["args"]
                    target = version.rsplit(".", 1)[0]
                    self.assertEqual(args[0], "skills/php-go-interpreter/scripts/compare_php.py")
                    self.assertEqual(args[args.index("--php-version") + 1], version)
                    wrapper = checkout.resolve() / "containers/php-version/bin" / ("php" + target.replace(".", ""))
                    self.assertEqual(args[args.index("--php") + 1], str(wrapper))
                    self.assertEqual(args[args.index("--fixtures") + 1], "tests/fixtures")
                    self.assertEqual(args[args.index("--interpreter") + 1:],
                                     ["./bin/gophp", "--extra=preserved argument",
                                      "--php-version=" + target, "{file}"])


if __name__ == "__main__":
    unittest.main()
