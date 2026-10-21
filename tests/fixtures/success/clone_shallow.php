<?php
class CloneBox {
    public int $value = 1;
    public object $child;
    public function __clone(): void { $this->value++; }
}
$original = new CloneBox;
$original->child = new stdClass;
$copy = clone $original;
$copy->child->value = 3;
echo $original->value, '|', $copy->value, '|', $original->child->value;
