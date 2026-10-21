<?php
class ProtectedParent { protected int $value = 9; }
class ProtectedChild extends ProtectedParent {
    public function readValue() { $name = 'value'; return $this->{$name}; }
}
$box = new ProtectedChild;
echo $box->readValue(), '|';
$name = 'value';
echo isset($box->{$name}) ? 'yes' : 'no';
try { echo $box->{$name}; } catch (Error $e) { echo '|Error'; }
echo "\n";
