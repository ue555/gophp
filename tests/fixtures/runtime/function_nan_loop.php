<?php
function check($n) { while ($n) { return 1; } return 0; }
$inf = 1e308 * 10; echo check($inf - $inf), "\n";
