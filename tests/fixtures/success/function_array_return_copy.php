<?php
function copyArray($a) { return $a; }
$a = [[1]]; $b = copyArray($a); $c = $b;
$b[0][0] += 4; $c[0][] = 6;
echo ($a === [[1]]) + 0, ($b === [[5]]) + 0, ($c === [[1, 6]]) + 0, ':', copyArray($a)[0][0], "\n";
