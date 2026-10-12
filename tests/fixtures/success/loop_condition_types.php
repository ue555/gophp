<?php
$values = [null, false, 0, 0.0, -0.0, '', '0', true, -1, 0.5, ' ', '00', [], [1], new stdClass];
for ($i = 0; $i < 15; ++$i) {
    $count = 0; while ($values[$i]) { ++$count; break; }
    echo $count;
}
echo "\n";
