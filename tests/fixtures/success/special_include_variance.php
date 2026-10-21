<?php
require __DIR__ . '/../support/special_type_classes.inc';
class IncludedContract { public function make(): IncludedBase { return new IncludedBase; } }
class IncludedImplementation extends IncludedContract { public function make(): IncludedDerived { return new IncludedDerived; } }
echo (new IncludedImplementation)->make() instanceof IncludedBase ? 'included' : 'wrong';
