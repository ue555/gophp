<?php
class ForwardThing implements FutureContract { use FutureTrait; }
interface FutureContract { public function value(); }
trait FutureTrait { public function value() { return 12; } }
echo (new ForwardThing)->value(), "\n";
