<?php
class UntypedPromotion {
    public function __construct(public $value = 'seed') { echo $this->value, '|'; }
}
$a = new UntypedPromotion;
$b = new UntypedPromotion(9);
echo gettype($a->value), '|', $b->value, "\n";
