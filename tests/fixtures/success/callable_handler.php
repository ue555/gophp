<?php
function reportCallableWarning(int $level, string $message): bool {
    echo $level, ':', $message, '|';
    return true;
}
$handler = reportCallableWarning(...);
set_error_handler($handler);
trigger_error('handled', E_USER_NOTICE);
restore_error_handler();
echo 'done';
