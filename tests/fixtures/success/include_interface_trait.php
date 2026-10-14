<?php
require __DIR__ . '/../support/item7_contracts.inc';
class IncludedThing implements IncludedContract { use IncludedImplementation; }
$x = new IncludedThing;
echo $x->value(), '|', $x instanceof IncludedContract, "\n";
