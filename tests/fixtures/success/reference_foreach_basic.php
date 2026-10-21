<?php
$values = [1, 2, 3];
foreach ($values as $key => &$item) {
    $item += $key;
}
echo $values[0], '|', $values[1], '|', $values[2], '|';
$item = 9;
echo $values[2];
unset($item);
