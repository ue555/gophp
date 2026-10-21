<?php
class RefStaticBase { public static $value = 1; }
class RefStaticChild extends RefStaticBase {}
class RefStaticOwn extends RefStaticBase { public static $value = 7; }
$alias =& RefStaticBase::$value;
RefStaticChild::$value = 4;
echo $alias, '|', RefStaticBase::$value, '|';
$alias = 5;
echo RefStaticChild::$value, '|', RefStaticOwn::$value;
