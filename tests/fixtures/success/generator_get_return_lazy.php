<?php
function returnLazy() {
    echo 'entered|';
    yield 'first';
    echo 'resumed|';
    return 'done';
}
$generator = returnLazy();
try {
    $generator->getReturn();
} catch (Exception $e) {
    echo 'pending|';
}
echo $generator->current(), '|';
try {
    $generator->getReturn();
} catch (Exception $e) {
    echo 'still-pending|';
}
$generator->next();
echo $generator->getReturn();
function returnWithoutYield() {
    echo '|empty-entered|';
    if (false) {
        yield 1;
    }
    return 7;
}
$empty = returnWithoutYield();
echo $empty->getReturn();
