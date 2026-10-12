<?php
$a = ['x' => [1, 2], 'n' => null]; $b = $a;
unset($b['x'][0], $b['missing']['x'], $b['n'][0]);
echo count($a['x']), ':', count($b['x']), ':', $a['x'][0], ':', isset($b['missing']) + 0, "\n";
$b['x'][] = 3; foreach ($b['x'] as $k => $v) echo $k, ':', $v, ';'; echo "\n";
