<?php
interface VarianceVariadicContract { public function choose(int|string ...$values): int|string; }
class VarianceVariadicImpl implements VarianceVariadicContract {
    public function choose(int|string|bool ...$values): int { return count($values); }
}
echo (new VarianceVariadicImpl)->choose(1, "two", true), "\n";
