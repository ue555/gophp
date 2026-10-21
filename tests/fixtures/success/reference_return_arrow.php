<?php
$value = 2;
$arrow = fn&(&$x) => $x;
$alias =& $arrow($value);
$alias = 6;
echo $value;
