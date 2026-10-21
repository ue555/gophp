<?php
class ReadonlyParentScope { public readonly int $value; }
class ReadonlyChildScope extends ReadonlyParentScope {
    public function initialize() { $this->value = 8; }
}
$c = new ReadonlyChildScope;
try { $c->initialize(); echo $c->value; } catch (Error $e) { echo 'blocked'; }
echo "\n";
