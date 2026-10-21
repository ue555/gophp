<?php
require __DIR__ . '/../support/variance_contract.inc';
class VarianceLoadedImpl implements VarianceLoadedContract {
    public function inspect(VarianceLoadedBase $value): VarianceLoadedChild { return new VarianceLoadedChild; }
}
echo gettype((new VarianceLoadedImpl)->inspect(new VarianceLoadedChild)), "\n";
