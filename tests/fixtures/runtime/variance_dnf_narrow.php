<?php
interface VarianceErrorA {}
interface VarianceErrorB {}
interface VarianceErrorContract { public function read(VarianceErrorA $value): VarianceErrorA; }
echo "before\n";
class VarianceErrorChild implements VarianceErrorContract {
    public function read(VarianceErrorA&VarianceErrorB $value): VarianceErrorA { return $value; }
}
echo "after\n";
