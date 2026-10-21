<?php
$v = 2;
$a = [];
$a[] =& $v;
$a[] = 4;
$v = 3;
echo $a[0], '|', $a[1], '|';
$b = ['inner' => []];
$x =& $b['inner']['k'];
echo gettype($x), '|';
$x = 9;
$c = $b;
$c['inner']['k'] = 10;
echo $b['inner']['k'], '|', $c['inner']['k'];
