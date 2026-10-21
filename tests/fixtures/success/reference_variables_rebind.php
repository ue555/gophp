<?php
$a = 1;
$b =& $a;
$c =& $a;
$other = 9;
$b =& $other;
$b = 10;
echo $a, '|', $c, '|', $other, '|';
unset($c);
echo isset($c) ? 'set|' : 'unset|';
$a = 2;
$c = 7;
echo $a, '|', $c, '|', $b;
