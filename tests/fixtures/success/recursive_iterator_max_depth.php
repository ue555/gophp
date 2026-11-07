<?php
// RecursiveIteratorIterator::setMaxDepth()/getMaxDepth(): elements at the
// maximum depth are treated as leaves (Symfony Finder's depth() filter).
$data = ['a' => 1, 'b' => ['c' => 2, 'd' => ['e' => 3]], 'f' => 4];

final class Nested implements RecursiveIterator
{
    private array $keys;
    private int $position = 0;

    public function __construct(private array $items)
    {
        $this->keys = array_keys($items);
    }

    public function current(): mixed
    {
        return $this->items[$this->keys[$this->position]];
    }

    public function key(): mixed
    {
        return $this->keys[$this->position];
    }

    public function next(): void
    {
        $this->position++;
    }

    public function rewind(): void
    {
        $this->position = 0;
    }

    public function valid(): bool
    {
        return $this->position < count($this->keys);
    }

    public function hasChildren(): bool
    {
        return is_array($this->current());
    }

    public function getChildren(): RecursiveIterator
    {
        return new Nested($this->current());
    }
}

function walk(RecursiveIteratorIterator $iterator): string
{
    $parts = [];
    foreach ($iterator as $key => $value) {
        $parts[] = $key . '@' . $iterator->getDepth() . (is_array($value) ? '[]' : '=' . $value);
    }
    return implode(' ', $parts);
}

foreach ([RecursiveIteratorIterator::LEAVES_ONLY, RecursiveIteratorIterator::SELF_FIRST, RecursiveIteratorIterator::CHILD_FIRST] as $mode) {
    foreach ([-1, 0, 1] as $max) {
        $iterator = new RecursiveIteratorIterator(new Nested($data), $mode);
        $iterator->setMaxDepth($max);
        echo $mode, '/', var_export($iterator->getMaxDepth(), true), ': ', walk($iterator), "\n";
    }
}

$iterator = new RecursiveIteratorIterator(new Nested($data));
echo var_export($iterator->getMaxDepth(), true), "\n";
$iterator->setMaxDepth();
echo var_export($iterator->getMaxDepth(), true), "\n";
try {
    $iterator->setMaxDepth(-2);
} catch (ValueError $error) {
    echo get_class($error), ': ', $error->getMessage(), "\n";
}
