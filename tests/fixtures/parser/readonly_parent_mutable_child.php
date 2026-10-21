<?php
readonly class ReadonlyAncestor {
    public function __construct(public int $value) {}
}
class MutableDescendant extends ReadonlyAncestor {
    public int $counter = 0;
    public function bump() { $this->counter++; }
}
$o = new MutableDescendant(3);
$o->bump();
try { $o->value = 4; } catch (Error $e) { echo 'readonly|'; }
echo $o->value, '|', $o->counter, "\n";
