<?php
class TypedUnionConflict {
    public int|string $first = 'ready';
    public float|string $second = 'ready';
}
$box = new TypedUnionConflict;
$alias =& $box->first;
$box->second =& $alias;
try { $alias = 42; } catch (TypeError $e) { echo 'conflict|'; }
echo $box->first, '|', $box->second;
