<?php
class VarianceAnimal {}
class VarianceDog extends VarianceAnimal {}
class VarianceBase {
    public function take(VarianceDog $dog): VarianceAnimal { return $dog; }
}
class VarianceChild extends VarianceBase {
    public function take(VarianceAnimal $dog): VarianceDog { return new VarianceDog; }
}
echo gettype((new VarianceChild)->take(new VarianceDog)), "\n";
