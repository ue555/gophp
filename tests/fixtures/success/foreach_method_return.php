<?php
class C { public $calls = 0; function first($rows) { foreach ($rows as $row) { ++$this->calls; if (!empty($row)) return $row; } } }
$o = new C; $a = [[], [2, 3], [4]]; $b = $o->first($a); $b[] = 9;
echo $o->calls, ':', count($a[1]), ':', count($b), "\n";
