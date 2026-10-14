<?php
class StaticBase {
    public static $count = 1;
    public const LABEL = "base";
    public static function bump() { static::$count++; return static::LABEL; }
    public static function make() { return new static; }
    public static function names() { return self::class . "/" . static::class; }
}
class StaticChild extends StaticBase {
    public static $count = 10;
    public const LABEL = "child";
}
echo StaticBase::bump(), ':', StaticBase::$count, '|';
echo StaticChild::bump(), ':', StaticChild::$count, '|';
echo (StaticChild::make() instanceof StaticChild), '|', StaticChild::names(), "\n";
