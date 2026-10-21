<?php
$values = ['package' => ['providers' => ['App\\Provider'], 'aliases' => []], 7 => true, 'nil' => null, 'text' => "a'b\\c\0d", 'float' => 1.25];
echo var_export($values, true), "\n";
var_export(false);
echo "|", var_export(2.0, true), "|", var_export(-0.0, true), "\n";
echo var_export([1e15, 1e20, 1e-5], true), "\n";
