<?php
class DynamicKeywordTarget {
    public static function value() { return 7; }
    public static function check() {
        $class = 'self';
        try { return $class::value(); } catch (Error $e) { return 'missing'; }
    }
}
echo DynamicKeywordTarget::check(), '|';
$class = 'DynamicKeywordTarget';
echo $class::value(), "\n";
