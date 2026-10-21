<?php declare(strict_types=1);
require __DIR__ . '/../support/typed_property_defs.inc';
$object = new CrossFileTyped;
try { $object->count = '3'; } catch (TypeError $error) { echo 'strict|'; }
echo $object->count, '|';
weakTypedWrite($object);
echo $object->count;
