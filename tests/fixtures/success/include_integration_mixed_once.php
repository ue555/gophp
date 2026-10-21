<?php
$path = __DIR__ . '/../support/include_integration_return.inc';
$first = include $path;
$second = require_once __DIR__ . '/../support/./include_integration_return.inc';
$third = require $path;
echo $first, '|', $second, '|', $third, "\n";
