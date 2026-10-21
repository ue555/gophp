<?php
class TypedRebind { public int $value = 1; }
$box = new TypedRebind;
$old =& $box->value;
$next = 2;
$box->value =& $next;
$old = 'free';
$next = '5';
echo $old, '|', $box->value, '|', gettype($next);
