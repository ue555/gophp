import pathlib,subprocess,json,datetime,sys,uuid,os
ROOT=pathlib.Path(__file__).resolve().parent
def run(label,args,cwd=None,env=None):
 base=ROOT/(datetime.datetime.now().strftime('%Y%m%d-%H%M%S')+'-'+label+'-'+uuid.uuid4().hex[:8])
 p=subprocess.run(args,cwd=cwd,env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
 pathlib.Path(str(base)+'.stdout').write_bytes(p.stdout); pathlib.Path(str(base)+'.stderr').write_bytes(p.stderr)
 pathlib.Path(str(base)+'.json').write_text(json.dumps(dict(command=args,cwd=str(cwd or os.getcwd()),exit_code=p.returncode),indent=2)+'\n')
 print(label, 'exit',p.returncode, 'log',base,flush=True);print(p.stdout.decode(errors='replace'),end='',flush=True);print(p.stderr.decode(errors='replace'),end='',flush=True)
 return p
if __name__=='__main__':
 p=run(sys.argv[1],sys.argv[2:]);sys.exit(p.returncode if p.returncode>=0 else 128-p.returncode)
