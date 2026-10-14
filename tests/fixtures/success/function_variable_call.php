<?php
function triple($value) { return $value * 3; }
$callable = 'triple';
echo $callable(7), "\n";
