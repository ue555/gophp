<?php
function &firstElement(&$values) {
    return $values[0];
}
$items = [1, 2];
$alias =& firstElement($items);
$alias = 4;
echo $items[0], '|';
$copy = firstElement($items);
$copy = 9;
echo $items[0];
