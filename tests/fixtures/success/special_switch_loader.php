<?php
switch (1) {
    case 1:
        require __DIR__ . '/../support/special_type_classes.inc';
        break;
}
class SwitchContract { public function make(): IncludedBase { return new IncludedBase; } }
class SwitchImplementation extends SwitchContract { public function make(): IncludedDerived { return new IncludedDerived; } }
echo (new SwitchImplementation)->make() instanceof IncludedBase ? 'switch' : 'wrong';
