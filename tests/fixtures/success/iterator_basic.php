<?php
class ShortSequence implements Iterator {
    private int $position = 0;
    public function rewind(): void { $this->position = 0; }
    public function valid(): bool { return $this->position < 3; }
    public function key(): int { return $this->position + 4; }
    public function current(): int { return $this->position * 2; }
    public function next(): void { $this->position++; }
}
foreach (new ShortSequence as $key => $item) {
    echo $key, ':', $item, ';';
}
