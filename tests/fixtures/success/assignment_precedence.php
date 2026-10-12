<?php
$a = false || true; echo ($a === true) + 0;
$a = 0; true || $a = 1; echo $a;
$b = 0; false && $b = 1; echo $b;
$c = 0; false || $c = 2; echo $c;
$d = 0; true && $d = 3; echo $d;
echo (1 + $e = 4), ':', $e;
echo ':', ($f = 1 + 2 * 3), ':', ($g = 3 <=> 2), "\n";
