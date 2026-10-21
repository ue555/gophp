#!/usr/bin/env python3
"""Verify the untouched Laravel skeleton with dev dependencies over real TCP."""
import argparse,base64,hashlib,json,os,re,shlex,shutil,socket,sqlite3,subprocess,sys,tempfile,time,traceback,uuid
from http.cookies import SimpleCookie
from email.utils import parsedate_to_datetime
from pathlib import Path
from concurrent.futures import ThreadPoolExecutor
from urllib.request import Request,urlopen
from urllib.error import HTTPError,URLError
ROOT=Path(__file__).resolve().parents[2]
APP=Path(__file__).resolve().parent/'app'
VERSIONS={'8.3':'8.3.35','8.4':'8.4.26','8.5':'8.5.11'}
def check(ok,msg):
 if not ok:raise RuntimeError(msg)
def sha(path):return hashlib.sha256(Path(path).read_bytes()).hexdigest()
def db_state(db):
 with sqlite3.connect(db) as c:
  schema=c.execute("select name,sql from sqlite_master where type in ('table','index') and name not like 'sqlite_%' order by name").fetchall()
  migrations=c.execute('select migration,batch from migrations order by id').fetchall()
 return {'schema':schema,'migrations':migrations}
class Run:
 def __init__(self,out,args):self.out=out;self.args=args;self.commands=[]
 def command(self,label,cmd,cwd,env=None):
  cmd=list(map(str,cmd));print('RUN',label,flush=True)
  r=subprocess.run(cmd,cwd=cwd,env=env,capture_output=True,timeout=180)
  (self.out/(label+'.stdout')).write_bytes(r.stdout);(self.out/(label+'.stderr')).write_bytes(r.stderr)
  self.commands.append({'label':label,'command':cmd,'cwd':str(cwd),'exit_code':r.returncode})
  (self.out/'commands.json').write_text(json.dumps(self.commands,indent=2)+'\n')
  check(r.returncode==0,f'{label} exit {r.returncode}: {r.stderr[:600]!r}')
  check(not r.stderr,f'{label} stderr: {r.stderr[:600]!r}')
  return r.stdout.decode()
 def runtime(self,kind,version,temp):
  label=kind+'-'+version;work=Path(temp)/label;app=work/'app';work.mkdir()
  shutil.copytree(APP,app,ignore=shutil.ignore_patterns('vendor','.env','cache','logs','views','sessions','*.sqlite'))
  # The unchanged welcome view is source, unlike storage/framework/views.
  shutil.copytree(APP/'resources/views',app/'resources/views')
  (app/'vendor').symlink_to(APP/'vendor',target_is_directory=True)
  for d in ['bootstrap/cache','storage/framework/cache/data','storage/framework/views','storage/framework/sessions','storage/logs']:(app/d).mkdir(parents=True,exist_ok=True)
  db=work/'database.sqlite';db.touch()
  envtext=(APP/'.env.example').read_text().replace('APP_NAME=Laravel','APP_NAME="Fresh Laravel Probe"').replace('APP_KEY=','APP_KEY=base64:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=')
  envtext+='\nDB_DATABASE='+str(db)+'\n';(app/'.env').write_text(envtext)
  env=dict(os.environ,LC_ALL='C',LANG='C',TZ='UTC',COLUMNS='120',PATH='/usr/bin:/bin')
  for key in ['APP_NAME','APP_ENV','APP_KEY','APP_DEBUG','DB_DATABASE','APP_BASE_PATH','APP_CONFIG_CACHE','APP_PACKAGES_CACHE','APP_SERVICES_CACHE','APP_ROUTES_CACHE','APP_EVENTS_CACHE']:env.pop(key,None)
  user_args=['--user',f'{os.getuid()}:{os.getgid()}'] if Path(self.args.engine).name=='docker' else []
  prefix=[str(self.args.binary),'--php-version='+version] if kind=='go' else [self.args.engine,'run','--entrypoint=/usr/local/bin/docker-php-entrypoint','--rm',*user_args,'-v',f'{ROOT}:{ROOT}','-w',str(app),'-e','LC_ALL=C','-e','LANG=C','-e','TZ=UTC','-e','COLUMNS=120','-e','PATH=/usr/bin:/bin','localhost/gophp-dependency-research:'+VERSIONS[version],'/usr/local/bin/php']
  autoload=work/'autoload.php';autoload.write_text("<?php require __DIR__.'/app/vendor/autoload.php'; echo \"autoload OK\\n\";\n")
  check(self.command(label+'-autoload',prefix+[str(autoload)],app,env)=='autoload OK\n','autoload output')
  probe=work/'environment.php';probe.write_text("<?php require __DIR__.'/app/vendor/autoload.php'; $app=require __DIR__.'/app/bootstrap/app.php'; $app->make(Illuminate\\Contracts\\Console\\Kernel::class)->bootstrap(); echo json_encode([config('app.name'),config('database.connections.sqlite.database'),env('VITE_APP_NAME')]);\n")
  result=json.loads(self.command(label+'-environment',prefix+[str(probe)],app,env));check(result==['Fresh Laravel Probe',str(db),'Fresh Laravel Probe'],'dotenv values differ: '+repr(result))
  listing=self.command(label+'-list',prefix+['artisan','list','--no-ansi'],app,env)
  check('Laravel Framework 12.70.0' in listing and all(x in listing for x in ['migrate','pail','tinker','test','route:list']),'missing standard/dev commands')
  about=json.loads(self.command(label+'-about',prefix+['artisan','about','--json','--no-ansi'],app,env));check(about['environment']['application_name']=='Fresh Laravel Probe' and about['environment']['laravel_version']=='12.70.0','about')
  check(about['environment']['php_version']==VERSIONS[version],'wrong runtime patch version')
  routes=json.loads(self.command(label+'-routes',prefix+['artisan','route:list','--json','--no-ansi'],app,env));check(any(x['uri']=='/' and x['middleware']==['web'] for x in routes) and any(x['uri']=='up' for x in routes),'standard routes/middleware changed')
  migration=self.command(label+'-migrate',prefix+['artisan','migrate','--force','--no-ansi'],app,env)
  before=db_state(db);check(len(before['migrations'])==3,'standard migrations missing')
  check({'users','password_reset_tokens','sessions','cache','cache_locks','jobs','job_batches','failed_jobs','migrations'} <= {x[0] for x in before['schema']},'standard tables missing')
  again=self.command(label+'-migrate-again',prefix+['artisan','migrate','--force','--no-ansi'],app,env);check('Nothing to migrate' in again and db_state(db)==before,'migration rerun changed DB')
  (self.out/(label+'-db.json')).write_text(json.dumps(before,indent=2)+'\n')
  port_socket=socket.socket();port_socket.bind(('127.0.0.1',0));port=port_socket.getsockname()[1];port_socket.close()
  container='fresh-laravel-'+uuid.uuid4().hex[:12]
  if kind=='go':cmd=prefix+[f'--serve=127.0.0.1:{port}','public/index.php']
  else:cmd=[self.args.engine,'run','--entrypoint=/usr/local/bin/docker-php-entrypoint','--rm',*user_args,'--name',container,'-p',f'127.0.0.1:{port}:8080','-v',f'{ROOT}:{ROOT}','-w',str(app),'-e','TZ=UTC','-e','PATH=/usr/bin:/bin','localhost/gophp-dependency-research:'+VERSIONS[version],'/usr/local/bin/php','-S','0.0.0.0:8080','-t','public','public/index.php']
  log=(self.out/(label+'-server.log')).open('wb');process=subprocess.Popen(cmd,cwd=app,env=env,stdout=log,stderr=subprocess.STDOUT)
  entry={'label':label+'-server','command':cmd,'cwd':str(app),'pid':process.pid};self.commands.append(entry)
  responses=[]
  def request(path,cookie=''):
   try:r=urlopen(Request(f'http://127.0.0.1:{port}'+path,headers={'Accept':'text/html','Cookie':cookie}),timeout=60)
   except HTTPError as error:r=error
   with r:record={'path':path,'status':r.status,'headers':list(r.headers.items()),'body':r.read().decode()}
   return record
  def accept(record):
   responses.append(record)
   with (self.out/(label+'-http.jsonl')).open('a') as f:f.write(json.dumps(record,ensure_ascii=False)+'\n')
   check(record['status']==200,f'HTTP {record["path"]}: {record["status"]}; raw body saved')
   headers={k.lower():v for k,v in record['headers']};check(headers.get('content-type','').lower().startswith('text/html'),'HTML Content-Type')
   check(abs(parsedate_to_datetime(headers['date']).timestamp()-time.time())<120,'HTTP Date')
   cookies=SimpleCookie()
   for k,v in record['headers']:
    if k.lower()=='set-cookie':cookies.load(v)
   if record['path']=='/':
    check(set(cookies)=={'XSRF-TOKEN','fresh-laravel-probe-session'},'standard session cookies: '+repr(list(cookies)))
    for c in cookies.values():
     check(c['path']=='/' and c['samesite'].lower()=='lax' and 7190<=int(c['max-age'])<=7200,'cookie attributes')
     check(abs(parsedate_to_datetime(c['expires']).timestamp()-time.time()-7200)<120,'cookie expiry')
     check(abs(parsedate_to_datetime(c['expires']).timestamp()-parsedate_to_datetime(headers['date']).timestamp()-int(c['max-age']))<=2,'cookie expiry/max-age consistency')
    check(bool(cookies['fresh-laravel-probe-session']['httponly']) and not cookies['XSRF-TOKEN']['httponly'],'cookie HttpOnly')
   else:check(not cookies,'health request leaked cookies')
   check('Fatal error' not in record['body'] and 'Warning:' not in record['body'],'diagnostic in HTML')
   check('<!DOCTYPE html>' in record['body'] or '<!doctype html>' in record['body'],'missing HTML')
   if record['path']=='/':check('Laravel' in record['body'] and '<html lang="en">' in record['body'] and '@if' not in record['body'],'welcome Blade output')
  try:
   for _ in range(100):
    check(process.poll() is None,'server exited early')
    if kind=='official' and 'Development Server' not in (self.out/(label+'-server.log')).read_text():
     time.sleep(.1);continue
    try:
     with socket.create_connection(('127.0.0.1',port),timeout=.2):break
    except OSError:time.sleep(.1)
   else:raise RuntimeError('TCP server not ready')
   check(not list((app/'storage/framework/views').glob('*.php')),'view cache not initially empty')
   accept(request('/up'));first=request('/');accept(first)
   jar=SimpleCookie()
   for k,v in first['headers']:
    if k.lower()=='set-cookie':jar.load(v)
   cookie='; '.join(k+'='+v.value for k,v in jar.items())
   with sqlite3.connect(db) as c:initial_sessions=c.execute('select id,payload from sessions').fetchall()
   check(len(initial_sessions)==1,'first welcome must create one session')
   cache={str(f.relative_to(app)):sha(f) for f in (app/'storage/framework/views').glob('*.php')};check(cache,'Blade did not compile views')
   accept(request('/',cookie));check(cache=={str(f.relative_to(app)):sha(f) for f in (app/'storage/framework/views').glob('*.php')},'warm request changed compiled views')
   for i in range(self.args.sequential):accept(request('/' if i%2==0 else '/up',cookie if i%2==0 else ''))
   with ThreadPoolExecutor(max_workers=8) as pool:
    for record in pool.map(request,['/','/up']*4):accept(record)
   check(db_state(db)==before,'HTTP corrupted schema/migrations')
   with sqlite3.connect(db) as c:
    sessions=c.execute('select id,payload,last_activity from sessions order by id').fetchall()
    untouched={t:c.execute('select count(*) from '+t).fetchone()[0] for t in ['users','password_reset_tokens','jobs','job_batches','failed_jobs','cache','cache_locks']}
   check(len(sessions)==5,'cookie replay failed or requests leaked sessions: '+str(len(sessions)))
   check(all(n==0 for n in untouched.values()),'HTTP changed unrelated tables')
   tokens=set()
   for sid,payload,activity in sessions:
    check(re.fullmatch('[A-Za-z0-9]{40}',sid) and abs(activity-time.time())<180,'session id/timestamp')
    decoded=base64.b64decode(payload,validate=True).decode();match=re.search(r's:6:"_token";s:40:"([A-Za-z0-9]{40})";',decoded);check(match,'session CSRF token');tokens.add(match[1])
    if sid==initial_sessions[0][0]:check(payload==initial_sessions[0][1],'session payload changed on cookie replay')
   check(len(tokens)==5,'CSRF token leaked across independent sessions')
   (self.out/(label+'-sessions.json')).write_text(json.dumps({'sessions':sessions,'untouched':untouched},indent=2)+'\n')
   (self.out/(label+'-view-cache.json')).write_text(json.dumps(cache,indent=2)+'\n')
   print('PASS',label,'CLI/migration/HTTP/Blade cold/warm',len(responses),'requests',flush=True)
  finally:
   shutil.copytree(app/'storage/logs',self.out/(label+'-application-logs'),dirs_exist_ok=True)
   shutil.copytree(app/'storage/framework/views',self.out/(label+'-compiled-views'),dirs_exist_ok=True)
   if kind=='official':
    r=subprocess.run([self.args.engine,'rm','-f',container],capture_output=True);entry['container_cleanup_exit_code']=r.returncode
   else:process.terminate()
   try:entry['exit_code']=process.wait(timeout=15)
   except subprocess.TimeoutExpired:process.kill();entry['exit_code']=process.wait()
   log.close();(self.out/'commands.json').write_text(json.dumps(self.commands,indent=2)+'\n')
  serverlog=(self.out/(label+'-server.log')).read_text();check(not any(s in serverlog for s in ['Fatal error','Warning:','Uncaught','Parse error']),'server diagnostic')
  check(all(not p.read_bytes() for p in (self.out/(label+'-application-logs')).glob('*.log')),'application log has diagnostics')
  tokenprobe=work/'blade-tokens.php';tokenprobe.write_text("<?php require __DIR__.'/app/vendor/autoload.php'; $files=[__DIR__.'/app/resources/views/welcome.blade.php',__DIR__.'/app/vendor/laravel/framework/src/Illuminate/Foundation/resources/health-up.blade.php']; foreach(glob(__DIR__.'/app/storage/framework/views/*.php') as $f){$files[]=$f;} $out=[];foreach($files as $f){$text=file_get_contents($f);$tokens=token_get_all($text);$out[]=['text'=>$text,'tokens'=>$tokens];}echo json_encode($out);\n")
  tokenized=json.loads(self.command(label+'-blade-tokens',prefix+[str(tokenprobe)],app,env))
  return {'about':about,'routes':routes,'db':before,'responses':responses,'app_root':str(app),'tokens':tokenized,'session_count':len(sessions)}

