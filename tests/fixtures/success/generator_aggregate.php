<?php
class Sequence implements IteratorAggregate {
    public function getIterator(): Traversable {
        yield 1;
        yield 2;
    }
}
foreach (new Sequence as $value) {
    echo $value, '|';
}
