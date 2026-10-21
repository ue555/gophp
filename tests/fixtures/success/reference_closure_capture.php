<?php
$number = 1;
$next = function () use (&$number) {
    $number++;
    return $number;
};
echo $next(), '|', $number, '|';
$number = 8;
echo $next(), '|', $number;
