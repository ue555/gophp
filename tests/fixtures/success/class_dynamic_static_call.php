<?php
class DynamicCallBase {
    public static function join($first, $second = 'default', ...$rest) {
        echo static::class, ':', $first, ':', $second, ':', count($rest), '|';
    }
}
class DynamicCallChild extends DynamicCallBase {}
$class = 'DynamicCallChild';
$method = 'join';
$class::join('a');
DynamicCallChild::$method(second: 'b', first: 'a');
$class::$method(...['first' => 'x', 'second' => 'y']);
$class::{"jo" . "in"}('z');
