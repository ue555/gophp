<?php
interface VarianceNarrowReturn { public function make(): int|string; }
echo "before\n";
class VarianceWideReturn implements VarianceNarrowReturn {
    public function make(): int|string|bool { return true; }
}
echo "after\n";
