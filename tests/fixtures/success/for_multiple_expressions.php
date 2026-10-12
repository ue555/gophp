<?php
$log = '';
for ($i = 0, $log .= 'I'; $log .= 'C', $i < 2; $log .= 'U', ++$i) { $log .= 'B'; }
echo $log, ':', $i, "\n";
for ($a = 0, $b = 5; false, $a < 2; ++$a, --$b) {}
echo $a, ':', $b, "\n";
