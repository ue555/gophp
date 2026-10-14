<?php
$x = 2;
$f = function ($n) use ($x) { return $n + $x; };
$x = 10;
$arrow = fn ($n) => $n + $x;
$x = 100;
echo $f(3), '|', $arrow(1), '|', $x, "\n";
