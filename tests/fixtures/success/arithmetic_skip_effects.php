<?php
$a = 1; false && ($a += 1 / 0); true || ++$missing;
if (false) { $a %= 0; --$other; }
false && ($a[] += 1); true || ($a .= 'bad');
echo $a, "\n";
