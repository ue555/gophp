<?php
class ThrowTarget { public static function run() {} }
echo "before\n";
if (false) {
    ThrowTarget::{throw new Exception('stop')}();
}
