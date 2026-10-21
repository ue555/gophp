#!/usr/bin/env python3
"""Run a command and preserve its output and actual exit code in a unique log."""
import argparse
from datetime import datetime
from pathlib import Path
import shlex
import subprocess
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[2]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--label', required=True, help='target/runtime/check name')
    parser.add_argument('command', nargs=argparse.REMAINDER)
    args = parser.parse_args()
    cmd = args.command[1:] if args.command[:1] == ['--'] else args.command
    if not cmd:
        parser.error('command is required after --')
    label = ''.join(ch if ch.isalnum() or ch in '.-' else '-' for ch in args.label)
    with tempfile.NamedTemporaryFile(mode='w', prefix=f'task-api-{label}-{datetime.now():%Y%m%d-%H%M%S}-', suffix='.log', dir=ROOT/'docs/validation', delete=False) as log:
        print('log:', log.name, flush=True)
        def write(text):
            log.write(text); log.flush(); print(text, end='', flush=True)
        write('target: '+args.label+'\ncommand: '+shlex.join(cmd)+'\ncwd: '+str(Path.cwd())+'\n')
        try:
            process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
            for line in process.stdout:
                write(line)
            code = process.wait()
        except OSError as error:
            write(repr(error)+'\n'); code = 127
        write(f'exit_code: {code}\n')
        return code if code >= 0 else 128-code


if __name__ == '__main__':
    raise SystemExit(main())
