<?php
$inf = 1e308 * 10; $nan = $inf - $inf;
while ($nan) { echo 'w'; break; }
for (; $nan; ) { echo 'f'; break; }
$i = 0; do { ++$i; if ($i === 2) break; echo 'd'; } while ($nan);
echo ':', $i, "\n";
