<?php
class TypedArrayMember { public array $items = [1]; }
$box = new TypedArrayMember;
$alias =& $box->items[0];
$alias = 3;
echo $box->items[0], '|';
try { $box->items = 'wrong'; } catch (TypeError $e) { echo 'blocked|'; }
echo $alias, '|', gettype($box->items);
