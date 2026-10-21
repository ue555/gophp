<?php
function &countDown() {
    $value = 3;
    while ($value > 0) {
        yield $value;
    }
}
foreach (countDown() as &$number) {
    echo --$number, '|';
}
unset($number);
