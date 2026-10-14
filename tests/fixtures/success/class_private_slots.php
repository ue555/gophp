<?php
class SlotParent {
    private $number = 2;
    private $items = [1];
    public function __construct() { $this->number = 8; }
    public function parentNumber() { return $this->number; }
    public function addParent() { $this->number += 3; $this->items[] = 2; }
    public function parentItems() { return count($this->items); }
}
class SlotChild extends SlotParent {
    private $number = 5;
    private $items = [7];
    public function childNumber() { return $this->number; }
    public function addChild() { ++$this->number; $this->items[] = 8; }
    public function childItems() { return count($this->items); }
}
$first = new SlotChild;
$second = new SlotChild;
$shared = $first;
$first->addParent();
$shared->addChild();
echo $first->parentNumber(), '|', $first->childNumber(), '|', $first->parentItems(), '|', $first->childItems(), '|';
echo $second->parentNumber(), '|', $second->childNumber(), "\n";
