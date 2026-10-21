#!/usr/bin/env python3
"""Prepare a persistent manual-use DB and caches with Go PHP; never reset the DB."""
import argparse
import os
from pathlib import Path
from verify import APP, ROOT, VERSIONS, command


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--version', choices=VERSIONS, default='8.4')
    parser.add_argument('--interpreter', type=Path, default=ROOT/'bin/gophp')
    args = parser.parse_args()
    env = os.environ.copy()
    # Match Laravel's application-relative default paths; explicit overrides are
    # resolved here and passed through to every preparation command.
    database = Path(env.get('DB_DATABASE', APP/'database/database.sqlite')).resolve()
    storage = Path(env.get('LARAVEL_STORAGE_PATH', APP/'storage')).resolve()
    database.parent.mkdir(parents=True, exist_ok=True)
    database.touch(exist_ok=True)
    for directory in ('framework/cache/data','framework/views','framework/sessions','logs'):
        (storage/directory).mkdir(parents=True, exist_ok=True)
    (APP/'bootstrap/cache').mkdir(parents=True, exist_ok=True)
    env.update(DB_DATABASE=str(database), LARAVEL_STORAGE_PATH=str(storage))
    prefix = [str(args.interpreter.resolve()), '--php-version='+args.version, 'artisan']
    command(prefix+['package:discover','--no-interaction'], env)
    command(prefix+['migrate','--force','--no-interaction'], env)
    command(prefix+['migrate:status','--no-interaction'], env)
    print('Prepared DB:', database)
    print('Storage:', storage)
    print('Existing rows retained; no reset or deletion performed.')


if __name__ == '__main__':
    main()
