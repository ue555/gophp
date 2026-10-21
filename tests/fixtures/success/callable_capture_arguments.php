<?php
$value = 4;
$closure = function (int $n) use ($value): int { return $value + $n; };
$value = 9;
$arrow = fn(int $n): int => $value + $n;
$box = [$closure, $arrow];
echo $box[0](2), '|', $box[0](3), '|', $box[1](1), '|';
function namedPair(string $first, string $second = 'b'): string { return $first . $second; }
$first = namedPair(...);
echo $first(second: '2', first: '1'), '|', call_user_func_array('namedPair', ['second' => 'Y', 'first' => 'X']);
