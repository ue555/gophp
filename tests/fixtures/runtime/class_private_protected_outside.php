<?php
class OuterHidden {
    private $value = 1;
    public function parentValue() { return $this->value; }
}
class OuterProtected extends OuterHidden {
    protected $value = 2;
    public function childValue() { return $this->value; }
}
$object = new OuterProtected;
echo $object->parentValue(), '|', $object->childValue(), "\n";
echo $object->value;
