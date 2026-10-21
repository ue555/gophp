<?php
class FloatUnionReturnParent {
    public function value(): float|string { return 1.0; }
}
echo "before\n";
class FloatUnionReturnChild extends FloatUnionReturnParent {
    public function value(): int|string { return 1; }
}
echo "after\n";
