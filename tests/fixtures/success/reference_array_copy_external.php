<?php
$a = [1];
$outside =& $a[0];
$b = $a;
$outside = 4;
echo $a[0], '|', $b[0], '|';
$b[0] = 6;
echo $a[0], '|', $outside;
