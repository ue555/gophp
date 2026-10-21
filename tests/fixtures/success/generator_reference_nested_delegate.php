<?php
function &outerReferenceGenerator() {
    $value = 1;
    $nested = function () { yield from [2]; };
    foreach ($nested() as $unused) {}
    yield $value;
}
foreach (outerReferenceGenerator() as &$item) {
    $item = 3;
    echo $item;
}
