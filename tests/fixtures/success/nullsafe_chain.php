<?php
class Box { public $next = null; function get($x) { echo "call"; return $this; } }
$b = new Box; $n = null;
echo ($n?->get(1 / 0)?->next ?? "nil"), "|";
echo ($b?->get(2)?->next ?? "end"), "\n";
