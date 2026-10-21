<?php
interface VarianceObjectContract { public function take(bool $value): mixed; }
class VarianceObjectImpl implements VarianceObjectContract {
    public function take($value): true { return true; }
}
interface VarianceReturnContract { public function make(): object; }
class VarianceReturnImpl implements VarianceReturnContract {
    public function make(): VarianceObjectImpl { return new VarianceObjectImpl; }
}
echo gettype((new VarianceObjectImpl)->take(false)), '|', gettype((new VarianceReturnImpl)->make()), "\n";
