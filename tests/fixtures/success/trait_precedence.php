<?php
class PrecedenceParent { public function value() { return "parent"; } }
trait PrecedenceTrait { public function value() { return "trait"; } }
class PrecedenceChild extends PrecedenceParent { use PrecedenceTrait; }
class PrecedenceOverride extends PrecedenceChild { public function value() { return "child"; } }
echo (new PrecedenceChild)->value(), '|', (new PrecedenceOverride)->value(), "\n";
