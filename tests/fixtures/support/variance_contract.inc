<?php
class VarianceLoadedBase {}
class VarianceLoadedChild extends VarianceLoadedBase {}
interface VarianceLoadedContract {
    public function inspect(VarianceLoadedChild $value): VarianceLoadedBase;
}
