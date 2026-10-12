<?php
function build($n, $a = []) {
    if ($n === 0) return $a;
    for ($i = 0; $i < 2; ++$i) { $a[] = $n * 10 + $i; }
    return build($n - 1, $a);
}
$a = [9]; $b = build(3, $a); $c = $b; $c[0]++;
echo ($a === [9]) + 0, ($b === [9, 30, 31, 20, 21, 10, 11]) + 0, ($c[0] === 10) + 0, "\n";
