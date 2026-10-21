<?php
class FloatUnionParamParent {
    public function value(int|string $number) { return $number; }
}
echo "before\n";
class FloatUnionParamChild extends FloatUnionParamParent {
    public function value(float|string $number) { return $number; }
}
echo "after\n";
