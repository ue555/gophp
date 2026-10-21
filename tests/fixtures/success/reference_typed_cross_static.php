<?php
class TypedCrossStatic { public int $value = 1; public static int $saved = 2; }
$box = new TypedCrossStatic;
$alias =& $box->value;
TypedCrossStatic::$saved =& $alias;
$alias = '6';
echo $box->value, '|', TypedCrossStatic::$saved;
