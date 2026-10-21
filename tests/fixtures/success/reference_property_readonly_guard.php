<?php
class RefReadonlyGuard {
    public readonly int $value;
    public function __construct() { $this->value = 3; }
}
$box = new RefReadonlyGuard;
try { $alias =& $box->value; } catch (Error $e) { echo 'blocked|'; }
echo $box->value;
