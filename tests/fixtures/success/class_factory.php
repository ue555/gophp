<?php
class Factory { function create($n) { return new Item($n); } }
class Item { public $n; function __construct($n) { $this->n = $n; } }
$f = new Factory; $a = $f->create(3); $b = $f->create(4);
echo $a->n, $b->n, ($a !== $b) + 0, "\n";
