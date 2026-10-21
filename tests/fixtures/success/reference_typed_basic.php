<?php
class TypedBasic {
    public int $value = 1;
    public static int $saved = 2;
}
$box = new TypedBasic;
$alias =& $box->value;
$alias = '3';
$saved =& TypedBasic::$saved;
$saved = '4';
echo $box->value, '|', TypedBasic::$saved, '|', gettype($alias), '|', gettype($saved);
