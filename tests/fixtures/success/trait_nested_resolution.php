<?php
trait NestedLeft { public function value() { return 'left'; } }
trait NestedRight { public function value() { return 'right'; } }
trait ResolvedNested {
    use NestedLeft, NestedRight {
        NestedLeft::value insteadof NestedRight;
        NestedRight::value as protected rightValue;
    }
    public function both() { return $this->value() . '/' . $this->rightValue(); }
}
class NestedUser { use ResolvedNested; }
echo (new NestedUser)->both(), "\n";
