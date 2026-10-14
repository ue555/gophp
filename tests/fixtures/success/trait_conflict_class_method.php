<?php
trait FirstValue { public function value() { return 1; } }
trait SecondValue { public function value() { return 2; } }
class OwnValue {
    use FirstValue, SecondValue;
    public function value() { return 3; }
}
echo (new OwnValue)->value(), "\n";
