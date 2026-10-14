<?php
class ParentItem {
    protected $value;
    public function __construct($value) { $this->value = $value; }
    public function read() { return $this->value; }
}
class ChildItem extends ParentItem {}
echo (new ChildItem(9))->read(), "\n";
