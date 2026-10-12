<?php
$i = 0; $sum = 0;
while ($i < 8) {
    ++$i; if ($i % 2 === 0) continue;
    if ($i === 5) break;
    $sum += $i;
}
echo $i, ':', $sum, "\n";
