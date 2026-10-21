<?php declare(strict_types=1);
class ReadonlyTypedState {
    public readonly int $value;
    public function initialize($value) { $this->value = $value; }
}
$o = new ReadonlyTypedState;
try { $o->initialize('bad'); } catch (TypeError $e) { echo 'type|'; }
echo isset($o->value) ? 'set|' : 'unset|';
$o->initialize(12);
echo $o->value, "\n";
