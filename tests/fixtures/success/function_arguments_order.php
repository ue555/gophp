<?php
function marker($s) { echo $s; return $s; }
function combine($a, $b, $c = '!') { return $a . $b . $c; }
echo ':', combine(marker('a'), marker('b')), "\n";
$x = 1; function pair($a, $b) { return $a . ':' . $b; }
echo pair($x, $x = 2), ':', $x, "\n";
