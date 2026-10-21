<?php
function receiveInside() {
    $sent = yield 1;
    yield $sent;
    return 7;
}
function forwardToInner() {
    $result = yield from receiveInside();
    yield $result;
    return 9;
}
$generator = forwardToInner();
echo $generator->send('sent'), '|';
$generator->next();
echo $generator->current(), '|';
$generator->next();
echo $generator->getReturn();
