<?php
class TypedPending { public int $value; }
$box = new TypedPending;
$source = '3';
$box->value =& $source;
echo $box->value, '|', gettype($source), '|';
$source = 5;
echo $box->value;
