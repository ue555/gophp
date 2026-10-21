<?php
class RelativeBase {}
class RelativeChild extends RelativeBase {
    public parent $item;
    public function create(): parent { return new RelativeBase; }
}
$child = new RelativeChild;
$child->item = $child->create();
echo $child->item instanceof RelativeBase ? 'parent' : 'wrong';
