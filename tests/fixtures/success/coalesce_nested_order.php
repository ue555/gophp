<?php
class Counter { public $n = 0; } $c = new Counter; $a = [];
function key_of($c) { $c->n++; return "x"; }
echo $a[key_of($c)] ?? 7, "|", $c->n, "|";
$a["x"] = 9; echo $a[key_of($c)] ?? 7, "|", $c->n, "\n";
