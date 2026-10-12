<?php
$a = ['nested' => ['x' => 1], 'list' => [2]];
$b = $a; $b['nested']['x'] = 9; $b['list'][] = 3;
echo $a['nested']['x'], $b['nested']['x'], ($a['list'] === [2]) + 0, ($b['list'] === [2, 3]) + 0;
$c['nested'][] = 4; $c['nested'][] = 5; $c[][] = 6;
echo ($c === ['nested' => [4, 5], 0 => [6]]) + 0, "\n";
