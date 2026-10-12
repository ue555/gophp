<?php
$a = 1; $b = 2; unset($a, $b, $missing);
echo isset($a) + 0, isset($b) + 0, empty($missing) + 0, "\n";
$a = [10, 20]; $i = 0; unset($a[$i++], $a[$i++]);
echo count($a), ':', $i, "\n";
