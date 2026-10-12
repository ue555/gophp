<?php
$a = [1, 2]; foreach ($a as $k => $v) { echo $k, ':', $v, ';'; unset($k, $v); }
echo isset($k) + 0, isset($v) + 0, ':', count($a), "\n";
foreach ([3] as $v) { unset($a); } echo $v, ':', isset($a) + 0, "\n";
