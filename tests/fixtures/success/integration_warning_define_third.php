<?php
function stopThirdNull($level, $message) { throw new Exception('stop'); }
set_error_handler('stopThirdNull');
try { define('AUDIT_THIRD_NULL', 9, null); } catch (Exception $e) { echo 'null|'; }
restore_error_handler();
echo defined('AUDIT_THIRD_NULL') ? 'created|' : 'clear|';

$sawWarning = false;
set_error_handler(function ($level, $message) {
    global $sawWarning;
    $sawWarning = true;
    return true;
});
define('AUDIT_THIRD_TRUE', 10, true);
restore_error_handler();
echo $sawWarning ? 'warned|' : 'silent|';
echo defined('AUDIT_THIRD_TRUE') ? 'created' : 'clear';
