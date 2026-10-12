<?php
function square($x) { return $x * $x; }
class Box {
    private $v = 3;
    function value() { return square($this->v); }
    function mixed($n) { return combineBox($this, $n); }
}
function combineBox($o, $n) { return $o->value() + $n; }
$o = new Box; echo $o->mixed(4), ':', $o->value(), "\n";
