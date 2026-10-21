<?php
class ArrayDynamicBox { public array $items = [1, 2, 3]; }
$box = new ArrayDynamicBox;
$name = 'items';
unset($box->{$name}[1]);
echo isset($box->{$name}[1]) ? 'yes' : 'no';
echo '|', $box->{$name}[0], ',', $box->{$name}[2], '|';
$box->{$name}[0] += 4;
echo $box->items[0], "\n";
