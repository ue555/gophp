<?php
class DynamicUnselected { public static $value = 1; public static function run() { return 4; } }
$class = 'DynamicUnselected'; $name = 'run';
if (false) { echo $missingClass::${$missingName}; echo $class::$missingMethod(); }
echo $class::$value, '|', $class::$name(), "\n";
