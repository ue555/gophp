<?php
$inf = 1e308 * 10; $nan = $inf - $inf;
echo $nan <=> 0, ':', 0 <=> $nan, ':', $nan <=> $nan, ':', $inf <=> $inf, ':', -$inf <=> $inf, "\n";
echo ($nan <> $nan) + 0, ([ $nan ] <> [ $nan ]) + 0;
$a = [$nan]; $b = $a;
echo ($a === $b) + 0, ($a == $b) + 0, "\n";
