<?php
$inner = new ArrayIterator(['first' => 1, 'second' => 2]);
$wrapper = new IteratorIterator($inner);
echo ($wrapper->getInnerIterator() === $inner) ? 'same|' : 'different|';
foreach ($wrapper as $key => $value) echo $key, ':', $value, ';';
echo '|';
foreach ($wrapper as $key => $value) echo $key, ':', $value, ';';
