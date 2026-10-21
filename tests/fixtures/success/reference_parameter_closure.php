<?php
$value = 1;
$callback = function (&$item) use (&$value) {
    $item += $value;
    $value++;
};
$callback($value);
echo $value, '|';
$callback($value);
echo $value;
