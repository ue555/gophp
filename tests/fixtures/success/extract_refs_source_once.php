<?php
$calls = 0;
$array = ['x' => 1, 'y' => 2];
function &sourceArray() {
    global $calls, $array;
    $calls++;
    return $array;
}
echo extract(sourceArray(), EXTR_REFS), '|', $calls, '|';
$x = 9;
echo $array['x'], '|';
$array['y'] = 7;
echo $y, '|', $calls;
