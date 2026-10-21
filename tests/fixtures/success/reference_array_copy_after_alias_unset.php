<?php
$a = [1];
$x =& $a[0];
unset($x);
$b = $a;
$b[0] = 4;
echo $a[0], '|', $b[0];
