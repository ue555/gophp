<?php
$set = function () use (&$value) {
    $value = 7;
};
echo isset($value) ? 'set' : 'missing', '|';
$set();
echo $value;
