<?php
if (false) {
    interface UnusedContract { public function missing(); }
    trait UnusedTrait { public function broken() { return 0; } }
}
if (true) {
    interface ActiveContract { public function value(); }
    trait ActiveTrait { public function value() { return 6; } }
}
class ActiveThing implements ActiveContract { use ActiveTrait; }
$x = new ActiveThing;
echo $x->value(), '|', $x instanceof ActiveContract, "\n";
