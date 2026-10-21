<?php
$shared = 10;
function restoreState() {
    static $n = 1;
    global $shared;
    $loads = 0;
    try {
        require __DIR__ . '/../support/include_integration_throw.inc';
    } catch (Exception $e) {
        echo $e->getMessage(), '|';
    }
    $shared++;
    $n++;
    echo (__FILE__ === __DIR__ . '/include_integration_restore_after_catch.php' ? 'file|' : 'wrong|');
    echo $n, '|', $shared, '|';
}
restoreState();
restoreState();
echo "done\n";
