<?php
function drop($a, $o) { unset($a[0], $o->x); return $a; }
$a = [3]; $o = new stdClass; $o->x = 4; $b = drop($a, $o);
echo count($a), ':', count($b), ':', isset($o->x) + 0, "\n";
