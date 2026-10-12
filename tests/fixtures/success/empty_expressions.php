<?php
foreach ([null, false, true, 0, 1, -0.0, '', '0', '00', [], [0], new stdClass] as $v) echo empty($v) + 0;
echo "\n", empty(1 + 2) + 0, empty('') + 0, empty($n = 0) + 0, ':', $n, "\n";
function f() { echo 'F'; return 0; } echo empty(f()) + 0, "\n";
