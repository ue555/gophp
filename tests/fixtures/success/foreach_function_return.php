<?php
function first($a) { foreach ($a as $v) { if ($v > 1) return $v; } }
function sum($a) { $n = 0; foreach ($a as $v) $n += $v; $a[0] = 99; return $n; }
$a = [1, 2, 3]; echo first($a), ':', sum($a), ':', $a[0], ':', is_null(first([])) + 0, "\n";
