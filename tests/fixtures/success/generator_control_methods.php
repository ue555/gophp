<?php
function controlSequence() {
    echo 'start|';
    yield 'first' => 10;
    echo 'middle|';
    yield 'second' => 20;
    echo 'end|';
    return 30;
}

$g = controlSequence();
echo 'created|';
echo $g->valid() ? 'valid|' : 'invalid|';
$g->rewind();
echo $g->key(), ':', $g->current(), '|';
try {
    $g->getReturn();
} catch (Exception $e) {
    echo 'not-returned|';
}
$g->next();
echo $g->key(), ':', $g->current(), '|';
try {
    $g->rewind();
} catch (Exception $e) {
    echo 'rewind-error|';
}
$g->next();
echo $g->valid() ? 'valid|' : 'closed|';
echo gettype($g->current()), '|', $g->getReturn();
