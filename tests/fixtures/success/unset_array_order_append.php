<?php
$a = [2 => 'a', 'x' => 'b', 5 => 'c']; unset($a[5], $a['missing']); $a[] = 'd';
unset($a['x']); $a['x'] = 'e';
foreach ($a as $k => $v) echo $k, ':', $v, ';'; echo "\n";
$a = [-3 => 1]; unset($a[-3]); $a[] = 2; foreach ($a as $k => $v) echo $k, ':', $v, "\n";
