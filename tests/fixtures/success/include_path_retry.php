<?php
set_error_handler(function ($level, $message) { return true; });
set_include_path(__DIR__ . '/../support/include_path/first');
$missing = include_once 'retry.inc';
echo ($missing ? 'found' : 'missing'), '|';
set_include_path(__DIR__ . '/../support/include_path/second');
echo include_once 'retry.inc', "\n";
restore_error_handler();
