<?php
$a = 1;
echo ($b =& $a), '|';
$a = 2;
echo $b, '|';
$b = 3;
echo $a, '|';
$copy = $b;
$b = 4;
echo $copy, '|', $a;
