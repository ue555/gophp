<?php
class RefWarningBox {}
$box = new RefWarningBox;
set_error_handler(function ($level, $message) { throw new Exception('stop'); });
try { $alias =& $box->missing; } catch (Exception $e) { echo 'caught|'; }
restore_error_handler();
$calls = 0;
function countRefWarning($level, $message) { global $calls; $calls++; return true; }
set_error_handler('countRefWarning');
$box->missing = 4;
restore_error_handler();
echo $calls, '|', $box->missing;
