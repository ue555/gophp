<?php
$items = ['one' => 1, 'two' => 2];
$callback = function (&$value, $key, $increment) { $value += $increment; };
echo array_walk(callback: $callback, arg: 3, array: $items) ? 'ok|' : 'bad|';
echo implode(',', $items), '|';

$tree = ['outer' => ['x' => 4], 'y' => 5];
array_walk_recursive(callback: function (&$value, $key) { $value += 2; }, array: $tree);
echo json_encode($tree);
