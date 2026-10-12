<?php
class Chain {
    public $n = 0;
    function set($n) { $this->n = $n; return $this; }
    function get() { return [$this->n]; }
}
$o = new Chain;
echo $o->set(3)->set(4)->get()[0], ':', $o->n, "\n";
