<?php
class TypedArray { public int $value = 1; }
$box = new TypedArray;
$array = [];
$array[0] =& $box->value;
$copy = $array;
$copy[0] = '9';
echo $box->value, '|', $array[0], '|', $copy[0];
