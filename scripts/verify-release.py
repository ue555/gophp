#!/usr/bin/env python3
"""Run release gates on a fresh checkout, preserving every command and artifact."""
import argparse
from datetime import datetime
import json
import os
from pathlib import Path
import re
import shlex
import subprocess
import sys
import tempfile
import traceback
import uuid

from release_evidence import digest, source_identity

ROOT = Path(__file__).resolve().parents[1]
VERSIONS = ('8.3.35','8.4.26','8.5.11')


class Runner:
    def __init__(self, directory):
        self.directory = directory
        self.commands = []

    def run(self, label, args, env=None):
        args = list(map(str,args))
        path = self.directory/f'{len(self.commands):02d}-{label}.log'
        entry = {'label':label,'command':args,'log':path.name,'target':label}
        self.commands.append(entry)
        print('RUN',label,flush=True)
        with path.open('x', buffering=1) as log:
            log.write('target: '+label+'\ncommand: '+shlex.join(args)+'\ncwd: '+str(ROOT)+'\n')
            try:
                result = subprocess.run(args,cwd=ROOT,env=env,stdout=log,stderr=subprocess.STDOUT)
                code = result.returncode
            except OSError as error:
                log.write(repr(error)+'\n'); code = 127
            log.write(f'exit_code: {code}\n')
        entry['exit_code'] = code
        (self.directory/'commands.json').write_text(json.dumps(self.commands,indent=2)+'\n')
        if code:
            raise RuntimeError(f'{label} exited {code}; see {path}')
        return path


def fixture_inventory():
    expected = json.loads((ROOT/'tests/fixtures/expected.json').read_text())
    actual = {str(path.relative_to(ROOT/'tests/fixtures')) for path in (ROOT/'tests/fixtures').rglob('*.php')}
    if actual != set(expected):
        raise RuntimeError(f'fixture/expectation mismatch: {actual ^ set(expected)}')
    return len(actual)


