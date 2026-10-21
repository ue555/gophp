<?php
class ReadonlyUninitializedUnset {
    public readonly int $value;
    public function resetBeforeFirstWrite() { unset($this->value); $this->value = 8; }
}
$o = new ReadonlyUninitializedUnset;
try { unset($o->value); } catch (Error $e) { echo 'outside|'; }
$o->resetBeforeFirstWrite();
echo $o->value, "\n";
