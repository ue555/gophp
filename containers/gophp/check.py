#!/usr/bin/env python3
"""Validate the container CLI against saved expectations, not official PHP."""

import argparse
import json
from pathlib import Path
import subprocess
import sys


ROOT = Path(__file__).resolve().parents[2]
RUNNER = ROOT / "containers/gophp/run.sh"
FIXTURES = ROOT / "tests/fixtures"
SMOKE = (
    "success/hello.php",
    "success/v1_report.php",
    "runtime/undefined.php",
    "runtime/division_zero.php",
    "parser/missing_expression.php",
    "runtime/logical_nan_not.php",
    "runtime/object_new_dereference.php",
)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--all", action="store_true", help="check every saved fixture")
    parser.add_argument("--timeout", type=float, default=15, help="seconds per execution")
    args = parser.parse_args()
    if not 0 < args.timeout < float("inf"):
        parser.error("--timeout must be finite and positive")
    expected = json.loads((FIXTURES / "expected.json").read_text())
    names = sorted(expected) if args.all else SMOKE
    count = 0

    def execute(arguments):
        result = subprocess.run(
            [str(RUNNER)] + arguments, cwd=ROOT, stdin=subprocess.DEVNULL,
            capture_output=True, timeout=args.timeout,
        )
        if result.returncode in (125, 126, 127) or result.returncode < 0:
            raise RuntimeError("container execution failed (exit {}): {}".format(
                result.returncode, result.stderr.decode(errors="replace")))
        return result.returncode, result.stdout, result.stderr

    try:
        for version in ("8.3", "8.4", "8.5"):
            for name in names:
                path = str(FIXTURES / name)
                want = expected[name]
                code = want.get("exit_by_version", {}).get(version, want["exit"])
                stdout = want.get("stdout_by_version", {}).get(version, want["stdout"])
                stderr = want.get("stderr_by_version", {}).get(version, want["stderr"])
                wanted = (code, stdout.encode(), stderr.replace("{file}", path).encode())
                actual = execute(["--php-version=" + version, path])
                if actual != wanted:
                    print("FAIL: PHP target {} {}\nexpected: {!r}\nactual: {!r}".format(
                        version, name, wanted, actual), file=sys.stderr)
                    return 1
                count += 1
            print("PASS: PHP target {}: {} fixtures".format(version, len(names)), flush=True)
        usage = b"usage: gophp [--php-version=8.3|8.4|8.5] <file.php>\n"
        for arguments, wanted in ((["--help"], (0, b"", usage)),
                                  ([], (2, b"", usage))):
            if execute(arguments) != wanted:
                print("FAIL: CLI arguments {!r}".format(arguments), file=sys.stderr)
                return 1
        print("PASS: {} fixture executions and 2 CLI checks (saved expectations)".format(count))
        return 0
    except (OSError, subprocess.TimeoutExpired, RuntimeError) as error:
        print("ERROR: {}".format(error), file=sys.stderr)
        return 2


if __name__ == "__main__":
    sys.exit(main())
