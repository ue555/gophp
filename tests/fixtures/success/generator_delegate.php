<?php
function inner() {
    yield 4 => 'four';
    return 8;
}
function outer() {
    $value = yield from inner();
    yield $value;
}
foreach (outer() as $key => $value) {
    echo $key, ':', $value, '|';
}
