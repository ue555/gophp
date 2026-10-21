<?php
class TypedBase {
    private int $value = 2;
    public function baseValue() { return $this->value; }
    public function writeBase(int $value) { $this->value = $value; }
}
class TypedChild extends TypedBase {
    private string $value = 'child';
    public function childValue() { return $this->value; }
    public function writeChild(string $value) { $this->value = $value; }
}
$object = new TypedChild;
echo $object->baseValue(), '|', $object->childValue();
$object->writeBase(5);
$object->writeChild('next');
echo '|', $object->baseValue(), '|', $object->childValue();
