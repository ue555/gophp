<?php
$o = new stdClass; $o->a = [1, 2]; $a = [1, [2, [3, 4]], [], $o];
echo count([]), ':', count($a), ':', count($a, COUNT_NORMAL), ':', count($a, COUNT_RECURSIVE), "\n";
echo count($a, true), ':', count($a, false), ':', count($a, '1'), ':', count($a, 1.0), "\n";
