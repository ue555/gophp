<?php
class Mutator {
    function run($n, $a, $o) { $n = 9; $a[0] = 9; $o->n = 9; return [$n, $a]; }
}
$n = 1; $a = [2]; $o = new stdClass; $o->n = 3;
$m = new Mutator; $result = $m->run($n, $a, $o);
echo $n, $a[0], $o->n, $result[0], $result[1][0], "\n";
