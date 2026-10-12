<?php
class TraceArguments {
    function inner($s, $f, $a, $o, $b, $n) { return 1 / 0; }
    function outer($n) {
        $n = 9;
        return $this->inner("a\nb\tc\\d", 1.0, [2], $this, false, null);
    }
}
$o = new TraceArguments;
$o->outer(3);
