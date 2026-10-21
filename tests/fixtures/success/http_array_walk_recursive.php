<?php
$items = ['one' => ['a' => 'A', 'b' => ['B']], 'two' => 'C'];
array_walk_recursive($items, function (&$value, $key) { $value .= $key; });
echo json_encode($items);
