<?php
class HiddenParent {
    private $value = 1;
    public function parentValue() { return $this->value; }
}
class PublicChild extends HiddenParent {
    public $value = 2;
    public function childValue() { return $this->value; }
}
$object = new PublicChild;
echo $object->parentValue(), '|', $object->childValue(), '|', $object->value, '|';
$object->value = 9;
echo $object->parentValue(), '|', $object->childValue(), "\n";
