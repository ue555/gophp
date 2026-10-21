<?php
class RefArrayBox { public $value = 1; }
$box = new RefArrayBox;
$items = [0];
$items[0] =& $box->value;
$items[0] = 4;
echo $box->value, '|';
$box->value = 6;
echo $items[0];
