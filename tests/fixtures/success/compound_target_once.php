<?php
class Box {
    public $data = [10, 20]; public $calls = 0;
    function target() { ++$this->calls; echo 'T'; return $this; }
    function key() { ++$this->calls; echo 'K'; return 0; }
    function rhs() { ++$this->calls; echo 'R'; return 3; }
}
$b = new Box;
$b->target()->data[$b->key()] += $b->rhs();
echo ':', $b->data[0], ':', $b->calls, "\n";
echo $b->target()->data[$b->key()]++, ':', $b->data[0], ':', $b->calls, "\n";
