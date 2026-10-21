<?php
class NormalModuloBox { public float $value = 5.5; }
$number = 5.5;
$array = [5.5];
$box = new NormalModuloBox;
$warnings = 0;
set_error_handler(function ($level, $message) { global $warnings; ++$warnings; return true; });
$number %= 2;
$array[0] %= 2;
$box->value %= 2;
restore_error_handler();
echo $number, '|', $array[0], '|', $box->value, '|', $warnings;
