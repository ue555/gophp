<?php
class Factorial {
    function calc($n) { if ($n < 2) return 1; return $n * $this->calc($n - 1); }
}
$o = new Factorial;
echo $o->calc(6), "\n";
