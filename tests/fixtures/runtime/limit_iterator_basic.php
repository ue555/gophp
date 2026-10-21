<?php
$inner = new ArrayIterator(['a' => 1, 'b' => 2, 'c' => 3]);
$limited = new LimitIterator($inner, 1, 1);
foreach ($limited as $key => $value) echo $key, ':', $value, ';';
echo '|';
foreach (new LimitIterator(new ArrayIterator([10, 20, 30]), 2) as $key => $value) echo $key, ':', $value, ';';
echo '|';
foreach (new LimitIterator(new ArrayIterator([10, 20, 30]), 0, 0) as $value) echo 'bad';
echo 'done';
