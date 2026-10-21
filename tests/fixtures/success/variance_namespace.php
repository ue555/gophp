<?php
namespace VarianceNames;
use VarianceNames\Base as ImportedBase;
class Base {}
class Child extends Base {}
interface Contract {
    public function build(Child $value): ImportedBase;
}
class Implementation implements Contract {
    public function build(ImportedBase $value): Child { return new Child; }
}
echo gettype((new Implementation)->build(new Child)), "\n";
