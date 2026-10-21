<?php
$calls = 0;
set_include_path(__DIR__ . '/../support/include_path/first');
set_error_handler(function ($level, $message) { throw new Exception('stopped'); });
try {
    include 'absent-path.inc';
    $calls++;
} catch (Exception $e) { echo $e->getMessage(), '|', $calls, "\n"; }
restore_error_handler();
