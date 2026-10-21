<?php
class ThrowTarget { public static function run() {} }
echo "before\n";
function neverCalled() {
    ThrowTarget::{throw new Exception('stop')}();
}
