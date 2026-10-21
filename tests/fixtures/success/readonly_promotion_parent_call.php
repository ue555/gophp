<?php
class PromotedParentCall {
    public function __construct(public readonly int $base) {}
}
class PromotedChildCall extends PromotedParentCall {
    public function __construct(int $base, public readonly int $child) {
        echo $this->child, '|';
        parent::__construct($base);
    }
}
$o = new PromotedChildCall(3, 5);
echo $o->base, '|', $o->child, "\n";
