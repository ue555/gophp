<?php
class ReferenceBag implements IteratorAggregate {
    public array $items = [1, 2];

    public function &getIterator() {
        yield $this->items[0];
        yield $this->items[1];
    }
}
$bag = new ReferenceBag;
foreach ($bag as &$item) {
    $item += 10;
}
unset($item);
echo $bag->items[0], '|', $bag->items[1];
