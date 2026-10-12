<?php
function update($a) { $a['2'] += 1; $a[false] = 'new'; $a[] = 'last'; return $a; }
$a = [2 => 7, 0 => 'old', '02' => 8]; $b = update($a);
echo ($a === [2 => 7, 0 => 'old', '02' => 8]) + 0;
echo ($b === [2 => 8, 0 => 'new', '02' => 8, 3 => 'last']) + 0, "\n";
