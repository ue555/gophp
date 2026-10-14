<?php
class ForwardBase {
    public static function who() { return static::class; }
    public static function viaSelf() { return self::who(); }
}
class ForwardChild extends ForwardBase {
    public static function viaParent() { return parent::who(); }
}
echo ForwardChild::viaSelf(), '|', ForwardChild::viaParent(), "\n";
