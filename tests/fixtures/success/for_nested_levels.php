<?php
$log = ''; for ($i = 0; $i < 3; $log .= 'U', ++$i) {
    for ($j = 0; $j < 3; ++$j) {
        $log .= $i . $j;
        continue 2;
    }
}
echo $log, ':', $i, ':', $j, "\n";
for (; true; 1 / 0) { while (true) { break 2; } }
echo 'done', "\n";
