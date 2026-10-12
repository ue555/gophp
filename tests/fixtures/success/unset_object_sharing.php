<?php
class C { public $x = 3; function clear() { unset($this->x); } }
$o = new C; $a = [$o]; $b = $a; $b[0]->clear();
echo isset($o->x) + 0, isset($a[0]->x) + 0, ':';
$o->x = 7; echo $a[0]->x, ':', $b[0]->x, "\n";
unset($b[0]); echo count($a), ':', count($b), ':', $o->x, "\n";
