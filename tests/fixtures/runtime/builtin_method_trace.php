<?php
class C { function size($a) { $a[] = 3; return count($a, 2); } }
$o = new C; $o->size([]);
