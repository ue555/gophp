<?php
$loads = 0;
try {
    require __DIR__ . '/../support/include_integration_nested.inc';
} catch (Exception $e) {
    echo 'outer-catch|';
} finally {
    echo 'outer-finally', "\n";
}
