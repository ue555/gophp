<?php
class PrivateReadonlyBase { private readonly int $value; public function read() { return $this->value; } }
class PrivateReadonlyChild extends PrivateReadonlyBase {
    public function initialize() { $this->value = 4; }
}
$o = new PrivateReadonlyChild;
try { $o->initialize(); echo 'set'; } catch (Error $e) { echo 'blocked'; }
echo "\n";
