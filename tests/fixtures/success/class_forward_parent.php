<?php
class ForwardChild extends ForwardBase {}
class ForwardBase {
    public function value() { return 7; }
}
echo (new ForwardChild)->value(), "\n";
