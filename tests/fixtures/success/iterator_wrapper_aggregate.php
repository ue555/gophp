<?php
$source = new ArrayObject(['x' => 4, 'y' => 5]);
$wrapper = new IteratorIterator($source);
echo ($wrapper->getInnerIterator() instanceof Iterator) ? 'inner|' : 'wrong|';
foreach ($wrapper as $key => $value) echo $key, ':', $value, ';';
