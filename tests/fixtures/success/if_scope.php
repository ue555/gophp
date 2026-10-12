<?php
$x = 1;
if ($x === 1) {
    $x = 2;
    $created = 'visible';
    if ($x === 2) { $x = $x + 3; }
} else {
    $x = 99;
}
if ($x !== 5) { $x = 88; } else { $created = $created . ':else'; }
echo $x, '|', $created, "\n";
