<?php declare(strict_types=1);
class StrictDynamicBox { public int $value = 4; }
$box = new StrictDynamicBox;
$name = 'value';
try { $box->{$name} = '5'; } catch (TypeError $e) { echo 'TypeError|'; }
echo $box->{$name}, "\n";
