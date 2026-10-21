<?php
trait DynamicStaticTrait {
    public static function title($value) { return static::class . ':' . $value; }
}
class DynamicStaticOwner { use DynamicStaticTrait; }
$class = 'DynamicStaticOwner'; $method = 'title';
echo $class::$method(value: 'ok'), '|', $class::{$method}(...['value' => 'yes']), "\n";
