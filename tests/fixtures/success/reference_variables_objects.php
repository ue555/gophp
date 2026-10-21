<?php
$object = new stdClass;
$object->value = 1;
$copy = $object;
$alias =& $object;
$alias->value = 2;
echo $copy->value, '|';
$alias = new stdClass;
$alias->value = 3;
echo $object->value, '|', $copy->value;
