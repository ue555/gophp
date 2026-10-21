<?php
class TypedArraySource { public int $value = 1; }
$box = new TypedArraySource;
$array = [2];
$box->value =& $array[0];
$array[0] = '7';
echo $box->value, '|', gettype($array[0]);
