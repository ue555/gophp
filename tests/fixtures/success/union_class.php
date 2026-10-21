<?php
namespace Sample;
interface Named {}
class Base {}
class Child extends Base implements Named {}
class Other implements Named {}
use Sample\Base as ParentType;
function select(ParentType|Named|null $value): Named|ParentType|null { return $value; }
echo gettype(select(null)), '|', gettype(select(new Child)), '|', gettype(select(new Other)), "\n";
