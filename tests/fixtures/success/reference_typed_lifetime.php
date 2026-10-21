<?php
class TypedLifetime { public int $value = 1; }
$box = new TypedLifetime;
$alias =& $box->value;
unset($box);
$alias = 'free';
echo $alias;
