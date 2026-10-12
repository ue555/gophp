<?php
class Counter {
    private $value = 0;
    public function __construct($start = 1) { $this->value = $start; }
    public function next($step = 1) { $this->value = $this->value + $step; return $this->value; }
    public function self() { return $this; }
}
$a = new Counter(2); $b = new counter;
echo $a->next(), ':', $a->NEXT(3), ':', $b->next(), ':', ($a->self() === $a) + 0, "\n";
