<?php
$i = 0; $checks = 0; $log = '';
do {
    ++$i;
    for ($j = 0; true; $missing++) {
        while (true) { $log .= $i; continue 3; }
    }
} while (++$checks < 3);
echo $log, ':', $checks, "\n";
for (;;) { do { while (true) { break 3; } } while (1 / 0); }
echo 'done', "\n";
