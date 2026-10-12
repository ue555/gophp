<?php
echo (($a = $b = 7) === 7) + 0, $a, $b;
echo (($s = 'x') === 'x') + 0, (($n = null) === null) + 0;
echo (($f = 1.5) === 1.5) + 0, (($t = false) === false) + 0;
echo (($arr = [1]) === [1]) + 0, "\n";
