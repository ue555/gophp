<?php
class RefPrivateBase {
    private $value = 1;
    public function setBase() { $alias =& $this->value; $alias = 4; }
    public function baseValue() { return $this->value; }
}
class RefPrivateChild extends RefPrivateBase {
    private $value = 2;
    public function setChild() { $alias =& $this->value; $alias = 5; }
    public function childValue() { return $this->value; }
}
$child = new RefPrivateChild;
$child->setBase();
$child->setChild();
echo $child->baseValue(), '|', $child->childValue();
