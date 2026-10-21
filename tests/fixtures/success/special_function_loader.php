<?php
function loadSpecialTypes() { require __DIR__ . '/../support/special_type_classes.inc'; }
loadSpecialTypes();
class FunctionContract { public function make(): IncludedBase { return new IncludedBase; } }
class FunctionImplementation extends FunctionContract { public function make(): IncludedDerived { return new IncludedDerived; } }
echo (new FunctionImplementation)->make() instanceof IncludedBase ? 'function' : 'wrong';
