<?php
$loads = 0;
try {
    include_once __DIR__ . '/../support/include_integration_throw.inc';
} catch (Exception $e) {
    echo $e->getMessage(), '|';
}
echo include_once __DIR__ . '/../support/./include_integration_throw.inc', '|', $loads, "\n";
