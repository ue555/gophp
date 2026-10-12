<?php
$inf = 1e308 * 10; $nan = $inf - $inf;
foreach ([0, 2.5, $inf, $nan, '0', ' +2.5e2 ', "\t12\n", '12tail', '', ' ', 'INF', 'NAN', '0x10', '1_000', false, null, [], new stdClass, '1e999', "\xc2\xa012"] as $v) echo (int)is_numeric($v);
echo "\n";
