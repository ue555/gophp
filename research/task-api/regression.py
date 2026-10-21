#!/usr/bin/env python3
"""Run the existing fixed app's Go TCP checks with isolated bootstrap caches."""
import argparse
import importlib.util
import os
from pathlib import Path
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[2]
LEGACY = ROOT/'research/laravel-api'
sys.path.insert(0, str(LEGACY))
spec = importlib.util.spec_from_file_location('legacy_go_verifier', LEGACY/'verify-gophp.py')
legacy = importlib.util.module_from_spec(spec)
spec.loader.exec_module(legacy)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--version', choices=('8.3','8.4','8.5'), action='append')
    parser.add_argument('--stress-requests', type=int, default=100)
    args = parser.parse_args()
    if not 100 <= args.stress_requests <= 2000:
        parser.error('--stress-requests must be between 100 and 2000')
    (ROOT/'.cache').mkdir(exist_ok=True)
    for version in args.version or ('8.3','8.4','8.5'):
        with tempfile.TemporaryDirectory(prefix='task-api-legacy-bootstrap-', dir=ROOT/'.cache') as temp:
            for key in ('PACKAGES','SERVICES','CONFIG','ROUTES','EVENTS'):
                os.environ['APP_'+key+'_CACHE'] = str(Path(temp)/(key.lower()+'.php'))
            legacy.verify_version(version, ROOT/'bin/gophp', stress_requests=args.stress_requests)


if __name__ == '__main__':
    main()
