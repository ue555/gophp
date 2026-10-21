<?php
interface VarianceTop { public function process(int|string $value): object; }
interface VarianceLeft extends VarianceTop {}
interface VarianceRight extends VarianceTop {}
interface VarianceDiamond extends VarianceLeft, VarianceRight {}
class VarianceDiamondImpl implements VarianceDiamond {
    public function process(int|string|bool $value): VarianceDiamondImpl { return $this; }
}
echo gettype((new VarianceDiamondImpl)->process(true)), "\n";
