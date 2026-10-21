<?php
class ParentPrivateBox {
    private int $value = 2;
    public function parentValue() { $name = 'value'; return $this->{$name}; }
    public function changeParent() { $name = 'value'; $this->$name = 4; }
}
class ChildPrivateBox extends ParentPrivateBox {
    private int $value = 3;
    public function childValue() { $name = 'value'; return $this->{$name}; }
    public function changeChild() { $name = 'value'; $this->$name = 7; }
}
$box = new ChildPrivateBox;
echo $box->parentValue(), '|', $box->childValue(), '|';
$box->changeParent(); $box->changeChild();
echo $box->parentValue(), '|', $box->childValue(), "\n";
