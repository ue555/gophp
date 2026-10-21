<?php
class TypedReadOnly { public readonly int $value; public function __construct() { $this->value = 1; } }
$box = new TypedReadOnly;
try { $alias =& $box->value; } catch (Error $e) { echo 'blocked|'; }
echo $box->value;
