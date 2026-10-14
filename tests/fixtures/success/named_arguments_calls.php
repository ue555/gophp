<?php
class NamedCall {
    public function __construct($left, $right) { echo $left, $right, '|'; }
    public function method($left, $right) { echo $left, $right, '|'; }
    public static function staticMethod($left, $right) { echo $left, $right, '|'; }
}
$object = new NamedCall(right: 2, left: 1);
$object->method(right: 4, left: 3);
NamedCall::staticMethod(right: 6, left: 5);
$closure = function($left, $right) { echo $left, $right, '|'; };
$closure(right: 8, left: 7);
echo strlen(string: 'abc'), '|', count(mode: COUNT_NORMAL, value: [1, 2]), "\n";
