<?php
$inf = 1e308 * 10; $nan = $inf - $inf;
do { break; } while ($nan);
for (; true; !$nan) { break; }
while (true || $nan) { break; }
echo 'done', "\n";
