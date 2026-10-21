<?php
$value = 1;
$closure = function &() use (&$value) {
    return $value;
};
$alias =& $closure();
$alias = 4;
$copy = $closure();
$copy = 9;
echo $value, '|', $closure();
