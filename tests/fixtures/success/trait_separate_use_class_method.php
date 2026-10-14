<?php
trait FirstOverride { public function value() { return 1; } }
trait SecondOverride { public function value() { return 2; } }
class SeparateOverride {
    use FirstOverride;
    use SecondOverride;
    public function value() { return 3; }
}
echo (new SeparateOverride)->value(), "\n";
