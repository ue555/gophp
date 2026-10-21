<?php
class DynamicGrand {
    public static function who() { return static::class; }
    public static function fromSelf($method) { return self::$method(); }
}
class DynamicParent extends DynamicGrand {
    public static function fromParent($method) { return parent::$method(); }
}
class DynamicChild extends DynamicParent {
    public static function fromStatic($method) { return static::$method(); }
}
$class = 'DynamicChild'; $method = 'who';
echo DynamicChild::fromSelf($method), '|', DynamicChild::fromParent($method), '|';
echo DynamicChild::fromStatic($method), '|', $class::$method(), '|';
echo DynamicGrand::$method(), "\n";
