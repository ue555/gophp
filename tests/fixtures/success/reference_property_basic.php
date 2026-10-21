<?php
class RefBox { public $value = 1; }
$box = new RefBox;
$alias =& $box->value;
$alias = 4;
echo $box->value, '|';
$box->value = 7;
echo $alias, '|';
$copy = $alias;
$copy = 8;
echo $alias, '|', $copy;
