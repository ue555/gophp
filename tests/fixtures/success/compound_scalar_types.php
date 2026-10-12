<?php
$a = ' 12 '; $a += '3'; echo ($a === 15) + 0;
$a = false; $a += true; echo ($a === 1) + 0;
$a = null; $a -= 2; echo ($a === -2) + 0;
$a = '2.5'; $a *= 2; echo ($a === 5.0) + 0;
$a = true; $a /= 2; echo ($a === 0.5) + 0;
$a = '10'; $a %= '3'; echo ($a === 1) + 0, "\n";
