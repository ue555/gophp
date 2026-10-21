<?php
class CompoundBox { public int $value = 1; }
$box = new CompoundBox;
$result = ($box->value .= "2");
echo gettype($result), '|', $result, '|', gettype($box->value), '|', $box->value;
