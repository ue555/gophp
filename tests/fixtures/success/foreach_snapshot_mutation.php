<?php
$a = [10, 20, 30];
foreach ($a as $k => $v) { echo $k, ':', $v, ';'; $a[1] = 99; $a[] = 40; $v = -1; }
echo "\n", count($a), ':', $a[1], ':', $k, ':', $v, "\n";
foreach ($a as $v) { $a = []; echo $v, ';'; }
echo ':', count($a), "\n";
