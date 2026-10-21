<?php
class TypedStaticParent { public static int $value = 1; }
class TypedStaticChild extends TypedStaticParent {}
$alias =& TypedStaticParent::$value;
TypedStaticChild::$value = '7';
echo $alias, '|', TypedStaticParent::$value, '|', TypedStaticChild::$value;
