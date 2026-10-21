<?php
echo random_int(5,5),':',random_int(PHP_INT_MIN,PHP_INT_MIN),':',random_int(PHP_INT_MAX,PHP_INT_MAX),"\n";
$ok=true;
for($i=0;$i<100;$i++){ $x=random_int(-7,8); $ok=$ok&&is_int($x)&&$x>=-7&&$x<=8; $x=random_int(PHP_INT_MIN,PHP_INT_MAX);$ok=$ok&&is_int($x); }
echo $ok ? "bounded\n":"bad\n";
try {random_int(2,1);} catch(Throwable $e){echo $e->getMessage(),"\n";}
try {random_int([],1);} catch(Throwable $e){echo $e->getMessage(),"\n";}
