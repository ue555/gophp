<?php
class RefTargetBox { public $value = 1; public static $saved = 2; }
$box = new RefTargetBox;
$source = 3;
$box->value =& $source;
$source = 4;
echo $box->value, '|', $source, '|';
$other = 5;
RefTargetBox::$saved =& $other;
$other = 6;
echo RefTargetBox::$saved, '|', $other;
