<?php
class RefOwner { public $alias; public static $saved; }
$array = [1];
$owner = new RefOwner;
$owner->alias =& $array[0];
$copy = $array;
$copy[0] = 4;
echo $array[0], '|', $owner->alias, '|';
$other = [2];
RefOwner::$saved =& $other[0];
$clone = $other;
$clone[0] = 6;
echo $other[0], '|', RefOwner::$saved;
