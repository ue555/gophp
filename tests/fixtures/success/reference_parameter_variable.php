<?php
function bump(&$value) { $value += 2; }
$a = 3;
bump($a);
echo $a, '|';
$alias =& $a;
bump(value: $alias);
echo $a, '|', $alias;
