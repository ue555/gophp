<?php
class ObservedCachingInput implements Iterator {
    private int $position = 0;

    public function rewind(): void { echo 'r'; $this->position = 0; }
    public function valid(): bool { echo 'v'; return $this->position < 2; }
    public function current(): int { echo 'c'; return 10 + $this->position; }
    public function key(): int { echo 'k'; return $this->position; }
    public function next(): void { echo 'n'; ++$this->position; }
}

$cache = new CachingIterator(new ObservedCachingInput, 0);
$cache->rewind();
echo '|', $cache->current(), ':', $cache->key(), ':', $cache->hasNext() ? 'more' : 'end';
$cache->next();
echo '|', $cache->current(), ':', $cache->key(), ':', $cache->hasNext() ? 'more' : 'end';
