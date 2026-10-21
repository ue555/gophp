<?php
$values = [1, 2];
foreach ($values as &$item) {
    if ($item === 1) {
        $values[] = 3;
    }
    $item++;
}
unset($item);
echo $values[0], '|', $values[1], '|', $values[2];
