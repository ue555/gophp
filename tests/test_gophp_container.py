"""Check container launcher argument/stream/exit preservation without Podman."""

import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest


ROOT = Path(__file__).resolve().parents[1]


class ContainerLauncherTests(unittest.TestCase):
    def test_run_preserves_arguments_streams_and_exit(self):
        with tempfile.TemporaryDirectory(prefix="gophp container ") as directory:
            directory = Path(directory)
            calls = directory / "calls.json"
            stub = directory / "podman"
            stub.write_text(
                "#!" + sys.executable + "\n"
                "import json, os, sys\n"
                "with open(os.environ['GOPHP_TEST_CALLS'], 'w') as log:\n"
                "    json.dump(sys.argv[1:], log)\n"
                "sys.stdout.buffer.write(b'output\\x00\\n')\n"
                "sys.stderr.buffer.write(b'diagnostic\\n')\n"
                "sys.exit(int(os.environ['GOPHP_TEST_EXIT']))\n"
            )
            stub.chmod(0o755)
            for code in (0, 2, 125, 255):
                with self.subTest(code=code):
                    env = dict(os.environ, PATH=str(directory) + os.pathsep + os.environ["PATH"],
                               GOPHP_TEST_CALLS=str(calls), GOPHP_TEST_EXIT=str(code),
                               GOPHP_IMAGE="localhost/gophp:test")
                    arguments = ["--php-version=8.5", "tests/空白 file.php"]
                    result = subprocess.run(
                        [str(ROOT / "containers/gophp/run.sh")] + arguments,
                        cwd=directory, env=env, capture_output=True, timeout=10,
                    )
                    self.assertEqual((result.returncode, result.stdout, result.stderr),
                                     (code, b"output\x00\n", b"diagnostic\n"))
                    args = json.loads(calls.read_text())
                    self.assertEqual(args[0], "run")
                    self.assertEqual(args[-3:], ["localhost/gophp:test"] + arguments)
                    self.assertEqual(args[args.index("--volume") + 1], "{0}:{0}:ro".format(ROOT))
                    self.assertEqual(args[args.index("--workdir") + 1], str(ROOT))
                    self.assertIn("--rm", args)
                    self.assertNotIn("-t", args)


if __name__ == "__main__":
    unittest.main()
