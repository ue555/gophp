<?php
foreach ([SORT_REGULAR,SORT_NUMERIC,SORT_STRING,SORT_NATURAL,SORT_STRING|SORT_FLAG_CASE] as $flag) {
 $a=['first'=>2,10,'02','x10','x2','A','a',-1]; echo json_encode(rsort($a,$flag)),':',json_encode($a),"\n";
}
$a=[[3,1,2]];$i=0;$ref=&$a[0];rsort($a[$i++]);echo $i,':',json_encode($ref),"\n";
try { rsort($ref,[]); } catch(Throwable $e) {echo $e->getMessage(),"\n";}
echo json_encode($ref),"\n";
