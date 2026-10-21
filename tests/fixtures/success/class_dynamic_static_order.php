<?php
class DynamicOrder { public static $value = 0; public static function take($x) { return $x; } }
function orderClass() { echo 'class|'; return 'DynamicOrder'; }
function orderMember() { echo 'member|'; return 'take'; }
function orderArgument() { echo 'arg|'; return 7; }
function orderProperty() { echo 'property|'; return 'value'; }
function orderValue() { echo 'rhs|'; return 3; }
function orderThrow() { throw new Exception('stop'); }
echo orderClass()::{orderMember()}(orderArgument()), '|';
orderClass()::${orderProperty()} = orderValue();
echo DynamicOrder::$value, '|';
try { orderClass()::{orderThrow()}(orderArgument()); }
catch (Exception $e) { echo 'caught'; }
echo "\n";
