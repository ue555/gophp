<?php
$a = null; echo ($a .= true), ':'; $a .= false; $a .= '0'; echo $a, "\n";
$a = 12; echo ($a .= 2 + 3), ':', ($a === '125') + 0, "\n";
$a = 'v='; echo ($a .= 7 % 4), "\n";
