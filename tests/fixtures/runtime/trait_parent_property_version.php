<?php
class InheritedProperty { public $value = 1; }
trait DifferentProperty { public $value = 2; }
class PropertyUser extends InheritedProperty { use DifferentProperty; }
echo (new PropertyUser)->value, "\n";
