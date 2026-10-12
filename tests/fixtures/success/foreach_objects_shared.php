<?php
$o = new stdClass; $o->n = 1; $a = [$o, $o];
foreach ($a as $v) { ++$v->n; $v = null; }
echo $o->n, ':', $a[0]->n, ':', is_null($v) + 0, "\n";
