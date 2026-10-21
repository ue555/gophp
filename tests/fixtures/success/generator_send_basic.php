<?php
function receive() {
    $input = yield 1;
    yield $input;
    return 7;
}
$generator = receive();
echo $generator->send('hello'), '|';
$generator->next();
echo $generator->getReturn();
