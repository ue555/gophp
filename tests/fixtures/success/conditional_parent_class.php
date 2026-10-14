<?php
if (true) {
    class ConditionalBase {
        public function value() { return 5; }
    }
}
class ConditionalChild extends ConditionalBase {}
echo (new ConditionalChild)->value(), "\n";
