<?php
interface VarianceWideParam { public function take(int|string $value); }
echo "before\n";
class VarianceNarrowParam implements VarianceWideParam {
    public function take(int $value) {}
}
echo "after\n";
