<?php
$source = ['b' => 2, 'a' => 1];
asort($source);
echo implode(',', array_keys($source)), '|';
echo implode(',', array_pad([1], 3, 4)), '|';
echo array_search('b', ['a', 'b', 'c']), '|';
echo array_reduce([1, 2, 3], static fn ($carry, $item) => $carry + $item, 0), '|';
$hasTwo = function_exists('array_any')
    ? array_any([1, 2, 3], static fn ($value) => $value === 2)
    : in_array(2, [1, 2, 3], true);
echo $hasTwo ? 'yes' : 'no', '|';
echo implode(',', array_keys(array_fill_keys(['x', 'y'], 9))), '|';
echo json_decode('{"value":7}', true)['value'], '|';
echo rawurldecode('a%2Bb'), '|';
echo gmdate('D, d M Y H:i:s', 0);
