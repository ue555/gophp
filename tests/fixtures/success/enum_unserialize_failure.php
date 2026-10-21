<?php
enum RestoreStatus { case Ready; public const Alias = 1; }
class PlainStatus {}
$warnings = 0;
set_error_handler(function ($level, $message) use (&$warnings) {
    if ($level === E_WARNING) { $warnings++; }
    return true;
});
foreach (['E:13:"Missing:Ready";', 'E:17:"PlainStatus:Ready";', 'E:20:"RestoreStatus:Absent";', 'E:19:"RestoreStatus:Alias";', 'E:6:"Status";'] as $encoded) {
    echo unserialize($encoded) === false ? 'false|' : 'wrong|';
}
restore_error_handler();
echo $warnings, '|';
set_error_handler(function () { throw new Exception('stop'); });
try {
    unserialize('E:13:"Missing:Ready";');
    echo 'continued';
} catch (Exception $e) {
    echo 'caught';
}
restore_error_handler();
