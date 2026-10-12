<?php
$inf = 1e308 * 10;
$nan = $inf - $inf;
echo (0.0 === -0.0) + 0;
echo (0 === -0.0) + 0;
echo ($inf === $inf) + 0;
echo ($inf === -$inf) + 0;
echo ($nan === $nan) + 0;
echo ($nan !== $nan) + 0;
echo (5e-324 !== 0.0) + 0;
echo (1e-9999 === 0.0) + 0;
echo "\n";
