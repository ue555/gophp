<?php
$o = new Later;
echo $o->read(), "\n";
class Later { public $n = 7; function read() { return $this->n; } }
