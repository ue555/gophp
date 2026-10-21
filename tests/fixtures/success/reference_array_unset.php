<?php
$a = [3];
$x =& $a[0];
unset($a[0]);
echo $x, '|', isset($a[0]) ? 'yes' : 'no', '|';
$a[0] = 8;
$x = 5;
echo $a[0], '|', $x, '|';
unset($x);
echo $a[0];
