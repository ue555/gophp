<?php
class YieldingCollection implements IteratorAggregate {
    public function getIterator(): Traversable {
        $sent = yield 1;
        yield $sent;
        return 4;
    }
}
function delegateAggregate() {
    $result = yield from new YieldingCollection();
    yield $result;
}
$generator = delegateAggregate();
echo $generator->send('sent'), '|';
$generator->next();
echo $generator->current();
