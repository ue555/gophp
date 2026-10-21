#!/usr/bin/env python3
"""Record official observations for new fixtures; never overwrite existing expectations."""
import json,os,shlex,subprocess,sys,tempfile
from pathlib import Path
ROOT=Path(__file__).resolve().parents[2]
fixtures=ROOT/'tests/fixtures'
expected_path=fixtures/'expected.json'
expected=json.loads(expected_path.read_text())
baseline=json.loads(subprocess.check_output(["git","show","HEAD:tests/fixtures/expected.json"],cwd=ROOT))
paths=[Path(p).resolve() for p in sys.argv[1:]]
if not paths: paths=sorted(p for p in fixtures.rglob('*.php') if str(p.relative_to(fixtures)) not in baseline)
out=Path(tempfile.mkdtemp(prefix='fresh-laravel-official-observations-',dir=ROOT/'docs/validation'))
print('evidence:',out,flush=True)
observed={};commands=[]
for v,patch in [('8.4','8.4.26'),('8.3','8.3.35'),('8.5','8.5.11')]:
 php=ROOT/'containers/php-version/bin'/('php'+v.replace('.',''))
 env=dict(os.environ,LC_ALL='C',LANG='C',TZ='UTC')
 assert subprocess.check_output([str(php),'-n','-r','echo PHP_VERSION;'],env=env,text=True)==patch
 for p in paths:
  key=str(p.relative_to(fixtures));name=v+'-'+p.stem
  cmd=[str(php),'-n','-d','display_errors=stderr','-d','log_errors=0','-d','html_errors=0','-d','error_reporting=-1','-d','date.timezone=UTC',str(p)]
  r=subprocess.run(cmd,capture_output=True,env=env,timeout=30)
  (out/(name+'.stdout')).write_bytes(r.stdout);(out/(name+'.stderr')).write_bytes(r.stderr)
  commands.append({'command':cmd,'version':patch,'fixture':key,'exit_code':r.returncode})
  value={'stdout':r.stdout.decode().replace(str(p),'{file}'),'stderr':r.stderr.decode().replace(str(p),'{file}'),'exit':r.returncode}
  if v=='8.4': observed[key]=value
  else:
   for k,x in value.items():
    if x!=observed[key][k]: observed[key].setdefault(k+'_by_version',{})[v]=x
  print(patch,key,'exit',r.returncode,flush=True)
(out/'commands.json').write_text(json.dumps(commands,indent=2)+'\n')
for key,value in observed.items():
 if key in baseline: assert expected[key]==value, 'existing expectation differs: '+key
 else: expected[key]=value
expected_path.write_text(json.dumps(expected,ensure_ascii=False,indent=2)+'\n')
(out/'observed.json').write_text(json.dumps(observed,ensure_ascii=False,indent=2)+'\n')
