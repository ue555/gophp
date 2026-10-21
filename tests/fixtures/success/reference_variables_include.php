<?php
$source = 1;
$alias =& $source;
$result = require __DIR__ . '/../support/reference_variables_include.inc';
echo $source, '|', $alias, '|', $result;
