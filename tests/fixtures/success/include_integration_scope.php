<?php
$shared = 3;
function readScope() {
    $local = 4;
    try {
        $result = require __DIR__ . '/../support/include_integration_scope.inc';
    } finally {
        echo 'finally|';
    }
    echo (__FILE__ === __DIR__ . '/include_integration_scope.php' ? 'restored|' : 'wrong-file|');
    echo $local, '|', $result, '|';
}
readScope();
echo $shared, "\n";
