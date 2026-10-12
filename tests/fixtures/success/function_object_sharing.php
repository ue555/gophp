<?php
class Counter { public $n = 0; function next() { return ++$this->n; } }
function mutate($o, $a) { $o->next(); $a[0] += 10; return [$o, $a]; }
$o = new Counter; $a = [2]; $r = mutate($o, $a); $r[0]->n += 3;
echo $o->n, ':', $a[0], ':', $r[1][0], ':', ($r[0] === $o) + 0, "\n";
