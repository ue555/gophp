<?php
try {
    new LimitIterator(new ArrayIterator([]), -1);
    echo 'bad';
} catch (ValueError $e) {
    echo 'offset|';
}
try {
    new LimitIterator(new ArrayIterator([]), 0, -2);
    echo 'bad';
} catch (ValueError $e) {
    echo 'limit';
}
