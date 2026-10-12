<?php
$inf = 1e308 * 10; $nan = $inf - $inf; $a = [1];
unset($a[$nan]); echo count($a), "\n";
