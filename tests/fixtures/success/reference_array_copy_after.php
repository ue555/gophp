<?php
$a = [1];
$x =& $a[0];
$b = $a;
unset($x);
$b[0] = 4;
echo $a[0], '|', $b[0], '|';
$a[0] = 6;
echo $b[0];
