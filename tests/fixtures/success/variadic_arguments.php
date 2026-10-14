<?php
function collect($head, ...$tail) {
    echo count($tail), '|';
    foreach ($tail as $item) echo $item;
    return $head;
}
class VariadicBox {
    public function size(...$items) { return count($items); }
}
$number = 0;
function nextItem() { global $number; return ++$number; }
echo collect('A', nextItem(), nextItem()), '|';
echo (new VariadicBox)->size(4, 5), '|', $number, "\n";
