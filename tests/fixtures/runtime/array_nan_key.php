<?php
$inf = 1e308 * 10; $nan = $inf - $inf;
$a = [$nan => 2];
echo $a[0], "\n";
