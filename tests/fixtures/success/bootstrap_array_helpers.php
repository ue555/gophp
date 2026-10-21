<?php
$values = ['x', '', 'x', 'y'];
$filtered = array_values(array_filter(array_unique($values)));
echo implode(',', $filtered), '|';
echo ltrim('/tmp', '/'), '|', rtrim('/tmp///', '/'), '|';
$numbers = [10, 20, 30];
echo array_shift($numbers), '|', current($numbers), '|', end($numbers), '|';
$tail = array_slice($numbers, -1);
echo current($tail);
