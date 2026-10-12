<?php
class NanMethod {
    function condition($n) { if ($n) return true; return false; }
}
$inf = 1e308 * 10; $nan = $inf - $inf; $o = new NanMethod;
echo $o->condition($nan) + 0, "\n";
