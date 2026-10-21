<?php
class ThrowTarget { public static function run() {} }
echo "before\n";
ThrowTarget::{throw new Exception('stop')}();
echo "after\n";
