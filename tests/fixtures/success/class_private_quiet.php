<?php
class QuietParent {
    private $value = null;
    public function parentState() {
        echo isset($this->value) ? 'yes' : 'no', '|';
        $this->value ??= 4;
        echo $this->value, '|';
        unset($this->value);
        echo empty($this->value) ? 'empty' : 'full', '|';
        $this->value ??= 7;
        echo $this->value, '|';
    }
}
class QuietChild extends QuietParent {
    private $value = 9;
    public function childValue() { return $this->value; }
}
$object = new QuietChild;
echo isset($object?->value) ? 'yes' : 'no', '|';
echo $object?->value ?? 'fallback', '|';
$object->parentState();
echo $object->childValue(), "\n";
