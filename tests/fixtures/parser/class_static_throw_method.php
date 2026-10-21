<?php
class ThrowTarget { public static function run() {} }
echo "before\n";
class NeverCalled {
    public function unused() {
        ThrowTarget::{throw new Exception('stop')}();
    }
}
