<?php
$counter = 0;
$first = include_once __DIR__ . '/../support/include_value.inc';
$second = include_once __DIR__ . '/../support/./include_value.inc';
echo $first, '|', $second, '|', $counter, "\n";
