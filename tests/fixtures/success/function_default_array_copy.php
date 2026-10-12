<?php
function appendDefault($a = ['x' => [1]]) { $a['x'][] = 2; return $a; }
echo (appendDefault() === ['x' => [1, 2]]) + 0, (appendDefault() === ['x' => [1, 2]]) + 0, "\n";
$a = ['x' => [5]]; $b = appendDefault($a); $b['x'][0]++;
echo ($a === ['x' => [5]]) + 0, ($b === ['x' => [6, 2]]) + 0, "\n";
