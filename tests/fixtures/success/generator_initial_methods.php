<?php
function initialSequence() {
    echo 'begin|';
    yield 1;
    yield 2;
}
function emptySequence() {
    echo 'empty-begin|';
    if (false) yield 1;
}
function returnSequence() {
    echo 'return-begin|';
    yield 3;
    return 4;
}
$g = initialSequence();
echo $g->current(), '|';
$g->rewind();
echo $g->key(), '|';
$h = initialSequence();
$h->next();
echo $h->current(), '|';
$empty = emptySequence();
echo $empty->valid() ? 'valid' : 'empty';
echo '|';
$returns = returnSequence();
try {
    $returns->getReturn();
} catch (Exception $e) {
    echo 'pending|';
}
$returns->next();
echo $returns->getReturn();
