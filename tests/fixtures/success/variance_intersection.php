<?php
interface VariancePartA {}
interface VariancePartB {}
class VarianceBoth implements VariancePartA, VariancePartB {}
class VarianceIntersectionBase {
    public function check(VariancePartA&VariancePartB $value): VariancePartA { return $value; }
}
class VarianceIntersectionChild extends VarianceIntersectionBase {
    public function check(VariancePartA $value): VariancePartB&VariancePartA { return new VarianceBoth; }
}
echo gettype((new VarianceIntersectionChild)->check(new VarianceBoth)), "\n";
