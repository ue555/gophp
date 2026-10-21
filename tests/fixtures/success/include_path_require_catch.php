<?php
set_include_path(__DIR__ . '/../support/include_path/first');
function quietMissing($level, $message) { echo $level, '|'; return true; }
set_error_handler('quietMissing');
try { require 'absent-path.inc'; } catch (Error $e) { echo 'caught|'; }
restore_error_handler();
echo "after\n";
