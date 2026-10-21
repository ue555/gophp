<?php
class RefScopeBase {
    public static $value = 1;
    public static function changeCalled() { $alias =& static::$value; $alias = 4; }
}
class RefScopeChild extends RefScopeBase {
    public static $value = 2;
    public static function changeParent() { $alias =& parent::$value; $alias = 5; }
}
RefScopeChild::changeCalled();
RefScopeChild::changeParent();
echo RefScopeBase::$value, '|', RefScopeChild::$value;
