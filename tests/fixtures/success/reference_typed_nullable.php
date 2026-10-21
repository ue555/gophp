<?php
class TypedNullable { public ?int $value; }
$box = new TypedNullable;
$alias =& $box->value;
echo isset($box->value) ? 'set' : 'null', '|';
$alias = '6';
echo $box->value, '|', gettype($alias);
