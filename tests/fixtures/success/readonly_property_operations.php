<?php
class ReadonlyBox {
    public readonly int $number;
    public readonly object $inner;
    public readonly array $items;
    public function __construct() {
        $this->number = 4;
        $this->inner = new stdClass;
        $this->items = [1];
    }
}
$b = new ReadonlyBox;
$name = 'number';
$b->inner->value = 9;
echo $b->{$name}, '|', $b->inner->value, '|';
try { $b->{$name} = 5; } catch (Error $e) { echo 'write|'; }
try { $b->items[] = 2; } catch (Error $e) { echo 'array|'; }
try { unset($b->number); } catch (Error $e) { echo 'unset|'; }
echo $b->number, "\n";
