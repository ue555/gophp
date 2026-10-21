<?php
class WrappedValues implements IteratorAggregate {
    public function getIterator(): Traversable {
        return new ArrayIterator([2 => 'two', 4 => 'four']);
    }
}
function delegateIterator() {
    yield from new WrappedValues();
    yield 'after';
}
foreach (delegateIterator() as $key => $item) {
    echo $key, ':', $item, '|';
}
