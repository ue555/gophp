<?php
class NoConstructor { public $n = 5; }
$x = 1; $o = new NoConstructor($x = 2, $x = 3);
echo $x, ':', $o->n, "\n";
