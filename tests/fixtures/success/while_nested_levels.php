<?php
$i = 0; $seen = '';
while (++$i <= 3) {
    $j = 0;
    while (++$j <= 3) {
        if ($j === 2) continue 2;
        $seen .= $i . $j;
    }
    echo 'bad';
}
echo $seen, ':', $i, "\n";
while (true) { while (true) { echo 'stop'; break 2; } echo 'bad'; }
echo "\n";
