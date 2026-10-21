<?php
class DynamicVisibleBase {
    private static $secret = 'base';
    protected static $shared = 2;
    private static function hidden() { return 'base-method'; }
    public static function baseValue($property, $method) { return self::${$property} . ':' . self::$method(); }
    public static function shared() { return static::$shared; }
}
class DynamicVisibleChild extends DynamicVisibleBase {
    private static $secret = 'child';
    public static function childValue($property) { return self::${$property}; }
    public static function inherited($property) { return parent::${$property}; }
}
$class = 'DynamicVisibleChild';
echo $class::baseValue('secret', 'hidden'), '|', $class::childValue('secret'), '|';
echo $class::inherited('shared'), '|', $class::shared(), '|';
try { echo $class::${'secret'}; } catch (Error $e) { echo 'private'; }
echo '|';
try { echo $class::${'shared'}; } catch (Error $e) { echo 'protected'; }
echo "\n";