def compare_pair(results,out,version):
 def normalize(result):
  root=result['app_root']
  def text(s):return s.replace(root,'{app}').replace(root.lower(),'{app}')
  about=json.loads(text(json.dumps(result['about'])))
  responses=[]
  for r in result['responses']:
   body=r['body']
   if r['path']=='/up':
    check(len(re.findall(r'Response rendered in [0-9]+ms\.',body))==1,'health timing format')
    body=re.sub(r'Response rendered in [0-9]+ms\.','Response rendered in {milliseconds}ms.',body)
   # Transport Date/server/connection/length and randomized encrypted cookies
   # are validated separately; compare the application-owned stable headers.
   headers={k.lower():v for k,v in r['headers'] if k.lower() in {'content-type','cache-control','vary','access-control-allow-origin'}}
   responses.append({'path':r['path'],'status':r['status'],'headers':headers,'body':body})
  tokens=[json.loads(text(json.dumps(t))) for t in result['tokens']]
  tokens.sort(key=lambda x:x['text'])
  return {'about':about,'routes':result['routes'],'db':result['db'],'responses':responses,'tokens':tokens,'session_count':result['session_count']}
 official,go=normalize(results['official']),normalize(results['go'])
 for name,r in [('official',official),('go',go)]:(out/(name+'-'+version+'-normalized.json')).write_text(json.dumps(r,ensure_ascii=False,indent=2)+'\n')
 check(official==go,'official/Go comparison differs; normalized artifacts saved')
 print('MATCH',version,'CLI/DB/HTTP/Blade token streams',flush=True)
