<?php
abstract class AbstractBase {
    abstract public function value();
    public function double() { return $this->value() * 2; }
}
class AbstractChild extends AbstractBase {
    public function value() { return 6; }
}
echo (new AbstractChild)->double(), "\n";
