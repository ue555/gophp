<?php
class TypedArray {
    public array $items = [];
    public static array $shared = [];
}
$x = new TypedArray;
$y = $x;
$x->items[] = 1;
$y->items['key'] = 2;
TypedArray::$shared[] = 3;
TypedArray::$shared[] = 4;
echo count($x->items), '|', $x->items['key'], '|', count(TypedArray::$shared);
unset($x->items['key']);
echo '|', count($x->items);
