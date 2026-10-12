<?php
$x = 100;
function calc($x) { $local = 3; ++$x; return $x + $local; }
echo calc(1), ':', calc(2), ':', $x, "\n";
