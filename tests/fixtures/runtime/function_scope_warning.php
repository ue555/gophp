<?php
$globalValue = 10;
function readLocal() { echo $globalValue; $globalValue = 2; return $globalValue; }
echo readLocal(), ':', readLocal(), ':', $globalValue, "\n";
