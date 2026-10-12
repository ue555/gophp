<?php
function total($a) {
    $sum = 0;
    foreach ($a as $k => $v) {
        if (is_array($v)) $sum += total($v);
        else if (is_numeric($v)) $sum += (float)$v;
        unset($a[$k]);
    }
    return $sum;
}
$a = [1, ['2.5', null, [3]], 'skip']; $b = $a;
echo total($a), ':', count($a, COUNT_RECURSIVE), ':', ($a === $b) + 0, "\n";
