<?php
interface VarianceTraitA {}
interface VarianceTraitB {}
class VarianceTraitBoth implements VarianceTraitA, VarianceTraitB {}
abstract class VarianceTraitBase {
    abstract public function run(VarianceTraitA&VarianceTraitB $value): VarianceTraitA;
}
trait VarianceTraitBody {
    public function run(VarianceTraitA $value): VarianceTraitA&VarianceTraitB { return new VarianceTraitBoth; }
}
class VarianceTraitFinal extends VarianceTraitBase { use VarianceTraitBody; }
echo gettype((new VarianceTraitFinal)->run(new VarianceTraitBoth)), "\n";
