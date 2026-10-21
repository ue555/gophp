<?php
class DynamicPropertyBase { public static $count = 1; public static $other = 9; public static $maybe = null; }
class DynamicPropertyChild extends DynamicPropertyBase { public static $count = 10; }
$class = 'DynamicPropertyChild';
$name = 'count';
echo $class::$count, '|', $class::${$name}, '|';
$class::${$name} += 2;
echo ++$class::${$name}, '|', $class::${$name}++, '|';
echo isset($class::${$name}) ? 'yes' : 'no', '|';
echo empty($class::${$name}) ? 'empty' : 'set', '|';
echo $class::${$name} ?? 'fallback', '|';
echo DynamicPropertyBase::$count, '|', DynamicPropertyChild::$count, '|';
echo --$class::${$name}, '|', $class::${$name}--, '|';
$maybe = 'maybe';
echo $class::${$maybe} ??= 5;
echo '|', DynamicPropertyChild::$maybe, "\n";
