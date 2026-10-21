<?php
class TypedOperations { public int $value = 2; }
$box = new TypedOperations;
$alias =& $box->value;
$result = ($alias += '3');
echo $result, '|', $box->value, '|';
echo ++$alias, '|', $box->value, '|';
$alias ??= 99;
echo $box->value;
