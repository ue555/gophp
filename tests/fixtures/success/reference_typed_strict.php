<?php declare(strict_types=1);
class TypedStrict { public int $value = 2; }
$box = new TypedStrict;
$alias =& $box->value;
try { $alias = '3'; } catch (TypeError $e) { echo 'rejected|'; }
echo $box->value, '|';
$alias = 4;
echo $box->value;
