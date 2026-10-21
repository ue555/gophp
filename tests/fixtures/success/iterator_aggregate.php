<?php
class InnerSequence implements Iterator {
    private int $position = 0;
    public function rewind(): void { $this->position = 0; }
    public function valid(): bool { return $this->position < 2; }
    public function key(): int { return $this->position; }
    public function current(): string { return $this->position === 0 ? 'a' : 'b'; }
    public function next(): void { $this->position++; }
}
class SequenceWrapper implements IteratorAggregate {
    public function getIterator(): Traversable { return new InnerSequence; }
}
foreach (new SequenceWrapper as $item) echo $item;
