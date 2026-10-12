<?php
$inf = 1e308 * 10; $nan = $inf - $inf; $before = $nan++;
echo ($nan !== $nan) + 0, ($before !== $before) + 0, (++$inf === $inf) + 0, "\n";
