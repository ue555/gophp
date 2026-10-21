<?php
function recoverInside() {
    try {
        yield 'ready';
    } catch (Exception $error) {
        yield $error->getMessage();
    }
}
function forwardException() {
    yield from recoverInside();
    yield 'done';
}
$generator = forwardException();
echo $generator->throw(new Exception('caught')), '|';
$generator->next();
echo $generator->current();
