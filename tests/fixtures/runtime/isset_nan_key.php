<?php
$inf = 1e308 * 10; $nan = $inf - $inf; $a = [1];
echo isset($a[$nan]) + 0, "\n";
