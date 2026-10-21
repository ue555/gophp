<?php
class MagicBox {
    private array $items = [];
    public function __get(string $name): mixed { return $this->items[$name]; }
    public function __set(string $name, mixed $value): void { $this->items[$name] = $value; }
}
$box = new MagicBox;
$box->answer = 42;
echo $box->answer;
