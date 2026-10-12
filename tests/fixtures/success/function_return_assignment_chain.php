<?php
function adjust($a) { $a['x'] += 2; return $a; }
$a = $b = adjust(['x' => 1]); $a['x']++;
echo $a['x'], ':', $b['x'], "\n";
function returnAssignment() { return $a = $b = [2]; }
echo (returnAssignment() === [2]) + 0, "\n";
