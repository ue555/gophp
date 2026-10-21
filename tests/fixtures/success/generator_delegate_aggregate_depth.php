<?php
class AggregateLayers implements IteratorAggregate {
    public function __construct(public int $depth) {}

    public function getIterator(): Traversable {
        if ($this->depth === 0) {
            return new ArrayIterator([9]);
        }
        return new AggregateLayers($this->depth - 1);
    }
}

function delegateLayers() {
    yield from new AggregateLayers(24);
}

echo implode(',', iterator_to_array(delegateLayers()));
