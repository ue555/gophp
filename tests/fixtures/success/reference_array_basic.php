<?php
$a = [1];
$x =& $a[0];
$x = 2;
$a[0] = 3;
$copy = $a[0];
$copy = 8;
echo $x, '|', $a[0], '|', $copy, '|', gettype($a[0]);
