<?php
class Arguments { function pair($a, $b) { return [$a, $b]; } }
$o = new Arguments; $x = 1;
$a = $o->pair($x, $x = 2);
echo $a[0], $a[1], $x;
$a = $o->pair($x = 3, $x = 4);
echo $a[0], $a[1], $x, "\n";
