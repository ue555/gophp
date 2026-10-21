<?php
function stopLoading($level) { throw new Exception('stop'); }
set_error_handler('stopLoading');
try {
    require __DIR__ . '/support/error_handler_absent.inc';
    echo 'unreachable';
} catch (Exception $e) {
    echo $e->getMessage(), '|';
}
restore_error_handler();
echo "continued\n";
