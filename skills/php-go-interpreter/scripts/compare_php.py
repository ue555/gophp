#!/usr/bin/env python3
"""Compare PHP fixtures byte for byte against a pinned official PHP binary."""

import argparse
import math
import os
from pathlib import Path
import signal
import subprocess
import sys


def run(command, timeout, env):
    # On POSIX, terminate child processes as well when a runner times out.
    with subprocess.Popen(
        command, stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
        stderr=subprocess.PIPE, env=env, start_new_session=(os.name == "posix"),
    ) as process:
        try:
            stdout, stderr = process.communicate(timeout=timeout)
        except subprocess.TimeoutExpired:
            if os.name == "posix":
                try:
                    os.killpg(process.pid, signal.SIGKILL)
                except ProcessLookupError:
                    pass
            else:
                process.kill()
            process.communicate()
            raise
        if process.returncode < 0:
            raise RuntimeError("process terminated by signal {}".format(-process.returncode))
        return process.returncode, stdout, stderr


def preview(data):
    suffix = " ... ({} bytes total)".format(len(data)) if len(data) > 1000 else ""
    return repr(data[:1000]) + suffix


def podman_environment_error(result):
    # Exit 125 alone can also be a program's exit status. Require Podman's
    # connection or stopped-container diagnostic on stderr as well.
    code, _, stderr = result
    if code != 125:
        return None
    if (stderr.startswith(b"Cannot connect to Podman.")
            and b"Error: unable to connect to Podman socket:" in stderr):
        return "Podman connection failed"
    if stderr.startswith(b"Error: can only create exec sessions on running containers"):
        return "Podman container is not running"
    return None


def verify_php(command, expected_version, timeout, env):
    code, stdout, stderr = run(command + ["-r", "echo PHP_VERSION;"], timeout, env)
    actual_version = stdout.decode("ascii").strip()
    if code != 0 or stderr or actual_version != expected_version:
        raise RuntimeError("expected PHP {}; got {!r}, exit {}, stderr {}".format(
            expected_version, actual_version, code, preview(stderr)))
    return actual_version


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--php", default="php", help="official PHP CLI executable")
    parser.add_argument("--php-version", required=True, help="exact PHP_VERSION, including patch")
    parser.add_argument("--fixtures", required=True, type=Path, help="one .php file or directory")
    parser.add_argument("--timeout", type=float, default=5.0, help="seconds per process (default: 5)")
    parser.add_argument("--interpreter", required=True, nargs=argparse.REMAINDER,
                        help="runner argv with one standalone {file}; put this option last")
    args = parser.parse_args()
    if not math.isfinite(args.timeout) or args.timeout <= 0:
        parser.error("--timeout must be a finite positive number")
    if not args.interpreter or args.interpreter.count("{file}") != 1:
        parser.error("--interpreter must contain exactly one standalone {file} argument")

    fixture_path = args.fixtures.resolve()
    if fixture_path.is_file() and fixture_path.suffix == ".php":
        fixtures = [fixture_path]
    elif fixture_path.is_dir():
        fixtures = sorted(path for path in fixture_path.rglob("*.php") if path.is_file())
    else:
        print("ERROR: fixtures must be a .php file or directory", file=sys.stderr)
        return 2
    if not fixtures:
        print("ERROR: no PHP fixtures found", file=sys.stderr)
        return 2

    env = dict(os.environ, LC_ALL="C", LANG="C", TZ="UTC")
    php_command = [args.php, "-n", "-d", "display_errors=stderr", "-d", "log_errors=0",
                   "-d", "html_errors=0", "-d", "error_reporting=-1", "-d", "date.timezone=UTC"]
    try:
        actual_version = verify_php(php_command, args.php_version, args.timeout, env)
    except (OSError, RuntimeError, subprocess.TimeoutExpired, UnicodeError) as error:
        print("ERROR: could not verify official PHP: {}".format(error), file=sys.stderr)
        return 2

    print("PHP: {} ({})".format(actual_version, args.php))
    print("Working directory: {}".format(Path.cwd()))
    passed = failed = errors = 0
    not_run = 0
    for fixture_index, fixture in enumerate(fixtures):
        commands = {
            "PHP": php_command + [str(fixture)],
            "Go": [str(fixture) if arg == "{file}" else arg for arg in args.interpreter],
        }
        results = {}
        reference_unavailable = False
        for label, command in commands.items():
            try:
                result = run(command, args.timeout, env)
            except (OSError, RuntimeError, subprocess.TimeoutExpired) as error:
                print("ERROR {} [{}]: {}".format(fixture, label, error))
                errors += 1
                if label == "PHP":
                    # A single fixture may hang while PHP remains healthy. Only
                    # stop the version run if a fresh version probe also fails.
                    try:
                        verify_php(php_command, args.php_version, args.timeout, env)
                    except (OSError, RuntimeError, subprocess.TimeoutExpired, UnicodeError) as probe_error:
                        print("ERROR: official PHP availability check failed: {}".format(probe_error))
                        reference_unavailable = True
                break
            environment_error = podman_environment_error(result) if label == "PHP" else None
            if environment_error:
                print("ERROR {} [PHP]: {} (exit {})".format(fixture, environment_error, result[0]))
                print("  PHP stdout: {}; stderr: {}".format(preview(result[1]), preview(result[2])))
                errors += 1
                reference_unavailable = True
                break
            results[label] = result
        if reference_unavailable:
            not_run = len(fixtures) - fixture_index - 1
            print("ERROR: stopping this version because official PHP is unavailable; {} fixtures not run".format(not_run))
            break
        if len(results) != 2:
            continue
        if fixture.parent.name == "success" and any(result[0] != 0 for result in results.values()):
            print("ERROR {}: success fixture must exit 0 in both interpreters (PHP {}, Go {})".format(
                fixture, results["PHP"][0], results["Go"][0]))
            for label, result in results.items():
                print("  {} stdout: {}; stderr: {}".format(label, preview(result[1]), preview(result[2])))
            errors += 1
            continue
        if results["PHP"] == results["Go"]:
            passed += 1
            print("MATCH {} (exit {})".format(fixture, results["PHP"][0]))
            continue
        failed += 1
        print("DIFF {}".format(fixture))
        for index, field in enumerate(("exit code", "stdout", "stderr")):
            reference, actual = results["PHP"][index], results["Go"][index]
            if reference != actual:
                print("  {}:".format(field))
                print("    PHP: {}".format(preview(reference) if isinstance(reference, bytes) else reference))
                print("    Go:  {}".format(preview(actual) if isinstance(actual, bytes) else actual))
    summary = "{} matched, {} different, {} errors; {} fixtures".format(passed, failed, errors, len(fixtures))
    if not_run:
        summary += "; {} not run".format(not_run)
    print(summary)
    return 2 if errors else (1 if failed else 0)


if __name__ == "__main__":
    sys.exit(main())
