<?php
class Store implements ArrayAccess {
    private array $items = [];
    public function offsetExists(mixed $offset): bool { return isset($this->items[$offset]); }
    public function offsetGet(mixed $offset): mixed { return $this->items[$offset]; }
    public function offsetSet(mixed $offset, mixed $value): void { $this->items[$offset] = $value; }
    public function offsetUnset(mixed $offset): void { unset($this->items[$offset]); }
}
$store = new Store;
$store['x'] = 5;
echo $store['x'];