def main():
 global APP
 parser=argparse.ArgumentParser(description=__doc__);parser.add_argument('--version',choices=VERSIONS,action='append');parser.add_argument('--runtime',choices=['both','official','go'],default='both');parser.add_argument('--binary',type=Path,default=ROOT/'bin/gophp');parser.add_argument('--engine',default='podman',choices=['podman','docker']);parser.add_argument('--sequential',type=int,default=100)
 parser.add_argument('--app',type=Path,default=APP)
 args=parser.parse_args();args.binary=args.binary.resolve();args.engine=shutil.which(args.engine) or args.engine;APP=args.app.resolve();check(args.sequential>=1,'sequential must be positive')
 out=Path(tempfile.mkdtemp(prefix='fresh-laravel-tcp-',dir=ROOT/'docs/validation'));print('evidence:',out,flush=True)
 run=Run(out,args);code=0
 sys.path.insert(0,str(ROOT/'scripts'))
 from release_evidence import source_identity
 source=source_identity(ROOT);source['binary_sha256']=sha(args.binary);source['locks']['fresh-laravel']=sha(APP/'composer.lock')
 (out/'source.json').write_text(json.dumps(source,indent=2)+'\n')
 identity={'head':subprocess.check_output(['git','rev-parse','HEAD'],cwd=ROOT,text=True).strip(),'diff_sha256':hashlib.sha256(subprocess.check_output(['git','diff','HEAD','--binary'],cwd=ROOT)).hexdigest(),'binary_sha256':sha(args.binary),'lock_sha256':sha(APP/'composer.lock')};(out/'identity.json').write_text(json.dumps(identity,indent=2)+'\n')
 try:
  with tempfile.TemporaryDirectory(prefix='run-',dir=ROOT/'.cache/fresh-laravel') as temp:
   for v in args.version or VERSIONS:
    results={}
    for runtime in ['official','go'] if args.runtime=='both' else [args.runtime]:results[runtime]=run.runtime(runtime,v,temp)
    (out/('summary-'+v+'.json')).write_text(json.dumps(results,ensure_ascii=False,indent=2)+'\n')
    if args.runtime=='both':compare_pair(results,out,v)
 except Exception:traceback.print_exc();code=1
 (out/'result.json').write_text(json.dumps({'command':sys.argv,'exit_code':code},indent=2)+'\n');return code
if __name__=='__main__':raise SystemExit(main())
