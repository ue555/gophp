<?php
class FloatParamParent {
    public function value(int $number) { return $number; }
}
echo "before\n";
class FloatParamChild extends FloatParamParent {
    public function value(float $number) { return $number; }
}
echo "after\n";
