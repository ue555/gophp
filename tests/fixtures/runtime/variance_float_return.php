<?php
class FloatReturnParent {
    public function value(): float { return 1.0; }
}
echo "before\n";
class FloatReturnChild extends FloatReturnParent {
    public function value(): int { return 1; }
}
echo "after\n";
