<?php
class TypedDynamic { public int $value = 1; }
$box = new TypedDynamic;
$calls = 0;
function typedName() { global $calls; $calls++; return 'value'; }
$alias =& $box->{typedName()};
$alias = '6';
echo $box->value, '|', $calls;
