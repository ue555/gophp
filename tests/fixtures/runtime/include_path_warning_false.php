<?php
set_include_path(__DIR__ . '/../support/include_path/first');
function reportMissing($level, $message) { echo $level, '|'; return false; }
set_error_handler('reportMissing');
$result = include 'absent-path.inc';
echo ($result ? 'found' : 'missing'), "\n";
restore_error_handler();
