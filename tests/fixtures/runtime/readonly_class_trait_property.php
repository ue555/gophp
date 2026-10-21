<?php
trait ReadonlyTraitField {
    public int $value;
    public function initialize() { $this->value = 5; }
}
readonly class ReadonlyTraitUser { use ReadonlyTraitField; }
$o = new ReadonlyTraitUser;
$o->initialize();
echo $o->value, '|';
try { $o->value = 6; } catch (Error $e) { echo 'blocked'; }
echo "\n";
