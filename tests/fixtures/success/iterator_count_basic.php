<?php
echo iterator_count([1, 2, 3]), '|';
$items = new ArrayIterator(['x' => 1, 'y' => 2]);
echo iterator_count($items), '|';
echo gettype($items->current()), '|';
echo iterator_count(new LimitIterator(new ArrayIterator([4, 5, 6]), 1, 1));
