<?php
function show($v) { echo json_encode($v), "\n"; }
$base = [2=>'two', 10=>'ten', '02'=>'zero', 'x10'=>'x10', 'x2'=>'x2', 'A'=>'A', 'a'=>'a', -1=>'negative'];
foreach ([SORT_REGULAR, SORT_NUMERIC, SORT_STRING, SORT_NATURAL, SORT_STRING|SORT_FLAG_CASE, SORT_LOCALE_STRING] as $flag) {
 $a=$base; show(krsort($a,$flag)); echo json_encode(array_keys($a)), "\n";
}
$box=[['a'=>1,'c'=>3,'b'=>2]]; $i=0; $alias=&$box[0];
show(krsort($box[$i++])); echo $i,':',json_encode($alias),"\n";
krsort(flags: SORT_STRING, array: $alias); echo json_encode($box),"\n";
try { krsort($alias, []); } catch (Throwable $e) { echo $e->getMessage(),"\n"; }
echo json_encode($alias),"\n";
$a=['x2'=>1,'x02'=>2,'x000'=>3,'x0'=>4,'x10'=>5,'x1'=>6];krsort($a,SORT_NATURAL);echo json_encode(array_keys($a)),"\n";
$x=17;try{krsort($x);}catch(Throwable $e){echo $e->getMessage(),"\n";}echo $x,"\n";
try{krsort(['a'=>1]);}catch(Throwable $e){echo $e->getMessage(),"\n";}
