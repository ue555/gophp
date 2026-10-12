<?php
class Defaults {
    public $a, $b = 2 + 3 * 4, $data = ['x' => [1]];
    function change() { $this->data['x'][] = 2; }
}
$a = new Defaults; $b = new Defaults;
$a->change();
echo ($a->a === null) + 0, ':', $a->b, ':', ($a->data === ['x' => [1, 2]]) + 0, ($b->data === ['x' => [1]]) + 0, "\n";
