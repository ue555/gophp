<?php
class OneItem implements Iterator {
    private int $position = 0;
    public function rewind(): void { $this->position = 0; }
    public function valid(): bool { return $this->position === 0; }
    public function key(): int { return 0; }
    public function current(): int { return 7; }
    public function next(): void { $this->position++; }
}
function keepIterable(iterable $input): iterable { return $input; }
echo gettype(keepIterable([1])), '|';
$item = new OneItem;
echo keepIterable($item) === $item ? 'same|' : 'copy|';
echo is_iterable($item) ? 'yes|' : 'no|';
echo is_iterable(new stdClass) ? 'wrong|' : 'plain|';
try {
    keepIterable(1);
} catch (TypeError $e) {
    echo 'rejected';
}
