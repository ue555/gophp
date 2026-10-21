<?php
class CompoundStaticBox { public static int $value = 1; }
$result = (CompoundStaticBox::$value .= "2");
echo gettype($result), '|', $result, '|', gettype(CompoundStaticBox::$value), '|', CompoundStaticBox::$value;
