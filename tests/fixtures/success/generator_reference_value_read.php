<?php
function &sharedValue() {
    global $stored;
    yield $stored;
}
$stored = 1;
$g = sharedValue();
$copy = $g->current();
$copy = 9;
echo $stored, '|', gettype($g->current()), '|';
foreach (sharedValue() as $copy) {
    $copy = 10;
}
echo $stored;
