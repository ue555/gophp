<?php
class RefLife { public $value = 1; }
$first = new RefLife;
$alias =& $first->value;
unset($first->value);
$alias = 3;
$first->value = 5;
echo $alias, '|', $first->value, '|';
$second = new RefLife;
$held =& $first->value;
$first = $second;
$held = 8;
echo $first->value, '|', $held;
