<?php
$inf = 1e308 * 10; $nan = $inf - $inf;
echo (int)$nan, ':', (int)$inf, ':', (int)(bool)$nan, ':', (string)$nan, "\n";
