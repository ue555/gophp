<?php
function &arrayValues(&$values) {
    foreach ($values as &$item) {
        yield $item;
    }
}
$items = [1, 2];
foreach (arrayValues($items) as &$item) {
    $item += 10;
}
unset($item);
echo $items[0], '|', $items[1];
