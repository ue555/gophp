<?php
class SharedRoot {
    protected static $counter = 1;
    public static function add($n) { static::$counter += $n; return static::$counter; }
    public static function read() { return self::$counter; }
}
class SharedChild extends SharedRoot {}
class SeparateChild extends SharedRoot { protected static $counter = 20; }
echo SharedChild::add(2), '|', SharedRoot::read(), '|';
echo SeparateChild::add(3), '|', SeparateChild::read(), '|', SharedRoot::read(), "\n";
