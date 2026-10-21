<?php
function recoverInside() {
    try {
        yield 'ready';
    } catch (Exception $e) {
        yield $e->getMessage();
    }
}
$generator = recoverInside();
echo $generator->throw(new Exception('caught'));
