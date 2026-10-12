<?php
class Adder {
    function add($a, $b = 2, $c = [3]) { $c[] = 4; return $a + $b + $c[0] + $c[1]; }
}
$o = new Adder;
echo $o->add(1), ':', $o->add(1, 5, [6]), ':', $o->add(1), "\n";
