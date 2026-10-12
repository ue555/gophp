<?php
$s = 'a0b';
foreach ([0, 1, -1, 3, -4, '0', 'x', '1.0', [], new stdClass] as $k) echo isset($s[$k]) + 0, empty($s[$k]) + 0, ';';
echo "\n";
