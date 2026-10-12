<?php
class ConstructorReturn {
    public $n;
    function __construct($n = 1) { $this->n = $n; return 7; }
}
$o = new ConstructorReturn(2);
echo $o->n, ':', $o->__construct(3), ':', $o->n, "\n";
