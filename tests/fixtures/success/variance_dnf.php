<?php
interface VarianceDnfA {}
interface VarianceDnfB {}
interface VarianceDnfC {}
class VarianceDnfBoth implements VarianceDnfA, VarianceDnfB, VarianceDnfC {}
interface VarianceDnfContract {
    public function apply((VarianceDnfA&VarianceDnfB)|int $value): (VarianceDnfA&VarianceDnfB)|int;
}
class VarianceDnfImpl implements VarianceDnfContract {
    public function apply(VarianceDnfA|int|string $value): (VarianceDnfC&VarianceDnfB&VarianceDnfA)|int { return $value === 3 ? 3 : new VarianceDnfBoth; }
}
echo gettype((new VarianceDnfImpl)->apply(new VarianceDnfBoth)), '|', (new VarianceDnfImpl)->apply(3), "\n";
