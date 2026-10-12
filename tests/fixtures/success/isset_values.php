<?php
echo isset($missing) + 0, ':', empty($missing) + 0, ':';
$n = null; $f = false; $z = 0; $s = ''; $a = [];
echo isset($n) + 0, isset($f, $z, $s, $a) + 0, ':', isset($a[0]) + 0, "\n";
$a = ['n' => null, 'f' => false, 'a' => [1]];
echo isset($a['n']) + 0, isset($a['f']) + 0, isset($a['a'][0]) + 0;
echo isset($missing['x']->p) + 0, empty($a['missing']['x']) + 0, "\n";
