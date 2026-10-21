<?php
class RefStaticVisible {
    private static $value = 1;
    public static function change() { $alias =& self::$value; $alias = 4; }
    public static function get() { return self::$value; }
}
RefStaticVisible::change();
echo RefStaticVisible::get(), '|';
try { $alias =& RefStaticVisible::$value; } catch (Error $e) { echo 'blocked'; }