def verify(args, runner):
    count = fixture_inventory()
    initial = source_identity(ROOT)
    (runner.directory/'source-before.json').write_text(json.dumps(initial,indent=2)+'\n')
    (runner.directory/'fixture-count.json').write_text(json.dumps({'fixtures':count,'expectations':count})+'\n')
    if args.require_fresh_vendor:
        for app in ('laravel-api','task-api'):
            if (ROOT/'research'/app/'vendor').exists() or (ROOT/'research'/app/'bootstrap/cache').exists():
                raise RuntimeError(f'fresh verification requires absent vendor and bootstrap/cache: {app}')
    runner.run('container-environment',[args.engine,'info'])
    versions = args.version or list(VERSIONS)
    for version in versions:
        image = args.image_template.format(version=version)
        if args.build_images:
            runner.run('build-dependencies-'+version,[args.engine,'build','--build-arg','PHP_VERSION='+version,
                        '-f','containers/dependency-research/Containerfile','-t',image,'.'])
        user_args = ['--user',f'{os.getuid()}:{os.getgid()}'] if args.engine == 'docker' else []
        for app in ('laravel-api','task-api'):
            app_path = ROOT/'research'/app
            mounts = ['-v',f'{app_path}:/app',*user_args,'-e','COMPOSER_HOME=/tmp/composer']
            runner.run(f'install-{app}-{version}',[args.engine,'run','--rm',*mounts,image,'composer','install',
                        '--no-dev','--no-interaction','--prefer-dist','--no-scripts','--no-plugins'])
            runner.run(f'platform-{app}-{version}',[args.engine,'run','--rm',*mounts,image,'composer','check-platform-reqs','--no-dev'])
            installed = json.loads((app_path/'vendor/composer/installed.json').read_text())['packages']
            locked = json.loads((app_path/'composer.lock').read_text())['packages']
            if sorted((p['name'],p['version']) for p in installed) != sorted((p['name'],p['version']) for p in locked):
                raise RuntimeError(f'vendor does not match lock: {app}')
            for directory in ('bootstrap/cache','storage/framework/cache/data','storage/framework/views','storage/framework/sessions','storage/logs'):
                (app_path/directory).mkdir(parents=True,exist_ok=True)
    # Vendor prerequisites are mandatory, so the existing integration tests
    # cannot silently skip their app checks on a fresh CI checkout.
    for app in ('laravel-api','task-api'):
        if not (ROOT/'research'/app/'vendor/autoload.php').is_file():
            raise RuntimeError('missing fixed vendor: '+app)
    runner.run('local-test-vet-build',['make','test','vet','build'])
    runner.run('diff-check',['git','diff','--check'])
    built = source_identity(ROOT)
    (runner.directory/'source-built.json').write_text(json.dumps(built,indent=2)+'\n')
    for version in versions:
        target = version.rsplit('.',1)[0]
        # Own oracle container, preserving absolute diagnostic paths. Never
        # reuse/stop the user's existing php-version containers.
        name = 'gophp-release-oracle-'+uuid.uuid4().hex[:12]
        started = False
        try:
            runner.run('oracle-start-'+version,[args.engine,'run','-d','--name',name,
                        '-v',f'{ROOT}:{ROOT}:ro','-w',ROOT,f'php:{version}-cli-bookworm',
                        'php','-r','while (true) { sleep(3600); }'])
            started = True
            wrapper = runner.directory/('php-'+version)
            wrapper.write_text('#!/bin/sh\nexec '+shlex.join([args.engine,'exec','-i','-w',str(ROOT),'-e','LC_ALL=C','-e','LANG=C','-e','TZ=UTC',name,'php'])+' "$@"\n')
            wrapper.chmod(0o755)
            output = runner.run('fixtures-'+version,[sys.executable,'skills/php-go-interpreter/scripts/compare_php.py',
                        '--php',wrapper,'--php-version',version,'--fixtures','tests/fixtures','--timeout','15',
                        '--interpreter','./bin/gophp','--php-version='+target,'{file}'])
            text = output.read_text()
            if text.count('\nMATCH ') != count or f'{count} matched, 0 different, 0 errors; {count} fixtures' not in text:
                raise RuntimeError('fixture comparison summary/individual count mismatch')
        finally:
            if started:
                runner.run('oracle-remove-'+version,[args.engine,'rm','-f',name])
        http_root = runner.directory/('http-'+version)
        runner.run('task-http-'+version,[sys.executable,'research/task-api/verify.py','--version',target,
                    '--engine',args.engine,'--image-template',args.image_template,'--output-root',http_root])
        artifacts = list(http_root.glob('task-api-*'))
        if len(artifacts) != 1:
            raise RuntimeError('missing/unexpected HTTP artifact directory')
        runner.run('task-http-audit-'+version,[sys.executable,'research/task-api/audit-evidence.py',artifacts[0],'--version',target])
        runner.run('legacy-http-'+version,[sys.executable,'research/task-api/regression.py','--version',target,
                    '--stress-requests','2000'])
    final = source_identity(ROOT)
    (runner.directory/'source-after.json').write_text(json.dumps(final,indent=2)+'\n')
    if initial['source_sha256'] != final['source_sha256'] or initial['locks'] != final['locks']:
        raise RuntimeError('source or locks changed during verification')
    if built['binary_sha256'] != final['binary_sha256']:
        raise RuntimeError('binary changed during verification')
    print(f'PASS {count} fixtures, task API 191 requests/runtime, legacy 2000 sequential/8 parallel, versions={versions}',flush=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--engine',choices=('podman','docker'),default='podman')
    parser.add_argument('--version',choices=VERSIONS,action='append')
    parser.add_argument('--image-template',default='localhost/gophp-dependency-research:{version}')
    parser.add_argument('--build-images',action='store_true')
    parser.add_argument('--require-fresh-vendor',action='store_true')
    parser.add_argument('--output-root',type=Path,default=ROOT/'docs/validation')
    args = parser.parse_args()
    args.output_root.mkdir(parents=True,exist_ok=True)
    directory = Path(tempfile.mkdtemp(prefix='release-task-api-'+datetime.now().strftime('%Y%m%d-%H%M%S')+'-',dir=args.output_root)).resolve()
    print('evidence:',directory,flush=True)
    runner = Runner(directory)
    code = 0
    try:
        verify(args,runner)
    except Exception:
        traceback.print_exc(); code=1
    (directory/'result.json').write_text(json.dumps({'command':sys.argv,'exit_code':code,'versions':args.version or list(VERSIONS)},indent=2)+'\n')
    print('exit_code:',code,flush=True)
    return code


if __name__ == '__main__':
    raise SystemExit(main())
