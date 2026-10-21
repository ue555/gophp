<?php
interface VarianceUnionBase { public function change(int|string $value): int|string|bool; }
class VarianceUnionChild implements VarianceUnionBase {
    public function change(bool|string|int $value): string|int { return $value; }
}
echo gettype((new VarianceUnionChild)->change(4)), '|', gettype((new VarianceUnionChild)->change('x')), "\n";
